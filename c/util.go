package main

/*
#include <string.h>
*/
import "C"
import "unsafe"

func goByteSlice(s *C.char, length C.size_t) []byte {
	tmpSlice := (*[1 << 16]C.char)(unsafe.Pointer(s))[:length:length]
	results := make([]byte, length)
	for i, s := range tmpSlice {
		results[i] = byte(s)
	}
	return results
}
