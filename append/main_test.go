package main

import (
	"testing"
)

func BenchmarkAppendMake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AppendMake()
	}
}

func BenchmarkAppendSimple(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AppendSimple()
	}
}
