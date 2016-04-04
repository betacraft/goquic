goquic, QUIC support for Go
===========================

This is a work-in-progress QUIC implementation for Go. This is based on
[libquic](https://github.com/devsisters/libquic) library, which is in turn based
on original QUIC implementation on [Chromium](http://www.chromium.org/quic).

QUIC is an experimental protocol aimed at reducing web latency over that of TCP.
On the surface, QUIC is very similar to TCP+TLS+SPDY implemented on UDP. Because
TCP is implement in operating system kernels, and middlebox firmware, making
significant changes to TCP is next to impossible. However, since QUIC is built
on top of UDP, it suffers from no such limitations.

Key features of QUIC over existing TCP+TLS+SPDY include

  * Dramatically reduced connection establishment time
  * Improved congestion control
  * Multiplexing without head of line blocking
  * Forward error correction
  * Connection migration

## Project Status

*This library is highly experimental.* Although `libquic` sources are from
Chromium (which are tested), the Go bindings are still highly pre-alpha state.

Known issues:

  * No support for read streaming. All request must fit in memory.
  * Secure QUIC not fully tested. May not support ECDSA certificates.

Things to do:

  * Read streaming support

## Preliminary Benchmarks

A very primitive benchmark testing have been done. Testing environments below:

| Items        | Description                                               |
| ------------ | --------------------------------------------------------- |
| Optimization | libquic built with `-O3` parameters                       |
| CPU          | Intel(R) Core(TM) i7-4930K CPU @ 3.40GHz                  |
| Server Code  | https://github.com/devsisters/goquic/blob/master/example/server.go |
| Server Parms | `GOMAXPROCS=12 ./server -port 9090 -n 12`                 |
| Client Code  | https://github.com/devsisters/quicbench/blob/master/quicbench.go |
| Client Parms | `./quicbench -u="https://example.com:9090/" -c 200 -r 1000` |

The server code is modified to create 30B, 1kB, 5kB, 10kB HTTP body payload.
Concurrency is 200 and each thread requests 1,000 requests. It is designed to
measure ideal throughput of the server. Naturally the throughput goes down when
concurrency increases.

Benchmark results:

| Payload Size | Requests per Second |
| ------------ | ------------------- |
| 30B Payload  | 23832.18 RPS        |
| 1kB Payload  | 21704.84 RPS        |
| 5kB Payload  | 9343.58 RPS         |
| 10kB Payload | 5312.75 RPS         |

On 10kB case, calculating the total network throughput is `435Mbps`.

How many connections per second can this server process?

`./gobench -u="https://example.com:9090/" -c 200 -r 100 -qk=false`

Turning off keepalive using `qk` option results in a pure new QUIC connection
per request. The benchmark results are `2905.58 CPS`.


Getting Started
===============

## Get source files
Very important follow the exact method,

```
mkdir $GOPATH/src/github.com/devsisters
cd $GOPATH/src/github.com/devsisters
git clone https://github.com/betacraft/goquic.git
```

## Build static library files

Although prebuilt static library files already exists in the repository for
convenience, it is always good practice to build library files from source. You
should not trust any unverifiable third-party binaries.

To build the library files for your architecture and OS:

```bash
./build_libs.sh
```

This will fetch `libquic` master and build all the binaries from source. The
C/C++ files for Go bindings will be all built too.

Currently Linux and Mac OS X and FreeBSD is supported.

#### Requirements
* Cmake
* Ninja

## How to build

If you are using Go >= 1.5, you can build goquic binaries without any extra work.

```bash
go build $GOPATH/src/github.com/devsisters/goquic/example/server.go
```

If you are using Go 1.4, you should open goquic.go and manually edit ${SRCDIR}
with your real path (maybe /YOUR/GOPATH/src/github.com/devsisters/goquic).



SPDY/QUIC support
=================

We have a experimental SPDY/QUIC implementation as a library.
You can use this library to add SPDY/QUIC support for your existing Go HTTP server.

See our SPDY-QUIC server/client implementation [here](example/).

## How to use server

When running a HTTP server, do:

```go
goquic.ListenAndServe(":8080", 1, nil)
```

instead of

```go
http.ListenAndServe(":8080", nil)
```

## How to use client

You need to create http.Client with Transport changed, do:

```go
client := &http.Client{
	Transport: goquic.NewRoundTripper(false),
}
resp, err := client.Get("http://example.com/")
```

instead of

```go
resp, err := http.Get("http://example.com/")
```


### Pushing code and Upgrading code

Code must pe pushed to the fork remote, while pushing code, do
```
git push fork master
```

Taking latest code, do:
```
git pull fork master
```
