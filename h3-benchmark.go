// h3-benchmark.go
// Copyright (c) 2020 FORTH-ICS
// Copyright (c) 2017 Wasabi Technology, Inc.

package main

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lh3lib
// #include <stdlib.h>
// #include "h3wrapper.h"
import "C"
import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

// Global variables
var storage_uri, bucket string
var duration_secs, threads, loops int
var object_size uint64
var object_data []byte
var running_threads, upload_count, download_count, delete_count, upload_slowdown_count, download_slowdown_count, delete_slowdown_count int32
var endtime, upload_finish, download_finish, delete_finish time.Time

var handle [256]C.H3_Handle
var bucket_str *C.char

func logit(msg string) {
	fmt.Println(msg)
	logfile, _ := os.OpenFile("benchmark.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if logfile != nil {
		logfile.WriteString(time.Now().Format(http.TimeFormat) + ": " + msg + "\n")
		logfile.Close()
	}
}

func createBucket(ignore_errors bool) {
	status := C.h3lib_create_bucket(handle[0], bucket_str)
	if status != C.H3_SUCCESS && status != C.H3_EXISTS {
		log.Fatalf("FATAL: Unable to create bucket %s (error: %d)", bucket, status)
	}
}

func deleteAllObjects() {
	status := C.h3lib_purge_bucket(handle[0], bucket_str)
	if status != C.H3_SUCCESS {
		log.Fatalf("FATAL: Unable to create bucket %s (error: %d)", bucket, status)
	}
}

func runUpload(thread_num int) {
	for time.Now().Before(endtime) {
		objnum := atomic.AddInt32(&upload_count, 1)
		objName := fmt.Sprintf("Object-%d", objnum)

		objectName := C.CString(objName)
		defer C.free(unsafe.Pointer(objectName))

		status := C.h3lib_write_object(handle[thread_num], bucket_str, objectName, unsafe.Pointer(&object_data[0]), C.ulonglong(object_size), C.ulonglong(0))
		if status != C.H3_SUCCESS{
			log.Fatalf("FATAL: Error uploading object %s (error: %d)", objName, status)
		} else {
			// logit(fmt.Sprintf("Uploaded %s", objName))
		}
	}
	// Remember last done time
	upload_finish = time.Now()
	// One less thread
	atomic.AddInt32(&running_threads, -1)
}

func runDownload(thread_num int) {
	for time.Now().Before(endtime) {
		atomic.AddInt32(&download_count, 1)
		objnum := rand.Int31n(upload_count) + 1
		objName := fmt.Sprintf("Object-%d", objnum)

		objectName := C.CString(objName)
		defer C.free(unsafe.Pointer(objectName))

		status := C.h3lib_read_dummy_object(handle[thread_num], bucket_str, objectName)
		if status != C.H3_SUCCESS{
			log.Fatalf("FATAL: Error downloading object %s (error: %d)", objName, status)
		} else {
			// logit(fmt.Sprintf("Downloaded %s", objName))
		}
	}
	// Remember last done time
	download_finish = time.Now()
	// One less thread
	atomic.AddInt32(&running_threads, -1)
}

func runDelete(thread_num int) {
	for {
		objnum := atomic.AddInt32(&delete_count, 1)
		if objnum > upload_count {
			break
		}

		objName := fmt.Sprintf("Object-%d", objnum)
		objectName := C.CString(objName)
		defer C.free(unsafe.Pointer(objectName))

		status := C.h3lib_delete_object(handle[thread_num], bucket_str, objectName)
		if status != C.H3_SUCCESS {
			log.Fatalf("FATAL: Error deleting object %s (error: %d)", objName, status)
		} else {
			// logit(fmt.Sprintf("Deleted %s", objName))
		}
	}
	// Remember last done time
	delete_finish = time.Now()
	// One less thread
	atomic.AddInt32(&running_threads, -1)
}

