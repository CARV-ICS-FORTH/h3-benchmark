default: all

h3wrapper: h3wrapper.h h3wrapper.c
	gcc -c h3wrapper.c

all: h3wrapper h3-benchmark.go
	go build

clean:
	rm -f h3wrapper.o h3-benchmark
