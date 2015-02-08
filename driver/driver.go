package driver

import(
	//"fmt"

)


const N_BUTTONS=3
const (
	DIRN_DOWN int=-1
	DIRN_STOP int=0
	DIRN_UP int=1

)

func set_motor_direction(DIRN int){
	if DIRN == 0 {
		IO_write_analog(MOTOR,0)
	} else if DIRN>0 {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR,2800)
	}else if DIRN<0 {
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR, 2800)
	}


}