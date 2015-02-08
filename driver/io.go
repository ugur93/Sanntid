package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h
*/
import "C"

func IO_init() bool {

	return bool(int(C.io_init())!= 1)

}

func IO_set_bit(channel int){
	C.io_set_bit(C.int(channel))


}
