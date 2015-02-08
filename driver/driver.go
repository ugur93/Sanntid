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

func get_obstruction_signal() (int) {
	return IO_read_bit(OBSTRUCTION)

}
func get_stop_signal() (int) {
	return IO_read_bit(STOP)
}
func set_stop_lamp(value int) {
	if value {
		IO_set_bit(LIGHT_STOP)
	}else{
		IO_clear_bit(LIGHT_STOP)
	}
}
func get_floor_sensor_signal() int {
	if IO_read_bit(SENSOR_FLOOR1) {
		return 0
	}else if IO_read_bit(SENSOR_FLOOR2) {
		return 1
	}else if IO_read_bit(SENSOR_FLOOR3) {
		return 2
	}else if IO_read_bit(SENSOR_FLOOR4) {
		return 3
	}else {
		return -1
	}
}

func set_floor_indicator(floor int) {

	if floor & 0x02 {
		IO_set_bit(LIGHT_FLOOR_IND1)
	} else {
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
}