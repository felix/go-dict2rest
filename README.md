# dict2rest

A simple proxy service to provide an HTTP (REST) interface to a
[Dict protocol (RFC 2229)](http://tools.ietf.org/html/rfc2229) for the Dict
protocol. Written in Go.

An instance can _usually_ be access from
http://felixhanley.info/projects/lahu/dictionary/

## Installation

This project currently uses [gb](https://getgb.io) as its build tool. All
dependencies are in this repository.

Assuming you have `gb` installed it should be as simple as this:

```shell
git clone https://github.com/felix/go-dict2rest
cd go-dict2rest
gb build
./bin/dict2rest --dicthost dict.org
```

The server binary has the following options:

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

## Usage

The URL endpoints try to match the commands defined in RFC 2229. Results are
JSON and have the following format (newlines added for readability):

GET /databases

```json
[
    {
        "name":"eng-lhu",
        "description":"English-Lahu Dictionary ver. 20160128"
    }, {
        "name":"lhu-eng",
        "description":"Lahu-English Dictionary ver. 20160128"
    }
]
```

GET /define/{word}  
GET /define/{word}?dict=eng-lhu

```json
[
    {
        "dictionary":"WordNet (r) 3.0 (2006)",
        "word":"lahu",
        "definition":"Lahu\n    n 1: a Loloish language\n"
    }
]
```

Where 'eng-lhu' is one of the names of the server's dictionaries.

RFC 2229 error codes are passed through as JSON:

```json
{"code":552,"message":"no match"}
```

## TODO

- Add strategy listing and lookups
- Add server and status commands
- Add tests

## License

Copyright Felix Hanley, 2016

See LICENSE file.
