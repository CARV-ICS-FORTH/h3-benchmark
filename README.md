# H3 Benchmark

**This project is derived from [s3-benchmark](https://github.com/wasabi-tech/s3-benchmark), modified for H3. This README has been updated for H3, but if you need more information, please refer to the original code.**

`h3-benchmark` is a performance testing tool for performing [H3](https://github.com/CARV-ICS-FORTH/H3) operations (PUT, GET, and DELETE) on objects. Besides the storage and bucket configuration, the object size and number of threads can be varied for different tests.

## Building

`h3-benchmark` is written in [Go](https://golang.org). To build, make sure you have a fairly recent version of Go (tested with 1.13.12) and `h3lib` installed. Then run `make`.

## Command Line Arguments

Below are the command line arguments to the program (which can be displayed using -help):

```
  -b string
        Bucket for testing (default "h3-benchmark-bucket")
  -d int
        Duration of each test in seconds (default 60)
  -l int
        Number of times to repeat test (default 1)
  -s string
        H3 storage URI
  -t int
        Number of threads to run (default 1)
  -z string
        Size of objects in bytes with postfix K, M, and G (default "1M")
```

## Example Benchmark

Below is an example run of the benchmark with the default 1MB object size. The benchmark reports for each operation (PUT, GET and DELETE) the results in terms of data speed and operations per second. The program writes all results to the log file `h3-benchmark.log`.

```
$ ./h3-benchmark -s redis://127.0.0.1:6379 -b b1 -d 10
H3 benchmark program v1.0
Parameters: storage_uri=redis://127.0.0.1:6379, bucket=b1, duration=10, threads=1, loops=1, size=1M
Loop 1: PUT time 10.0 secs, objects = 4828, speed = 482.8MB/sec, 482.8 operations/sec. Slowdowns = 0
Loop 1: GET time 10.0 secs, objects = 4709, speed = 470.8MB/sec, 470.8 operations/sec. Slowdowns = 0
Loop 1: DELETE time 0.6 secs, 8201.4 deletes/sec. Slowdowns = 0
```

## Acknowledgements

This project has received funding from the European Union’s Horizon 2020 research and innovation programme under grant agreement No 825061 (EVOLVE - [website](https://www.evolve-h2020.eu>), [CORDIS](https://cordis.europa.eu/project/id/825061)).
