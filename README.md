# Mesh Hunter - A service mesh scanning tool

[![CI](https://github.com/praetorian-inc/mesh-hunter/workflows/CI/badge.svg)](actions?query=branch%3Adevelopment)
[![Release](https://github.com/praetorian-inc/mesh-hunter/workflows/Release/badge.svg)](releases)

Mesh Hunter gathers and analyzes the configuration of an Istio cluster and audits it for potential violations of security best practices.

## Why We Built Mesh Hunter

Like all cloud infrastructure, Istio requires some hardening effort beyond what a default deployment offers.
The [Istio Security Best Practices](https://istio.io/latest/docs/ops/best-practices/security/) document
covers this in great detail. This hardening process has a lot of moving parts and it's easy to miss
one of the steps that could assist an attacker in compromising a cluster. Mesh Hunter was built to make
the detection of these missing hardening steps as straightforward as possible.

The two usage modes can help engineers analyze their clusters from different perspectives:

* The perspective of an attacker that has just obtained code execution on an Istio workload but without any other context or permissions.
* The perspective of a systems engineer that has the ability to dump all relevant configuration information for analysis.

By implementing analysis methods for both of these perspectives, Mesh Hunter is able to gather a more "complete"
picture of the security posture of an Istio cluster.

For more information, please read [our blog post](https://www.praetorian.com/blog/wherever-this-will-live/).

## Install

You can install Mesh Hunter locally by using any one of the options listed below.

### Install with `go install`

```shell
$ go install github.com/praetorian-inc/mesh-hunter@latest
```

### Install a release binary

1. Download the binary for your OS from the [releases page](https://github.com/praetorian-inc/mesh-hunter/releases).

2. (OPTIONAL) Download the `checksums.txt` file to verify the integrity of the archive

```shell
# Check the checksum of the downloaded archive
$ shasum -a 256 mesh-hunter_${VERSION}_${ARCH}.tar.gz
b05c4d7895be260aa16336f29249c50b84897dab90e1221c9e96af9233751f22  mesh-hunter_${VERSION}_${ARCH}.tar.gz

$ cat mesh-hunter_${VERSION}_${ARCH}_checksums.txt | grep mithril_${VERSION}_${ARCH}.tar.gz
b05c4d7895be260aa16336f29249c50b84897dab90e1221c9e96af9233751f22  mesh-hunter_${VERSION}_${ARCH}.tar.gz
```

3. Extract the downloaded archive

```shell
$ tar -xvf mesh-hunter_${VERSION}_${ARCH}.tar.gz
```

4. Move the `mesh-hunter` binary into your path:

```shell
$ mv ./mesh-hunter /usr/local/bin/
```

### Clone and build yourself

```shell
# clone the Mesh Hunter repo
$ git clone https://github.com/praetorian-inc/mesh-hunter.git

# navigate into the repo directory and build
$ cd mesh-hunter
$ go build

# Move the Mesh Hunter binary into your path
$ mv ./mesh-hunter /usr/local/bin
```

## Usage

There are two main modes of operation for Mesh Hunter. With no positional argument,
Mesh Hunter will assume it is running inside of a cluster enabled with Istio, and
begin to enumerate the required data. Optionally, you can point mithril at a
directory containing Kubernets YAML files.

### Run Mesh Hunter against static configuration information

```shell
# running with a directory specified will cause it to run in file analysis mode
./mesh-hunter [options] <directory name>
```

### Run Mesh Hunter in an Istio workload container

```shell
./mesh-hunter [options]
```

### Get Help

```shell
mesh-hunter help
```

### Command Line Options

Mesh Hunter comes equipped with several command line options to influence the
operation of the tool. Additionally, many configuration options can be passed
to the tool through a configuration file. By default, Mesh Hunter looks for the
config file at `./mesh-hunter.yml` (the directory from which the tool is run), but
can be passed as a switch to specify an arbitrary file location.

Configuration of Mesh Hunter is handled by a combination of
[Cobra](https://github.com/spf13/cobra "Cobra") and
[Viper](https://github.com/spf13/viper "Viper"). This allows Mesh Hunter to be
configured through the following methods, in order of precedence.

1. Command Line Flag
2. Environment Variables
3. Configuration File

It should be noted that any data that is discovered during a run will overwrite
all configuration options.

The following configuration options can be specified:

* `-c <file>` `--config <file>` - the configuration file location (default:
  `./mesh-hunter.yml`)

* `-l <level>` `--log-level <level>` - log level for console output, because
  logging is handled by [Logrus](https://github.com/sirupsen/logrus "Logrus"),
  the currently supported levels are trace, debug, info, warning, error, fatal,
  and panic. (default: `info`)

* `-s` `--save-config` - if this switch is passed, the configuration of Mesh Hunter
  will be written out to the specified config file. This is useful if the tool
  is to be run multiple times on the same cluster to allow for fewer arguments
  to be passed in subsequent runs. NOTE: this will overwrite the existing config
  file every time.

* `--format [text|json]` - the output format for the tool, this is either `text`
  for human readable content, or `json` for structured output.

* `--export <directory>` - this flag will cause Mesh Hunter to output the discovered
  Kubernetes resources to a directory as YAML files

* `--istio-version <version>` - if the Istio control plane version is known prior
  to running the tool, it can be passed via this flag. Additionally, it binds to
  the configuration variable `istio-version` in the configuration file.

* `--istio-namespace <namespace>` - if the namespace running the Istio control
  plane is known prior to running the tool, it can be passed via this flag.
  Additionally, it binds to the configuration variable `istio-namespace` in the
  configuration file.

* `--discovery-address <ip:port>` - this specifies the address of the
  unauthenticated XDS port. It is bound to the configuration variable
  `discovery-address`.

* `--debugz-address <ip:port>` - this specifies the address of the Istiod's debug
  API. It is bound to the configuration variable `debugz-address`.

* `--kubelet-addresses <list of ip:port>` - this specifies a list of kubelet nodes
  read-only API ports. It is bound to the configuration variable
  `kubelet-addresses`

To set these flags with environment variables, simply uppercase the
configuration variable name, and replace dashes with underscores, for example:
`istio-version` -> `ISTIO_VERSION`