func main() {
	// Hello
	fmt.Println("H3 benchmark program v1.0")

	// Parse command line
	myflag := flag.NewFlagSet("myflag", flag.ExitOnError)
	myflag.StringVar(&storage_uri, "s", "", "H3 storage URI")
	myflag.StringVar(&bucket, "b", "h3-benchmark-bucket", "Bucket for testing")
	myflag.IntVar(&duration_secs, "d", 60, "Duration of each test in seconds")
	myflag.IntVar(&threads, "t", 1, "Number of threads to run")
	myflag.IntVar(&loops, "l", 1, "Number of times to repeat test")
	var sizeArg string
	myflag.StringVar(&sizeArg, "z", "1M", "Size of objects in bytes with postfix K, M, and G")
	if err := myflag.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	// Check the arguments
	if storage_uri == "" {
		log.Fatal("Missing argument -s for storage URI.")
	}
	if threads > 256 {
		log.Fatalf("Threads can be up to 256.")
	}
	var err error
	if object_size, err = bytefmt.ToBytes(sizeArg); err != nil {
		log.Fatalf("Invalid -z argument for object size: %v", err)
	}

	// Echo the parameters
	logit(fmt.Sprintf("Parameters: storage_uri=%s, bucket=%s, duration=%d, threads=%d, loops=%d, size=%s",
		storage_uri, bucket, duration_secs, threads, loops, sizeArg))

	storage_uri_str := C.CString(storage_uri)
	defer C.free(unsafe.Pointer(storage_uri_str))
	for n := 0; n < threads; n++ {
		handle[n] = C.H3_Init(storage_uri_str)
		if handle[n] == nil {
			log.Fatal("Unable to initialize h3lib")
		}
		defer C.H3_Free(handle[n])
	}

	bucket_str = C.CString(bucket)
	defer C.free(unsafe.Pointer(bucket_str))

	// Initialize data for the bucket
	object_data = make([]byte, object_size)
	rand.Read(object_data)
	hasher := md5.New()
	hasher.Write(object_data)

	// Create the bucket and delete all the objects
	createBucket(true)
	deleteAllObjects()

	// Loop running the tests
	for loop := 1; loop <= loops; loop++ {

		// reset counters
		upload_count = 0
		upload_slowdown_count = 0
		download_count = 0
		download_slowdown_count = 0
		delete_count = 0
		delete_slowdown_count = 0

		// Run the upload case
		running_threads = int32(threads)
		starttime := time.Now()
		endtime = starttime.Add(time.Second * time.Duration(duration_secs))
		for n := 0; n < threads; n++ {
			go runUpload(n)
		}

		// Wait for it to finish
		for atomic.LoadInt32(&running_threads) > 0 {
			time.Sleep(time.Millisecond)
		}
		upload_time := upload_finish.Sub(starttime).Seconds()

		bps := float64(uint64(upload_count)*object_size) / upload_time
		logit(fmt.Sprintf("Loop %d: PUT time %.1f secs, objects = %d, speed = %sB/sec, %.1f operations/sec. Slowdowns = %d",
			loop, upload_time, upload_count, bytefmt.ByteSize(uint64(bps)), float64(upload_count)/upload_time, upload_slowdown_count))

		// Run the download case
		running_threads = int32(threads)
		starttime = time.Now()
		endtime = starttime.Add(time.Second * time.Duration(duration_secs))
		for n := 0; n < threads; n++ {
			go runDownload(n)
		}

		// Wait for it to finish
		for atomic.LoadInt32(&running_threads) > 0 {
			time.Sleep(time.Millisecond)
		}
		download_time := download_finish.Sub(starttime).Seconds()

		bps = float64(uint64(download_count)*object_size) / download_time
		logit(fmt.Sprintf("Loop %d: GET time %.1f secs, objects = %d, speed = %sB/sec, %.1f operations/sec. Slowdowns = %d",
			loop, download_time, download_count, bytefmt.ByteSize(uint64(bps)), float64(download_count)/download_time, download_slowdown_count))

		// Run the delete case
		running_threads = int32(threads)
		starttime = time.Now()
		endtime = starttime.Add(time.Second * time.Duration(duration_secs))
		for n := 0; n < threads; n++ {
			go runDelete(n)
		}

		// Wait for it to finish
		for atomic.LoadInt32(&running_threads) > 0 {
			time.Sleep(time.Millisecond)
		}
		delete_time := delete_finish.Sub(starttime).Seconds()

		logit(fmt.Sprintf("Loop %d: DELETE time %.1f secs, %.1f deletes/sec. Slowdowns = %d",
			loop, delete_time, float64(upload_count)/delete_time, delete_slowdown_count))
	}

	// All done
}
