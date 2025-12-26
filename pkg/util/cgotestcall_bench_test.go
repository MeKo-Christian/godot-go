package util

import (
	"testing"
	"unsafe"
)

func BenchmarkCgoTestCall(b *testing.B) {
	var v int
	ptr := unsafe.Pointer(&v)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CgoTestCall(ptr)
	}
}

func BenchmarkGoNoop(b *testing.B) {
	var v int
	ptr := unsafe.Pointer(&v)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ptr
	}
}
