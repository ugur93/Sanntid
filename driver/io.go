package driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func IO_init() bool {
	return bool(int(C.io_init())== 1)
}

func IO_set_bit(channel int){
	C.io_set_bit(C.int(channel))
}

func IO_clear_bit(channel int){
	C.io_clear_bit(C.int(channel))
}

func IO_write_analog(channel int,value int){
	C.io_write_analog(C.int(channel),C.int(value))
}

func IO_read_bit(channel int) int{

	return int(C.io_read_bit(C.int(channel)))	
}

func IO_read_analog(channel int) int{
	return int(C.io_read_analog(C.int(channel)))
}