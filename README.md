# dict2rest

A simple proxy service to provide an HTTP (REST) interface to a Dict protocol
(RFC 2229) server. Written in Go.

## Installation

This project currently uses [gb](https://getgb.io) as its build tool. All
dependencies are in this repository.

Assuming you have `gb` installed it should be as simple as this:

```shell
git clone https://github.com/felix/go-dict2rest
cd go-dict2rest
gb build
./bin/dict2rest --host dict.org
```

## Usage

```
$ dict2rest --help

Usage of dict2rest:
  -dicthost string
        Dict server name (default "localhost")
  -dictport string
        Dict server port (default "2628")
  -gzip
        Enable gzip compression
  -port string
        Listen port (default "8080")
```

## License

Copyright Felix Hanley, 2016

See LICENSE file.
