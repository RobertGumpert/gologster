package logger

import (
	"strconv"
	"testing"
	"time"
)

type User struct {
	Name string
	Mail *Email
}

type Email struct {
	Box string
}

var (
	f = newBasicFile(new(loggerBasic), map[string]string{
		"1": "./tests/logs/logs_1.txt",
		"2": "./tests/logs/logs_2.txt",
		"3": "./tests/logs/logs_3.txt",
		"4": "./tests/logs/logs_4.txt",
	})
	fileAsync = newLoggerFileMultithreading(f)
	basicFile = newLoggerFileMutex(f)
)

// go tool pprof -pdf profile.out > profile_cpu.pdf
// go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out

func BenchmarkAsync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		asyncWriter(fileAsync, i, "1")
	}
}

func BenchmarkBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basicWriter(basicFile, i, "1")
	}
}

func asyncWriter(asnc *loggerFileMultithreading, i int, key string) {
	asnc.add(&User{
		Name: strconv.Itoa(i),
	}, levelInfo, time.Now().Format("Mon Jan _2 15:04:05 2006"), key)
}

func basicWriter(bsc *loggerFileMutex, i int, key string) {
	bsc.add(&User{
		Name: strconv.Itoa(i),
	}, levelInfo, time.Now().Format("Mon Jan _2 15:04:05 2006"), key)
}
