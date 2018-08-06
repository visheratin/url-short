# URL shortening microservice

This repository contains a prototype of a microservice that performs back and forth conversion of URLs into shortened 
links of a fixed length. Length of the code can be varied easily through the config file. When using this microservice, 
do not change the length of the code once it was set, because there are no handlers for such situations (mostly because 
such situations usually do not happen in production).

Microservice uses TCP connections for receiving and sending data to eliminate overheads of HTTP.

Microservice permanently stores all its data on the local disk, but this can be easily changed through creating new type  
of storage by implementing a **Storage** interface. Depending on requirements for permanent data storing, saving the input-code 
pairs can be done in goroutines (if the saving fails code will still be returned to the client, but after service restart it will not 
be available) or the code can be not saved permanently at all (pass **nil** into **Converter** constructor, execution will be faster 
by several orders of magnitude).

In order to find how fast the service is, there are three benchmarks for it. Results of their execution on notebook with 
Core i7-7700HQ, SSD hard drive and Windows OS:

```
// generating codes without writing them to local disk
BenchmarkLoad-8                1000000	      1038 ns/op	      75 B/op	       0 allocs/op
```
```
// generating codes with writing them to local disk
BenchmarkLoadWithStorage-8         500	   3561159 ns/op	    1243 B/op	      20 allocs/op
```
```
// making requests for generating codes through TCP
BenchmarkListener-8   	           200	   6958286 ns/op	    5670 B/op	      21 allocs/op
```

Results of CPU profiling (https://github.com/visheratin/url-short/blob/master/cpu-profile.svg) show that the most of the 
time is spent on writing files to the disk and the URL shortening operation is very fast.
