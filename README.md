resp
====

[![Build Status](https://travis-ci.org/stvp/resp.svg?branch=master)](https://travis-ci.org/stvp/resp)

resp is an (incomplete) Go package that provides helpful tools for
reading and writing [Redis protocol][resp] objects.

[Documentation][docs]

[resp]: http://redis.io/topics/protocol
[docs]: http://godoc.org/github.com/stvp/resp

Benchmarks
----------

Benchmarks run on a 2015 MacBook with 1.1 GHz Intel Core M:

    BenchmarkCommandSlices-4              	10000000	       211 ns/op	      48 B/op	       1 allocs/op
    BenchmarkReaderReadObjectSliceSmall-4 	20000000	        98.7 ns/op	       0 B/op	       0 allocs/op
    BenchmarkReaderReadObjectSliceMedium-4	 5000000	       277 ns/op	       0 B/op	       0 allocs/op
    BenchmarkParseLenLine-4               	100000000	        19.3 ns/op	       0 B/op	       0 allocs/op

