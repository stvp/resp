resp
====

[![Build Status](https://travis-ci.org/stvp/resp.svg?branch=master)](https://travis-ci.org/stvp/resp)

resp is an in-progress Go package that provides helpful tools for reading and
writing [Redis protocol][resp] objects. It is under active development and the
API is not stable.

[Documentation][docs]

[resp]: http://redis.io/topics/protocol
[docs]: http://godoc.org/github.com/stvp/resp

Benchmarks
----------

All benchmarks run on a MacBook Pro (Retina) 2.3 GHz Intel Core i7

**5/1/14**

    BenchmarkCommandSlices	10000000	       160 ns/op
    BenchmarkReaderReadObjectSliceSmall	20000000	        74.5 ns/op	(294 MB/sec)
    BenchmarkReaderReadObjectSliceMedium	10000000	       232 ns/op	(390 MB/sec)
    BenchmarkParseLenLine	100000000	        15.6 ns/op

