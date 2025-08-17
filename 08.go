package main

/*
#include <stdlib.h>
#include <stdio.h>
void hello() {
	printf("hello from C\n");
}

void goodbye(char* s) {
	printf("Goodbye %s\n", s);
}
*/
import "C"
import "unsafe"

func main() {
	C.hello()
	s := C.CString("cruel world!\n")
	C.goodbye(s)
	C.free(unsafe.Pointer(s))
}
