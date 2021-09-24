package discovery

import (
	"github.com/bmizerany/assert"
	"github.com/praetorian-inc/mithril/pkg/envoy"
	"testing"
)

const istioGkeDumpJson = `{
 "configs": [
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.BootstrapConfigDump",
   "bootstrap": {
    "node": {
     "id": "sidecar~10.124.2.21~httpbin1-6c97fbc967-lz4fb.default~default.svc.cluster.local",
     "cluster": "httpbin1.default",
     "metadata": {
      "NAME": "httpbin1-6c97fbc967-lz4fb",
      "NAMESPACE": "default",
      "INCLUDE_INBOUND_PORTS": "",
      "app": "httpbin1",
      "EXCHANGE_KEYS": "NAME,NAMESPACE,INSTANCE_IPS,LABELS,OWNER,PLATFORM_METADATA,WORKLOAD_NAME,CANONICAL_TELEMETRY_SERVICE,MESH_ID,SERVICE_ACCOUNT",
      "INSTANCE_IPS": "10.124.2.21",
      "pod-template-hash": "6c97fbc967",
      "POD_PORTS": "[\n]",
      "kubectl.kubernetes.io/restartedAt": "2021-09-24T12:02:29-05:00",
      "INTERCEPTION_MODE": "REDIRECT",
      "SERVICE_ACCOUNT": "default",
      "CONFIG_NAMESPACE": "default",
      "OWNER": "kubernetes://apis/apps/v1/namespaces/default/deployments/httpbin1",
      "WORKLOAD_NAME": "httpbin1",
      "ISTIO_VERSION": "1.4.10-gke.8",
      "POD_NAME": "httpbin1-6c97fbc967-lz4fb",
      "PLATFORM_METADATA": {
       "gcp_project": "service-mithril",
       "gcp_location": "us-central1-c",
       "gcp_gce_instance_id": "2876453198940358853",
       "gcp_gke_cluster_name": "cluster-gke-istio"
      },
      "CLUSTER_ID": "Kubernetes",
      "LABELS": {
       "app": "httpbin1",
       "pod-template-hash": "6c97fbc967"
      },
      "ISTIO_PROXY_SHA": "istio-proxy:967db67f6c22a1572ff53b554ac2079401fb25eb"
     },
     "locality": {
      "region": "us-central1",
      "zone": "us-central1-c"
     },
     "build_version": "967db67f6c22a1572ff53b554ac2079401fb25eb/1.12.6/Clean/RELEASE/BoringSSL-FIPS"
    },
    "static_resources": {
     "listeners": [
      {
       "address": {
        "socket_address": {
         "address": "0.0.0.0",
         "port_value": 15090
        }
       },
       "filter_chains": [
        {
         "filters": [
          {
           "name": "envoy.http_connection_manager",
           "config": {
            "route_config": {
             "virtual_hosts": [
              {
               "routes": [
                {
                 "match": {
                  "prefix": "/stats/prometheus"
                 },
                 "route": {
                  "cluster": "prometheus_stats"
                 }
                }
               ],
               "domains": [
                "*"
               ],
               "name": "backend"
              }
             ]
            },
            "codec_type": "AUTO",
            "http_filters": {
             "name": "envoy.router"
            },
            "stat_prefix": "stats"
           }
          }
         ]
        }
       ]
      }
     ],
     "clusters": [
      {
       "name": "prometheus_stats",
       "type": "STATIC",
       "connect_timeout": "0.250s",
       "hosts": [
        {
         "socket_address": {
          "address": "127.0.0.1",
          "port_value": 15000
         }
        }
       ]
      },
      {
       "name": "xds-grpc",
       "type": "STRICT_DNS",
       "connect_timeout": "10s",
       "hosts": [
        {
         "socket_address": {
          "address": "istio-pilot.istio-system",
          "port_value": 15010
         }
        }
       ],
       "circuit_breakers": {
        "thresholds": [
         {
          "max_connections": 100000,
          "max_pending_requests": 100000,
          "max_requests": 100000
         },
         {
          "priority": "HIGH",
          "max_connections": 100000,
          "max_pending_requests": 100000,
          "max_requests": 100000
         }
        ]
       },
       "http2_protocol_options": {},
       "dns_refresh_rate": "300s",
       "dns_lookup_family": "V4_ONLY",
       "upstream_connection_options": {
        "tcp_keepalive": {
         "keepalive_time": 300
        }
       }
      },
      {
       "name": "zipkin",
       "type": "STRICT_DNS",
       "connect_timeout": "1s",
       "hosts": [
        {
         "socket_address": {
          "address": "zipkin.istio-system",
          "port_value": 9411
         }
        }
       ],
       "dns_refresh_rate": "300s",
       "dns_lookup_family": "V4_ONLY"
      }
     ]
    },
    "dynamic_resources": {
     "lds_config": {
      "ads": {}
     },
     "cds_config": {
      "ads": {}
     },
     "ads_config": {
      "api_type": "GRPC",
      "grpc_services": [
       {
        "envoy_grpc": {
         "cluster_name": "xds-grpc"
        }
       }
      ]
     }
    },
    "tracing": {
     "http": {
      "name": "envoy.zipkin",
      "config": {
       "trace_id_128bit": "true",
       "shared_span_context": "false",
       "collector_endpoint": "/api/v1/spans",
       "collector_cluster": "zipkin"
      }
     }
    },
    "admin": {
     "access_log_path": "/dev/null",
     "address": {
      "socket_address": {
       "address": "127.0.0.1",
       "port_value": 15000
      }
     }
    },
    "stats_config": {
     "stats_tags": [
      {
       "tag_name": "cluster_name",
       "regex": "^cluster\\.((.+?(\\..+?\\.svc\\.cluster\\.local)?)\\.)"
      },
      {
       "tag_name": "tcp_prefix",
       "regex": "^tcp\\.((.*?)\\.)\\w+?$"
      },
      {
       "tag_name": "response_code",
       "regex": "(response_code=\\.=(.+?);\\.;)|_rq(_(\\.d{3}))$"
      },
      {
       "tag_name": "response_code_class",
       "regex": "_rq(_(\\dxx))$"
      },
      {
       "tag_name": "http_conn_manager_listener_prefix",
       "regex": "^listener(?=\\.).*?\\.http\\.(((?:[_.[:digit:]]*|[_\\[\\]aAbBcCdDeEfF[:digit:]]*))\\.)"
      },
      {
       "tag_name": "http_conn_manager_prefix",
       "regex": "^http\\.(((?:[_.[:digit:]]*|[_\\[\\]aAbBcCdDeEfF[:digit:]]*))\\.)"
      },
      {
       "tag_name": "listener_address",
       "regex": "^listener\\.(((?:[_.[:digit:]]*|[_\\[\\]aAbBcCdDeEfF[:digit:]]*))\\.)"
      },
      {
       "tag_name": "mongo_prefix",
       "regex": "^mongo\\.(.+?)\\.(collection|cmd|cx_|op_|delays_|decoding_)(.*?)$"
      },
      {
       "tag_name": "reporter",
       "regex": "(reporter=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_namespace",
       "regex": "(source_namespace=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_workload",
       "regex": "(source_workload=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_workload_namespace",
       "regex": "(source_workload_namespace=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_principal",
       "regex": "(source_principal=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_app",
       "regex": "(source_app=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "source_version",
       "regex": "(source_version=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_namespace",
       "regex": "(destination_namespace=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_workload",
       "regex": "(destination_workload=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_workload_namespace",
       "regex": "(destination_workload_namespace=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_principal",
       "regex": "(destination_principal=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_app",
       "regex": "(destination_app=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_version",
       "regex": "(destination_version=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_service",
       "regex": "(destination_service=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_service_name",
       "regex": "(destination_service_name=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "destination_service_namespace",
       "regex": "(destination_service_namespace=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "request_protocol",
       "regex": "(request_protocol=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "response_flags",
       "regex": "(response_flags=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "connection_security_policy",
       "regex": "(connection_security_policy=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "permissive_response_code",
       "regex": "(permissive_response_code=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "permissive_response_policyid",
       "regex": "(permissive_response_policyid=\\.=(.+?);\\.;)"
      },
      {
       "tag_name": "cache",
       "regex": "(cache\\.(.+?)\\.)"
      },
      {
       "tag_name": "component",
       "regex": "(component\\.(.+?)\\.)"
      },
      {
       "tag_name": "tag",
       "regex": "(tag\\.(.+?)\\.)"
      }
     ],
     "use_all_default_tags": false,
     "stats_matcher": {
      "inclusion_list": {
       "patterns": [
        {
         "prefix": "reporter="
        },
        {
         "prefix": "component"
        },
        {
         "prefix": "cluster_manager"
        },
        {
         "prefix": "listener_manager"
        },
        {
         "prefix": "http_mixer_filter"
        },
        {
         "prefix": "tcp_mixer_filter"
        },
        {
         "prefix": "server"
        },
        {
         "prefix": "cluster.xds-grpc"
        },
        {
         "suffix": "ssl_context_update_by_sds"
        }
       ]
      }
     }
    }
   },
   "last_updated": "2021-09-24T17:02:35.929Z"
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.ClustersConfigDump",
   "version_info": "2021-09-24T16:05:58Z/13",
   "static_clusters": [
    {
     "cluster": {
      "name": "prometheus_stats",
      "type": "STATIC",
      "connect_timeout": "0.250s",
      "hosts": [
       {
        "socket_address": {
         "address": "127.0.0.1",
         "port_value": 15000
        }
       }
      ]
     },
     "last_updated": "2021-09-24T17:02:35.936Z"
    },
    {
     "cluster": {
      "name": "xds-grpc",
      "type": "STRICT_DNS",
      "connect_timeout": "10s",
      "hosts": [
       {
        "socket_address": {
         "address": "istio-pilot.istio-system",
         "port_value": 15010
        }
       }
      ],
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 100000,
         "max_pending_requests": 100000,
         "max_requests": 100000
        },
        {
         "priority": "HIGH",
         "max_connections": 100000,
         "max_pending_requests": 100000,
         "max_requests": 100000
        }
       ]
      },
      "http2_protocol_options": {},
      "dns_refresh_rate": "300s",
      "dns_lookup_family": "V4_ONLY",
      "upstream_connection_options": {
       "tcp_keepalive": {
        "keepalive_time": 300
       }
      }
     },
     "last_updated": "2021-09-24T17:02:35.966Z"
    },
    {
     "cluster": {
      "name": "zipkin",
      "type": "STRICT_DNS",
      "connect_timeout": "1s",
      "hosts": [
       {
        "socket_address": {
         "address": "zipkin.istio-system",
         "port_value": 9411
        }
       }
      ],
      "dns_refresh_rate": "300s",
      "dns_lookup_family": "V4_ONLY"
     },
     "last_updated": "2021-09-24T17:02:35.967Z"
    }
   ],
   "dynamic_active_clusters": [
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "BlackHoleCluster",
      "type": "STATIC",
      "connect_timeout": "10s"
     },
     "last_updated": "2021-09-24T17:02:36.736Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "InboundPassthroughClusterIpv4",
      "type": "ORIGINAL_DST",
      "connect_timeout": "10s",
      "lb_policy": "CLUSTER_PROVIDED",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "upstream_bind_config": {
       "source_address": {
        "address": "127.0.0.6",
        "port_value": 0
       }
      }
     },
     "last_updated": "2021-09-24T17:02:36.737Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "PassthroughCluster",
      "type": "ORIGINAL_DST",
      "connect_timeout": "10s",
      "lb_policy": "CLUSTER_PROVIDED",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.736Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "inbound|15020|mgmt-15020|mgmtCluster",
      "type": "STATIC",
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "load_assignment": {
       "cluster_name": "inbound|15020|mgmt-15020|mgmtCluster",
       "endpoints": [
        {
         "lb_endpoints": [
          {
           "endpoint": {
            "address": {
             "socket_address": {
              "address": "127.0.0.1",
              "port_value": 15020
             }
            }
           }
          }
         ]
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.737Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15004||istio-policy.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15004||istio-policy.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "tls_context": {
       "common_tls_context": {
        "tls_certificates": [
         {
          "certificate_chain": {
           "filename": "/etc/certs/cert-chain.pem"
          },
          "private_key": {
           "filename": "/etc/certs/key.pem"
          }
         }
        ],
        "validation_context": {
         "trusted_ca": {
          "filename": "/etc/certs/root-cert.pem"
         },
         "verify_subject_alt_name": [
          "spiffe://cluster.local/ns/istio-system/sa/istio-mixer-service-account"
         ]
        },
        "alpn_protocols": [
         "istio",
         "h2"
        ]
       },
       "sni": "outbound_.15004_._.istio-policy.istio-system.svc.cluster.local"
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-policy"
        }
       }
      }
     },
     "last_updated": "2021-09-24T17:02:36.705Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15004||istio-telemetry.istio-system.svc.cluster.local",
      "type": "STRICT_DNS",
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "tls_context": {
       "common_tls_context": {
        "tls_certificates": [
         {
          "certificate_chain": {
           "filename": "/etc/certs/cert-chain.pem"
          },
          "private_key": {
           "filename": "/etc/certs/key.pem"
          }
         }
        ],
        "validation_context": {
         "trusted_ca": {
          "filename": "/etc/certs/root-cert.pem"
         },
         "verify_subject_alt_name": [
          "spiffe://cluster.local/ns/istio-system/sa/istio-mixer-service-account"
         ]
        },
        "alpn_protocols": [
         "istio",
         "h2"
        ]
       },
       "sni": "outbound_.15004_._.istio-telemetry.istio-system.svc.cluster.local"
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-telemetry"
        }
       }
      },
      "load_assignment": {
       "cluster_name": "outbound|15004||istio-telemetry.istio-system.svc.cluster.local",
       "endpoints": [
        {
         "lb_endpoints": [
          {
           "endpoint": {
            "address": {
             "socket_address": {
              "address": "10.0.0.136",
              "port_value": 15004
             }
            }
           }
          }
         ]
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.731Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15010||istio-pilot.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15010||istio-pilot.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      }
     },
     "last_updated": "2021-09-24T17:02:36.679Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15010||istiod-istio-1611.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15010||istiod-istio-1611.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      }
     },
     "last_updated": "2021-09-24T17:02:36.734Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15011||istio-pilot.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15011||istio-pilot.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.679Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15012||istiod-istio-1611.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15012||istiod-istio-1611.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.735Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-citadel.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istio-citadel.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.676Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-galley.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istio-galley.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.677Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-pilot.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istio-pilot.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.680Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-policy.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istio-policy.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "max_requests_per_connection": 10000,
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 10000,
         "max_retries": 4294967295
        }
       ]
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-policy"
        }
       }
      }
     },
     "last_updated": "2021-09-24T17:02:36.706Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-sidecar-injector.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istio-sidecar-injector.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.707Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istio-telemetry.istio-system.svc.cluster.local",
      "type": "STRICT_DNS",
      "connect_timeout": "10s",
      "max_requests_per_connection": 10000,
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 10000,
         "max_retries": 4294967295
        }
       ]
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-telemetry"
        }
       }
      },
      "load_assignment": {
       "cluster_name": "outbound|15014||istio-telemetry.istio-system.svc.cluster.local",
       "endpoints": [
        {
         "lb_endpoints": [
          {
           "endpoint": {
            "address": {
             "socket_address": {
              "address": "10.0.0.136",
              "port_value": 15014
             }
            }
           }
          }
         ]
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.731Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15014||istiod-istio-1611.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15014||istiod-istio-1611.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.735Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15020||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15020||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "protocol_selection": "USE_DOWNSTREAM_PROTOCOL"
     },
     "last_updated": "2021-09-24T17:02:36.677Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15029||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15029||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.678Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15030||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15030||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.678Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15031||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15031||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.678Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15032||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15032||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.679Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|15443||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|15443||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.679Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|31400||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|31400||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.678Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|42422||istio-telemetry.istio-system.svc.cluster.local",
      "type": "STRICT_DNS",
      "connect_timeout": "10s",
      "max_requests_per_connection": 10000,
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 10000,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-telemetry"
        }
       }
      },
      "protocol_selection": "USE_DOWNSTREAM_PROTOCOL",
      "load_assignment": {
       "cluster_name": "outbound|42422||istio-telemetry.istio-system.svc.cluster.local",
       "endpoints": [
        {
         "lb_endpoints": [
          {
           "endpoint": {
            "address": {
             "socket_address": {
              "address": "10.0.0.136",
              "port_value": 42422
             }
            }
           }
          }
         ]
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.732Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||istio-galley.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||istio-galley.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.676Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.678Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||istio-sidecar-injector.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||istio-sidecar-injector.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.706Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||istiod-istio-1611.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||istiod-istio-1611.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.735Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||kubernetes.default.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||kubernetes.default.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.675Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|443||metrics-server.kube-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|443||metrics-server.kube-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "protocol_selection": "USE_DOWNSTREAM_PROTOCOL"
     },
     "last_updated": "2021-09-24T17:02:36.734Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|53||kube-dns.kube-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|53||kube-dns.kube-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.733Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|8060||istio-citadel.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|8060||istio-citadel.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      }
     },
     "last_updated": "2021-09-24T17:02:36.676Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|8080||istio-pilot.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|8080||istio-pilot.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.679Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|80||default-http-backend.kube-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|80||default-http-backend.kube-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.733Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|80||istio-ingressgateway.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|80||istio-ingressgateway.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      }
     },
     "last_updated": "2021-09-24T17:02:36.677Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|8383||istio-operator.istio-operator.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|8383||istio-operator.istio-operator.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.676Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|853||istiod-istio-1611.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|853||istiod-istio-1611.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "protocol_selection": "USE_DOWNSTREAM_PROTOCOL"
     },
     "last_updated": "2021-09-24T17:02:36.735Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|9090||prometheus.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|9090||prometheus.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.736Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|9090||promsd.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|9090||promsd.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.733Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|9091||istio-policy.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-policy"
        }
       }
      }
     },
     "last_updated": "2021-09-24T17:02:36.680Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
      "type": "STRICT_DNS",
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      },
      "metadata": {
       "filter_metadata": {
        "istio": {
         "config": "/apis/networking/v1alpha3/namespaces/istio-system/destination-rule/istio-telemetry"
        }
       }
      },
      "load_assignment": {
       "cluster_name": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
       "endpoints": [
        {
         "lb_endpoints": [
          {
           "endpoint": {
            "address": {
             "socket_address": {
              "address": "10.0.0.136",
              "port_value": 9091
             }
            }
           }
          }
         ]
        }
       ]
      }
     },
     "last_updated": "2021-09-24T17:02:36.707Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "cluster": {
      "name": "outbound|9901||istio-galley.istio-system.svc.cluster.local",
      "type": "EDS",
      "eds_cluster_config": {
       "eds_config": {
        "ads": {}
       },
       "service_name": "outbound|9901||istio-galley.istio-system.svc.cluster.local"
      },
      "connect_timeout": "10s",
      "circuit_breakers": {
       "thresholds": [
        {
         "max_connections": 4294967295,
         "max_pending_requests": 4294967295,
         "max_requests": 4294967295,
         "max_retries": 4294967295
        }
       ]
      },
      "http2_protocol_options": {
       "max_concurrent_streams": 1073741824
      }
     },
     "last_updated": "2021-09-24T17:02:36.677Z"
    }
   ]
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.ListenersConfigDump",
   "version_info": "2021-09-24T16:05:58Z/13",
   "static_listeners": [
    {
     "listener": {
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15090
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "config": {
           "codec_type": "AUTO",
           "http_filters": {
            "name": "envoy.router"
           },
           "stat_prefix": "stats",
           "route_config": {
            "virtual_hosts": [
             {
              "domains": [
               "*"
              ],
              "name": "backend",
              "routes": [
               {
                "route": {
                 "cluster": "prometheus_stats"
                },
                "match": {
                 "prefix": "/stats/prometheus"
                }
               }
              ]
             }
            ]
           }
          }
         }
        ]
       }
      ]
     },
     "last_updated": "2021-09-24T17:02:36.007Z"
    }
   ],
   "dynamic_active_listeners": [
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.124.2.21_15020",
      "address": {
       "socket_address": {
        "address": "10.124.2.21",
        "port_value": 15020
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "inbound|15020|mgmt-15020|mgmtCluster",
           "cluster": "inbound|15020|mgmt-15020|mgmtCluster"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "INBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.852Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_31400",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 31400
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|31400||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|31400||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.866Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15029",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15029
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15029||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15029||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.888Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.0.1_443",
      "address": {
       "socket_address": {
        "address": "10.0.0.1",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "kubernetes.default.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "kubernetes"
             },
             "destination.service.namespace": {
              "string_value": "default"
             },
             "destination.service.uid": {
              "string_value": "istio://default/services/kubernetes"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||kubernetes.default.svc.cluster.local",
           "cluster": "outbound|443||kubernetes.default.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.898Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15031",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15031
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15031||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15031||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.905Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.8.40_443",
      "address": {
       "socket_address": {
        "address": "10.0.8.40",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istiod-istio-1611"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istiod-istio-1611"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||istiod-istio-1611.istio-system.svc.cluster.local",
           "cluster": "outbound|443||istiod-istio-1611.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.913Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.9.174_443",
      "address": {
       "socket_address": {
        "address": "10.0.9.174",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-galley.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-galley"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-galley"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||istio-galley.istio-system.svc.cluster.local",
           "cluster": "outbound|443||istio-galley.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.920Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15030",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15030
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15030||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15030||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.927Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15032",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15032
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15032||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15032||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.935Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.159_443",
      "address": {
       "socket_address": {
        "address": "10.0.2.159",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-sidecar-injector.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-sidecar-injector"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-sidecar-injector"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||istio-sidecar-injector.istio-system.svc.cluster.local",
           "cluster": "outbound|443||istio-sidecar-injector.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.942Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.14.235_15011",
      "address": {
       "socket_address": {
        "address": "10.0.14.235",
        "port_value": 15011
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-pilot.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-pilot"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-pilot"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15011||istio-pilot.istio-system.svc.cluster.local",
           "cluster": "outbound|15011||istio-pilot.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.951Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.8.40_15012",
      "address": {
       "socket_address": {
        "address": "10.0.8.40",
        "port_value": 15012
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istiod-istio-1611"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istiod-istio-1611"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15012||istiod-istio-1611.istio-system.svc.cluster.local",
           "cluster": "outbound|15012||istiod-istio-1611.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.959Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_443",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|443||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.966Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.0.10_53",
      "address": {
       "socket_address": {
        "address": "10.0.0.10",
        "port_value": 53
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "kube-dns.kube-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "kube-dns"
             },
             "destination.service.namespace": {
              "string_value": "kube-system"
             },
             "destination.service.uid": {
              "string_value": "istio://kube-system/services/kube-dns"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|53||kube-dns.kube-system.svc.cluster.local",
           "cluster": "outbound|53||kube-dns.kube-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.974Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15443",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15443||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15443||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.981Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.0.136_42422",
      "address": {
       "socket_address": {
        "address": "10.0.0.136",
        "port_value": 42422
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-telemetry.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-telemetry"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-telemetry"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|42422||istio-telemetry.istio-system.svc.cluster.local",
           "cluster": "outbound|42422||istio-telemetry.istio-system.svc.cluster.local"
          }
         }
        ]
       },
       {
        "filter_chain_match": {
         "application_protocols": [
          "http/1.0",
          "http/1.1",
          "h2"
         ]
        },
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_10.0.0.136_42422",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "istio-telemetry.istio-system.svc.cluster.local:42422"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters": [
       {
        "name": "envoy.listener.tls_inspector"
       },
       {
        "name": "envoy.listener.http_inspector"
       }
      ],
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:36.997Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_15010",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15010
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_15010",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "15010"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.006Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.2.124_15020",
      "address": {
       "socket_address": {
        "address": "10.0.2.124",
        "port_value": 15020
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istio-ingressgateway"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istio-ingressgateway"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|15020||istio-ingressgateway.istio-system.svc.cluster.local",
           "cluster": "outbound|15020||istio-ingressgateway.istio-system.svc.cluster.local"
          }
         }
        ]
       },
       {
        "filter_chain_match": {
         "application_protocols": [
          "http/1.0",
          "http/1.1",
          "h2"
         ]
        },
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_10.0.2.124_15020",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "istio-ingressgateway.istio-system.svc.cluster.local:15020"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters": [
       {
        "name": "envoy.listener.tls_inspector"
       },
       {
        "name": "envoy.listener.http_inspector"
       }
      ],
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.021Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_80",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 80
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_80",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "80"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.029Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_8383",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 8383
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_8383",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "8383"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.037Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_15004",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15004
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_15004",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "15004"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.045Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_9901",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 9901
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_9901",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "9901"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.053Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_8080",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 8080
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_8080",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "8080"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.071Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_9090",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 9090
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_9090",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "9090"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.086Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_9091",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 9091
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_9091",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "9091"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.097Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.13.91_443",
      "address": {
       "socket_address": {
        "address": "10.0.13.91",
        "port_value": 443
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "metrics-server.kube-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "metrics-server"
             },
             "destination.service.namespace": {
              "string_value": "kube-system"
             },
             "destination.service.uid": {
              "string_value": "istio://kube-system/services/metrics-server"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|443||metrics-server.kube-system.svc.cluster.local",
           "cluster": "outbound|443||metrics-server.kube-system.svc.cluster.local"
          }
         }
        ]
       },
       {
        "filter_chain_match": {
         "application_protocols": [
          "http/1.0",
          "http/1.1",
          "h2"
         ]
        },
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_10.0.13.91_443",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "metrics-server.kube-system.svc.cluster.local:443"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters": [
       {
        "name": "envoy.listener.tls_inspector"
       },
       {
        "name": "envoy.listener.http_inspector"
       }
      ],
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.115Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "10.0.8.40_853",
      "address": {
       "socket_address": {
        "address": "10.0.8.40",
        "port_value": 853
       }
      },
      "filter_chains": [
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
             },
             "destination.service.name": {
              "string_value": "istiod-istio-1611"
             },
             "destination.service.namespace": {
              "string_value": "istio-system"
             },
             "destination.service.uid": {
              "string_value": "istio://istio-system/services/istiod-istio-1611"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "outbound|853||istiod-istio-1611.istio-system.svc.cluster.local",
           "cluster": "outbound|853||istiod-istio-1611.istio-system.svc.cluster.local"
          }
         }
        ]
       },
       {
        "filter_chain_match": {
         "application_protocols": [
          "http/1.0",
          "http/1.1",
          "h2"
         ]
        },
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_10.0.8.40_853",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "istiod-istio-1611.istio-system.svc.cluster.local:853"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters": [
       {
        "name": "envoy.listener.tls_inspector"
       },
       {
        "name": "envoy.listener.http_inspector"
       }
      ],
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.132Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_8060",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 8060
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_8060",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "8060"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.140Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "0.0.0.0_15014",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15014
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "envoy.http_connection_manager",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager",
           "stat_prefix": "outbound_0.0.0.0_15014",
           "http_filters": [
            {
             "name": "mixer",
             "typed_config": {
              "@type": "type.googleapis.com/istio.mixer.v1.config.client.HttpClientConfig",
              "transport": {
               "network_fail_policy": {
                "policy": "FAIL_CLOSE",
                "base_retry_wait": "0.080s",
                "max_retry_wait": "1s"
               },
               "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
               "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
               "report_batch_max_entries": 100,
               "report_batch_max_time": "1s"
              },
              "service_configs": {
               "default": {
                "disable_check_calls": true
               }
              },
              "default_destination_service": "default",
              "mixer_attributes": {
               "attributes": {
                "context.proxy_version": {
                 "string_value": "1.4.10"
                },
                "context.reporter.kind": {
                 "string_value": "outbound"
                },
                "context.reporter.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                },
                "source.namespace": {
                 "string_value": "default"
                },
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              },
              "forward_attributes": {
               "attributes": {
                "source.uid": {
                 "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
                }
               }
              }
             }
            },
            {
             "name": "envoy.cors"
            },
            {
             "name": "envoy.fault"
            },
            {
             "name": "envoy.router"
            }
           ],
           "tracing": {
            "operation_name": "EGRESS",
            "client_sampling": {
             "value": 100
            },
            "random_sampling": {
             "value": 1
            },
            "overall_sampling": {
             "value": 100
            }
           },
           "use_remote_address": false,
           "generate_request_id": true,
           "upgrade_configs": [
            {
             "upgrade_type": "websocket"
            }
           ],
           "stream_idle_timeout": "0s",
           "normalize_path": true,
           "rds": {
            "config_source": {
             "ads": {}
            },
            "route_config_name": "15014"
           }
          }
         }
        ]
       }
      ],
      "deprecated_v1": {
       "bind_to_port": false
      },
      "listener_filters_timeout": "5s",
      "traffic_direction": "OUTBOUND",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.149Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "virtualOutbound",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15001
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "BlackHoleCluster",
           "cluster": "BlackHoleCluster"
          }
         }
        ]
       },
       {
        "filters": [
         {
          "name": "mixer",
          "typed_config": {
           "@type": "type.googleapis.com/istio.mixer.v1.config.client.TcpClientConfig",
           "transport": {
            "network_fail_policy": {
             "policy": "FAIL_CLOSE",
             "base_retry_wait": "0.080s",
             "max_retry_wait": "1s"
            },
            "check_cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
            "report_cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
            "report_batch_max_entries": 100,
            "report_batch_max_time": "1s"
           },
           "mixer_attributes": {
            "attributes": {
             "context.proxy_version": {
              "string_value": "1.4.10"
             },
             "context.reporter.kind": {
              "string_value": "outbound"
             },
             "context.reporter.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             },
             "destination.service.host": {
              "string_value": "PassthroughCluster"
             },
             "destination.service.name": {
              "string_value": "PassthroughCluster"
             },
             "source.namespace": {
              "string_value": "default"
             },
             "source.uid": {
              "string_value": "kubernetes://httpbin1-6c97fbc967-lz4fb.default"
             }
            }
           },
           "disable_check_calls": true
          }
         },
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "PassthroughCluster",
           "cluster": "PassthroughCluster"
          }
         }
        ]
       }
      ],
      "use_original_dst": true
     },
     "last_updated": "2021-09-24T17:02:37.157Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "listener": {
      "name": "virtualInbound",
      "address": {
       "socket_address": {
        "address": "0.0.0.0",
        "port_value": 15006
       }
      },
      "filter_chains": [
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "0.0.0.0",
           "prefix_len": 0
          }
         ]
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "InboundPassthroughClusterIpv4",
           "cluster": "InboundPassthroughClusterIpv4"
          }
         }
        ],
        "metadata": {
         "filter_metadata": {
          "pilot_meta": {
           "original_listener_name": "virtualInbound"
          }
         }
        }
       },
       {
        "filter_chain_match": {
         "prefix_ranges": [
          {
           "address_prefix": "10.124.2.21",
           "prefix_len": 32
          }
         ],
         "destination_port": 15020
        },
        "filters": [
         {
          "name": "envoy.tcp_proxy",
          "typed_config": {
           "@type": "type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy",
           "stat_prefix": "inbound|15020|mgmt-15020|mgmtCluster",
           "cluster": "inbound|15020|mgmt-15020|mgmtCluster"
          }
         }
        ],
        "metadata": {
         "filter_metadata": {
          "pilot_meta": {
           "original_listener_name": "10.124.2.21_15020"
          }
         }
        }
       }
      ],
      "listener_filters": [
       {
        "name": "envoy.listener.original_dst"
       }
      ],
      "listener_filters_timeout": "1s",
      "continue_on_listener_filters_timeout": true
     },
     "last_updated": "2021-09-24T17:02:37.158Z"
    }
   ]
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.ScopedRoutesConfigDump"
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.RoutesConfigDump",
   "static_route_configs": [
    {
     "route_config": {
      "virtual_hosts": [
       {
        "name": "backend",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/stats/prometheus"
          },
          "route": {
           "cluster": "prometheus_stats"
          }
         }
        ]
       }
      ]
     },
     "last_updated": "2021-09-24T17:02:36.006Z"
    }
   ],
   "dynamic_route_configs": [
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "15014",
      "virtual_hosts": [
       {
        "name": "istio-citadel.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-citadel.istio-system.svc.cluster.local",
         "istio-citadel.istio-system.svc.cluster.local:15014",
         "istio-citadel.istio-system",
         "istio-citadel.istio-system:15014",
         "istio-citadel.istio-system.svc.cluster",
         "istio-citadel.istio-system.svc.cluster:15014",
         "istio-citadel.istio-system.svc",
         "istio-citadel.istio-system.svc:15014",
         "10.0.6.162",
         "10.0.6.162:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-citadel.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-citadel.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-citadel.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-citadel"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-citadel"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-citadel.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-citadel"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-citadel"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-galley.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-galley.istio-system.svc.cluster.local",
         "istio-galley.istio-system.svc.cluster.local:15014",
         "istio-galley.istio-system",
         "istio-galley.istio-system:15014",
         "istio-galley.istio-system.svc.cluster",
         "istio-galley.istio-system.svc.cluster:15014",
         "istio-galley.istio-system.svc",
         "istio-galley.istio-system.svc:15014",
         "10.0.9.174",
         "10.0.9.174:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-galley.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-galley.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-galley.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-galley"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-galley"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-galley.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-galley"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-galley"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-pilot.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-pilot.istio-system.svc.cluster.local",
         "istio-pilot.istio-system.svc.cluster.local:15014",
         "istio-pilot.istio-system",
         "istio-pilot.istio-system:15014",
         "istio-pilot.istio-system.svc.cluster",
         "istio-pilot.istio-system.svc.cluster:15014",
         "istio-pilot.istio-system.svc",
         "istio-pilot.istio-system.svc:15014",
         "10.0.14.235",
         "10.0.14.235:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-pilot.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-pilot.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-policy.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-policy.istio-system.svc.cluster.local",
         "istio-policy.istio-system.svc.cluster.local:15014",
         "istio-policy.istio-system",
         "istio-policy.istio-system:15014",
         "istio-policy.istio-system.svc.cluster",
         "istio-policy.istio-system.svc.cluster:15014",
         "istio-policy.istio-system.svc",
         "istio-policy.istio-system.svc:15014",
         "10.0.9.230",
         "10.0.9.230:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-policy.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-policy.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-sidecar-injector.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-sidecar-injector.istio-system.svc.cluster.local",
         "istio-sidecar-injector.istio-system.svc.cluster.local:15014",
         "istio-sidecar-injector.istio-system",
         "istio-sidecar-injector.istio-system:15014",
         "istio-sidecar-injector.istio-system.svc.cluster",
         "istio-sidecar-injector.istio-system.svc.cluster:15014",
         "istio-sidecar-injector.istio-system.svc",
         "istio-sidecar-injector.istio-system.svc:15014",
         "10.0.2.159",
         "10.0.2.159:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-sidecar-injector.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-sidecar-injector.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-sidecar-injector.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-sidecar-injector"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-sidecar-injector"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-sidecar-injector.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-sidecar-injector"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-sidecar-injector"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-telemetry.istio-system.svc.cluster.local:15014",
        "domains": [
         "istio-telemetry.istio-system.svc.cluster.local",
         "istio-telemetry.istio-system.svc.cluster.local:15014",
         "istio-telemetry.istio-system",
         "istio-telemetry.istio-system:15014",
         "istio-telemetry.istio-system.svc.cluster",
         "istio-telemetry.istio-system.svc.cluster:15014",
         "istio-telemetry.istio-system.svc",
         "istio-telemetry.istio-system.svc:15014",
         "10.0.0.136",
         "10.0.0.136:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istio-telemetry.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-telemetry.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istiod-istio-1611.istio-system.svc.cluster.local:15014",
        "domains": [
         "istiod-istio-1611.istio-system.svc.cluster.local",
         "istiod-istio-1611.istio-system.svc.cluster.local:15014",
         "istiod-istio-1611.istio-system",
         "istiod-istio-1611.istio-system:15014",
         "istiod-istio-1611.istio-system.svc.cluster",
         "istiod-istio-1611.istio-system.svc.cluster:15014",
         "istiod-istio-1611.istio-system.svc",
         "istiod-istio-1611.istio-system.svc:15014",
         "10.0.8.40",
         "10.0.8.40:15014"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15014||istiod-istio-1611.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istiod-istio-1611.istio-system.svc.cluster.local:15014/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.183Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "8060",
      "virtual_hosts": [
       {
        "name": "istio-citadel.istio-system.svc.cluster.local:8060",
        "domains": [
         "istio-citadel.istio-system.svc.cluster.local",
         "istio-citadel.istio-system.svc.cluster.local:8060",
         "istio-citadel.istio-system",
         "istio-citadel.istio-system:8060",
         "istio-citadel.istio-system.svc.cluster",
         "istio-citadel.istio-system.svc.cluster:8060",
         "istio-citadel.istio-system.svc",
         "istio-citadel.istio-system.svc:8060",
         "10.0.6.162",
         "10.0.6.162:8060"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|8060||istio-citadel.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-citadel.istio-system.svc.cluster.local:8060/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-citadel.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-citadel"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-citadel"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-citadel.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-citadel"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-citadel"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.184Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "istiod-istio-1611.istio-system.svc.cluster.local:853",
      "virtual_hosts": [
       {
        "name": "istiod-istio-1611.istio-system.svc.cluster.local:853",
        "domains": [
         "istiod-istio-1611.istio-system.svc.cluster.local",
         "istiod-istio-1611.istio-system.svc.cluster.local:853",
         "istiod-istio-1611.istio-system",
         "istiod-istio-1611.istio-system:853",
         "istiod-istio-1611.istio-system.svc.cluster",
         "istiod-istio-1611.istio-system.svc.cluster:853",
         "istiod-istio-1611.istio-system.svc",
         "istiod-istio-1611.istio-system.svc:853",
         "10.0.8.40",
         "10.0.8.40:853"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|853||istiod-istio-1611.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istiod-istio-1611.istio-system.svc.cluster.local:853/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.184Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "9090",
      "virtual_hosts": [
       {
        "name": "prometheus.istio-system.svc.cluster.local:9090",
        "domains": [
         "prometheus.istio-system.svc.cluster.local",
         "prometheus.istio-system.svc.cluster.local:9090",
         "prometheus.istio-system",
         "prometheus.istio-system:9090",
         "prometheus.istio-system.svc.cluster",
         "prometheus.istio-system.svc.cluster:9090",
         "prometheus.istio-system.svc",
         "prometheus.istio-system.svc:9090",
         "10.0.3.164",
         "10.0.3.164:9090"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|9090||prometheus.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "prometheus.istio-system.svc.cluster.local:9090/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "prometheus.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "prometheus"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/prometheus"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "prometheus.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "prometheus"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/prometheus"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "promsd.istio-system.svc.cluster.local:9090",
        "domains": [
         "promsd.istio-system.svc.cluster.local",
         "promsd.istio-system.svc.cluster.local:9090",
         "promsd.istio-system",
         "promsd.istio-system:9090",
         "promsd.istio-system.svc.cluster",
         "promsd.istio-system.svc.cluster:9090",
         "promsd.istio-system.svc",
         "promsd.istio-system.svc:9090",
         "10.0.3.212",
         "10.0.3.212:9090"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|9090||promsd.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "promsd.istio-system.svc.cluster.local:9090/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "promsd.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "promsd"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/promsd"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "promsd.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "promsd"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/promsd"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.186Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "8080",
      "virtual_hosts": [
       {
        "name": "istio-pilot.istio-system.svc.cluster.local:8080",
        "domains": [
         "istio-pilot.istio-system.svc.cluster.local",
         "istio-pilot.istio-system.svc.cluster.local:8080",
         "istio-pilot.istio-system",
         "istio-pilot.istio-system:8080",
         "istio-pilot.istio-system.svc.cluster",
         "istio-pilot.istio-system.svc.cluster:8080",
         "istio-pilot.istio-system.svc",
         "istio-pilot.istio-system.svc:8080",
         "10.0.14.235",
         "10.0.14.235:8080"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|8080||istio-pilot.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-pilot.istio-system.svc.cluster.local:8080/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.186Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "9901",
      "virtual_hosts": [
       {
        "name": "istio-galley.istio-system.svc.cluster.local:9901",
        "domains": [
         "istio-galley.istio-system.svc.cluster.local",
         "istio-galley.istio-system.svc.cluster.local:9901",
         "istio-galley.istio-system",
         "istio-galley.istio-system:9901",
         "istio-galley.istio-system.svc.cluster",
         "istio-galley.istio-system.svc.cluster:9901",
         "istio-galley.istio-system.svc",
         "istio-galley.istio-system.svc:9901",
         "10.0.9.174",
         "10.0.9.174:9901"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|9901||istio-galley.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-galley.istio-system.svc.cluster.local:9901/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-galley.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-galley"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-galley"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-galley.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-galley"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-galley"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.186Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "15004",
      "virtual_hosts": [
       {
        "name": "istio-policy.istio-system.svc.cluster.local:15004",
        "domains": [
         "istio-policy.istio-system.svc.cluster.local",
         "istio-policy.istio-system.svc.cluster.local:15004",
         "istio-policy.istio-system",
         "istio-policy.istio-system:15004",
         "istio-policy.istio-system.svc.cluster",
         "istio-policy.istio-system.svc.cluster:15004",
         "istio-policy.istio-system.svc",
         "istio-policy.istio-system.svc:15004",
         "10.0.9.230",
         "10.0.9.230:15004"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15004||istio-policy.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-policy.istio-system.svc.cluster.local:15004/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-telemetry.istio-system.svc.cluster.local:15004",
        "domains": [
         "istio-telemetry.istio-system.svc.cluster.local",
         "istio-telemetry.istio-system.svc.cluster.local:15004",
         "istio-telemetry.istio-system",
         "istio-telemetry.istio-system:15004",
         "istio-telemetry.istio-system.svc.cluster",
         "istio-telemetry.istio-system.svc.cluster:15004",
         "istio-telemetry.istio-system.svc",
         "istio-telemetry.istio-system.svc:15004",
         "10.0.0.136",
         "10.0.0.136:15004"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15004||istio-telemetry.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-telemetry.istio-system.svc.cluster.local:15004/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.187Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "metrics-server.kube-system.svc.cluster.local:443",
      "virtual_hosts": [
       {
        "name": "metrics-server.kube-system.svc.cluster.local:443",
        "domains": [
         "metrics-server.kube-system.svc.cluster.local",
         "metrics-server.kube-system.svc.cluster.local:443",
         "metrics-server.kube-system",
         "metrics-server.kube-system:443",
         "metrics-server.kube-system.svc.cluster",
         "metrics-server.kube-system.svc.cluster:443",
         "metrics-server.kube-system.svc",
         "metrics-server.kube-system.svc:443",
         "10.0.13.91",
         "10.0.13.91:443"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|443||metrics-server.kube-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "metrics-server.kube-system.svc.cluster.local:443/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "metrics-server.kube-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "metrics-server"
              },
              "destination.service.namespace": {
               "string_value": "kube-system"
              },
              "destination.service.uid": {
               "string_value": "istio://kube-system/services/metrics-server"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "metrics-server.kube-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "metrics-server"
              },
              "destination.service.namespace": {
               "string_value": "kube-system"
              },
              "destination.service.uid": {
               "string_value": "istio://kube-system/services/metrics-server"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.184Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "istio-ingressgateway.istio-system.svc.cluster.local:15020",
      "virtual_hosts": [
       {
        "name": "istio-ingressgateway.istio-system.svc.cluster.local:15020",
        "domains": [
         "istio-ingressgateway.istio-system.svc.cluster.local",
         "istio-ingressgateway.istio-system.svc.cluster.local:15020",
         "istio-ingressgateway.istio-system",
         "istio-ingressgateway.istio-system:15020",
         "istio-ingressgateway.istio-system.svc.cluster",
         "istio-ingressgateway.istio-system.svc.cluster:15020",
         "istio-ingressgateway.istio-system.svc",
         "istio-ingressgateway.istio-system.svc:15020",
         "10.0.2.124",
         "10.0.2.124:15020"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15020||istio-ingressgateway.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-ingressgateway.istio-system.svc.cluster.local:15020/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-ingressgateway"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-ingressgateway"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-ingressgateway"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-ingressgateway"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.188Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "80",
      "virtual_hosts": [
       {
        "name": "default-http-backend.kube-system.svc.cluster.local:80",
        "domains": [
         "default-http-backend.kube-system.svc.cluster.local",
         "default-http-backend.kube-system.svc.cluster.local:80",
         "default-http-backend.kube-system",
         "default-http-backend.kube-system:80",
         "default-http-backend.kube-system.svc.cluster",
         "default-http-backend.kube-system.svc.cluster:80",
         "default-http-backend.kube-system.svc",
         "default-http-backend.kube-system.svc:80",
         "10.0.11.198",
         "10.0.11.198:80"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|80||default-http-backend.kube-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "default-http-backend.kube-system.svc.cluster.local:80/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "default-http-backend.kube-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "default-http-backend"
              },
              "destination.service.namespace": {
               "string_value": "kube-system"
              },
              "destination.service.uid": {
               "string_value": "istio://kube-system/services/default-http-backend"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "default-http-backend.kube-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "default-http-backend"
              },
              "destination.service.namespace": {
               "string_value": "kube-system"
              },
              "destination.service.uid": {
               "string_value": "istio://kube-system/services/default-http-backend"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-ingressgateway.istio-system.svc.cluster.local:80",
        "domains": [
         "istio-ingressgateway.istio-system.svc.cluster.local",
         "istio-ingressgateway.istio-system.svc.cluster.local:80",
         "istio-ingressgateway.istio-system",
         "istio-ingressgateway.istio-system:80",
         "istio-ingressgateway.istio-system.svc.cluster",
         "istio-ingressgateway.istio-system.svc.cluster:80",
         "istio-ingressgateway.istio-system.svc",
         "istio-ingressgateway.istio-system.svc:80",
         "10.0.2.124",
         "10.0.2.124:80"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|80||istio-ingressgateway.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-ingressgateway.istio-system.svc.cluster.local:80/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-ingressgateway"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-ingressgateway"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-ingressgateway.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-ingressgateway"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-ingressgateway"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.188Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "15010",
      "virtual_hosts": [
       {
        "name": "istio-pilot.istio-system.svc.cluster.local:15010",
        "domains": [
         "istio-pilot.istio-system.svc.cluster.local",
         "istio-pilot.istio-system.svc.cluster.local:15010",
         "istio-pilot.istio-system",
         "istio-pilot.istio-system:15010",
         "istio-pilot.istio-system.svc.cluster",
         "istio-pilot.istio-system.svc.cluster:15010",
         "istio-pilot.istio-system.svc",
         "istio-pilot.istio-system.svc:15010",
         "10.0.14.235",
         "10.0.14.235:15010"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15010||istio-pilot.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-pilot.istio-system.svc.cluster.local:15010/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-pilot.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-pilot"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-pilot"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istiod-istio-1611.istio-system.svc.cluster.local:15010",
        "domains": [
         "istiod-istio-1611.istio-system.svc.cluster.local",
         "istiod-istio-1611.istio-system.svc.cluster.local:15010",
         "istiod-istio-1611.istio-system",
         "istiod-istio-1611.istio-system:15010",
         "istiod-istio-1611.istio-system.svc.cluster",
         "istiod-istio-1611.istio-system.svc.cluster:15010",
         "istiod-istio-1611.istio-system.svc",
         "istiod-istio-1611.istio-system.svc:15010",
         "10.0.8.40",
         "10.0.8.40:15010"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|15010||istiod-istio-1611.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istiod-istio-1611.istio-system.svc.cluster.local:15010/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istiod-istio-1611.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istiod-istio-1611"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istiod-istio-1611"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.188Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "9091",
      "virtual_hosts": [
       {
        "name": "istio-policy.istio-system.svc.cluster.local:9091",
        "domains": [
         "istio-policy.istio-system.svc.cluster.local",
         "istio-policy.istio-system.svc.cluster.local:9091",
         "istio-policy.istio-system",
         "istio-policy.istio-system:9091",
         "istio-policy.istio-system.svc.cluster",
         "istio-policy.istio-system.svc.cluster:9091",
         "istio-policy.istio-system.svc",
         "istio-policy.istio-system.svc:9091",
         "10.0.9.230",
         "10.0.9.230:9091"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|9091||istio-policy.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-policy.istio-system.svc.cluster.local:9091/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-policy.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-policy"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-policy"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "istio-telemetry.istio-system.svc.cluster.local:9091",
        "domains": [
         "istio-telemetry.istio-system.svc.cluster.local",
         "istio-telemetry.istio-system.svc.cluster.local:9091",
         "istio-telemetry.istio-system",
         "istio-telemetry.istio-system:9091",
         "istio-telemetry.istio-system.svc.cluster",
         "istio-telemetry.istio-system.svc.cluster:9091",
         "istio-telemetry.istio-system.svc",
         "istio-telemetry.istio-system.svc:9091",
         "10.0.0.136",
         "10.0.0.136:9091"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|9091||istio-telemetry.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-telemetry.istio-system.svc.cluster.local:9091/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.185Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "8383",
      "virtual_hosts": [
       {
        "name": "istio-operator.istio-operator.svc.cluster.local:8383",
        "domains": [
         "istio-operator.istio-operator.svc.cluster.local",
         "istio-operator.istio-operator.svc.cluster.local:8383",
         "istio-operator.istio-operator",
         "istio-operator.istio-operator:8383",
         "istio-operator.istio-operator.svc.cluster",
         "istio-operator.istio-operator.svc.cluster:8383",
         "istio-operator.istio-operator.svc",
         "istio-operator.istio-operator.svc:8383",
         "10.0.11.12",
         "10.0.11.12:8383"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|8383||istio-operator.istio-operator.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-operator.istio-operator.svc.cluster.local:8383/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-operator.istio-operator.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-operator"
              },
              "destination.service.namespace": {
               "string_value": "istio-operator"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-operator/services/istio-operator"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-operator.istio-operator.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-operator"
              },
              "destination.service.namespace": {
               "string_value": "istio-operator"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-operator/services/istio-operator"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       },
       {
        "name": "allow_any",
        "domains": [
         "*"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "PassthroughCluster"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.name": {
               "string_value": "PassthroughCluster"
              }
             }
            }
           }
          }
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.187Z"
    },
    {
     "version_info": "2021-09-24T16:05:58Z/13",
     "route_config": {
      "name": "istio-telemetry.istio-system.svc.cluster.local:42422",
      "virtual_hosts": [
       {
        "name": "istio-telemetry.istio-system.svc.cluster.local:42422",
        "domains": [
         "istio-telemetry.istio-system.svc.cluster.local",
         "istio-telemetry.istio-system.svc.cluster.local:42422",
         "istio-telemetry.istio-system",
         "istio-telemetry.istio-system:42422",
         "istio-telemetry.istio-system.svc.cluster",
         "istio-telemetry.istio-system.svc.cluster:42422",
         "istio-telemetry.istio-system.svc",
         "istio-telemetry.istio-system.svc:42422",
         "10.0.0.136",
         "10.0.0.136:42422"
        ],
        "routes": [
         {
          "match": {
           "prefix": "/"
          },
          "route": {
           "cluster": "outbound|42422||istio-telemetry.istio-system.svc.cluster.local",
           "timeout": "0s",
           "retry_policy": {
            "retry_on": "connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes",
            "num_retries": 2,
            "retry_host_predicate": [
             {
              "name": "envoy.retry_host_predicates.previous_hosts"
             }
            ],
            "host_selection_retry_max_attempts": "5",
            "retriable_status_codes": [
             503
            ]
           },
           "max_grpc_timeout": "0s"
          },
          "decorator": {
           "operation": "istio-telemetry.istio-system.svc.cluster.local:42422/*"
          },
          "typed_per_filter_config": {
           "mixer": {
            "@type": "type.googleapis.com/istio.mixer.v1.config.client.ServiceConfig",
            "disable_check_calls": true,
            "mixer_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            },
            "forward_attributes": {
             "attributes": {
              "destination.service.host": {
               "string_value": "istio-telemetry.istio-system.svc.cluster.local"
              },
              "destination.service.name": {
               "string_value": "istio-telemetry"
              },
              "destination.service.namespace": {
               "string_value": "istio-system"
              },
              "destination.service.uid": {
               "string_value": "istio://istio-system/services/istio-telemetry"
              }
             }
            }
           }
          },
          "name": "default"
         }
        ]
       }
      ],
      "validate_clusters": false
     },
     "last_updated": "2021-09-24T17:02:37.189Z"
    }
   ]
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v2alpha.SecretsConfigDump"
  }
 ]
}`

