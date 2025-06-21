package main

import "testing"

func BenchmarkBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basic()
	}
}

func BenchmarkConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concurrent()
	}
}
