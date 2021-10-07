# Mithril (Real name TBD) - A service mesh scanning tool

[![CI](https://github.com/praetorian-inc/mithril/workflows/CI/badge.svg)](actions?query=branch%3Adevelopment)
[![Release](https://github.com/praetorian-inc/mithril/workflows/Release/badge.svg)](releases)

What does this do

## Why We Built Mithril

Like all cloud infrastructure, Istio requires some hardening effort beyond what a default deployment offers.
The [Istio Security Best Practices](https://istio.io/latest/docs/ops/best-practices/security/) document
covers this in great detail. This hardening process has a lot of moving parts and it's easy to miss
one of the steps that could assist an attacker in compromising a cluster. Mithril was built to make 
the detection of these missing hardening steps as straightforward as possible.

The two usage modes can help engineers analyze their clusters from different perspectives:

* The perspective of an attacker that has just obtained code execution on an Istio workload but without any other context or permissions.
* The perspective of a systems engineer that has the ability to dump all relevant configuration information for analysis.

By implementing analysis methods for both of these perspectives, Mithril is able to gather a more "complete"
picture of the security posture of an Istio cluster.

For more information, please read [our blog post](https://www.praetorian.com/blog/wherever-this-will-live/).

## Install

You can install Mithril locally by using any one of the options listed below.

### Install with `go install`

```shell
$ go install github.com/praetorian-inc/mithril@latest
```

### Install a release binary

1. Download the binary for your OS from the [releases page](https://github.com/praetorian-inc/mithril/releases).

2. (OPTIONAL) Download the `checksums.txt` file to verify the integrity of the archive

```shell
# Check the checksum of the downloaded archive
$ shasum -a 256 mithril_${VERSION}_${ARCH}.tar.gz
b05c4d7895be260aa16336f29249c50b84897dab90e1221c9e96af9233751f22  mithril_${VERSION}_${ARCH}.tar.gz

$ cat mithril_${VERSION}_${ARCH}_checksums.txt | grep mithril_${VERSION}_${ARCH}.tar.gz
b05c4d7895be260aa16336f29249c50b84897dab90e1221c9e96af9233751f22  mithril_${VERSION}_${ARCH}.tar.gz
```

3. Extract the downloaded archive

```shell
$ tar -xvf mitrhil_${VERSION}_${ARCH}.tar.gz
```

4. Move the `mithril` binary into your path:

```shell
$ mv ./mithril /usr/local/bin/
```

### Clone and build yourself

```shell
# clone the Mithril repo
$ git clone https://github.com/praetorian-inc/mithril.git

# navigate into the repo directory and build
$ cd mithril
$ go build

# Move the gokart binary into your path
$ mv ./mithril /usr/local/bin
```

## Usage

There are two main modes of operation for Mithril.

### Run Mithril against static configuration information

```shell
# running without a directory specified defaults to '.'
./mithril <stuff>
```

### Run Mithril in an Istio workload container

```shell
./mithril <stuff>
```

### Get Help

```shell
mithril help
```