const testJson = `
{
  "configs": [
    {
      "@type": "type.googleapis.com/envoy.admin.v3.BootstrapConfigDump",
      "bootstrap": {
        "node": {
          "cluster": "httpbin2.default",
          "id": "sidecar~10.4.2.16~httpbin2-5cb469675d-x2m7g.default~default.svc.cluster.local",
          "metadata": {
            "CLUSTER_ID": "Kubernetes",
            "INSTANCE_IPS": "10.4.2.16",
            "ISTIO_PROXY_SHA": "istio-proxy:4b528a87271e841bd64daf841a1a384ed4fcac68",
            "LABELS": {
              "app": "httpbin2",
              "istio.io/rev": "default",
              "pod-template-hash": "5cb469675d",
              "security.istio.io/tlsMode": "istio",
              "service.istio.io/canonical-name": "httpbin2",
              "service.istio.io/canonical-revision": "v1",
              "version": "v1"
            },
            "MESH_ID": "cluster.local",
            "NAME": "httpbin2-5cb469675d-x2m7g",
            "PLATFORM_METADATA": {
              "gcp_gce_instance_id": "7955581888538677628",
              "gcp_gke_cluster_name": "cluster-1",
              "gcp_gke_cluster_url": "https://container.googleapis.com/v1/projects/service-mithril/locations/us-central1-c/clusters/cluster-1",
              "gcp_location": "us-central1-c",
              "gcp_project": "service-mithril",
              "gcp_project_number": "199906641779"
            },
            "POD_PORTS": "[{\"containerPort\":80,\"protocol\":\"TCP\"}]",
            "PROXY_CONFIG": {
              "binaryPath": "/usr/local/bin/envoy",
              "concurrency": 2,
              "configPath": "./etc/istio/proxy",
              "controlPlaneAuthPolicy": "MUTUAL_TLS",
              "discoveryAddress": "istiod.istio-system.svc:15012",
              "drainDuration": "45s",
              "parentShutdownDuration": "60s",
              "proxyAdminPort": 15000,
              "serviceCluster": "httpbin2.default",
              "statNameLength": 189,
              "statusPort": 15020,
              "terminationDrainDuration": "5s",
              "tracing": {
                "zipkin": {
                  "address": "zipkin.istio-system:9411"
                }
              }
            },
            "SERVICE_ACCOUNT": "httpbin2",
            "kubectl.kubernetes.io/restartedAt": "2021-09-21T17:10:50Z"
          }
        }
      }
    }
  ]
}
`

func TestRetrieveEnvoyConfig(t *testing.T) {
	istioJson := []byte(testJson)
	istioGkeJson := []byte(istioGkeDumpJson)
	// Removing the URL check for now since you need an envoy instance hosted to test this
	//ec, _ := envoy.RetrieveConfig("http://localhost:15000/config_dump")
	ec, _ := envoy.LoadConfig(istioJson)
	address, _ := ec.DiscoveryAddress()
	assert.Equal(t, "istiod.istio-system.svc:15012", address)
	ec, _ = envoy.LoadConfig(istioGkeJson)
	address, _ = ec.DiscoveryAddress()
	assert.Equal(t, "istio-pilot.istio-system.svc.cluster.local:15010", address)

}
