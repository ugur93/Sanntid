package driver

import(
	"fmt"

)


const N_BUTTONS=3


button_array=[][]int({BUTTON_UP1,BUTTON_DOWN1,BUTTON_COMMAND1},
	{BUTTON_UP2,BUTTON_DOWN2,BUTTON_COMMAND2},
	{BUTTON_UP3,BUTTON_DOWN3,BUTTON_COMMAND3})

lamp_array=[][]int({LAMP_UP1,LAMP_DOWN1,LAMP_COMMAND1},
	{LAMP_UP2,LAMP_DOWN2,LAMP_COMMAND2},
	{LAMP_UP3,LAMP_DOWN3,LAMP_COMMAND3})

const(
	BUTTON_CALL_UP int =0
	BUTTON_CALL_DOWN int =1
	BUTTON_COMMAND int = 2

)
const (
	DIRN_DOWN int=-1
	DIRN_STOP int=0
	DIRN_UP int=1

)
func Driver_init(){
	if !IO_init() {
		fmt.Println("Could not initialize IO module")	
	}
	
	Set_stop_lamp(0)
	Set_door_lamp(0)
	Set_floor_indicator(0)


}
func Set_motor_direction(DIRN int){
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

func Get_obstruction_signal() (int) {
	return IO_read_bit(OBSTRUCTION)

}
func Get_stop_signal() (int) {
	return IO_read_bit(STOP)
}
func Set_stop_lamp(value int) {
	IO_clear_bit(LIGHT_STOP)
	if value==1 {
		IO_set_bit(LIGHT_STOP)
	}else{
		IO_clear_bit(LIGHT_STOP)
	}
}
func Set_door_lamp(value int){
	if value==1 {
		IO_set_bit(LIGHT_DOOR_OPEN)
	}else 
		IO_clear_bit(LIGHT_DOOR_OPEN)
}
func Get_floor_sensor_signal() int {
	if IO_read_bit(SENSOR_FLOOR1)==1 {
		return 0
	}else if IO_read_bit(SENSOR_FLOOR2)==1 {
		return 1
	}else if IO_read_bit(SENSOR_FLOOR3)==1 {
		return 2
	}else if IO_read_bit(SENSOR_FLOOR4)==1 {
		return 3
	}else {
		return -1
	}
}

func Set_floor_indicator(floor int) {

	if floor == 2 {
		IO_set_bit(LIGHT_FLOOR_IND1)
	} else {
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor==1 {
		IO_set_bit(LIGHT_FLOOR_IND2)
	}else {
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func Set_button_lamp(BUTTON_TYPE int, floor int, value int){
	if value==1 {
		IO_set_bit(lamp_channel[BUTTON_TYPE][floor])
	}else
		IO_clear_bit(lamp_channel[BUTTON_TYPE][floor])
}
func Get_button_signal(BUTTON_TYPE int,floor int, value int){


	if IO_read_bit(button_array[BUTTON_TYPE][floor])==1 {
		return 1;
	}else {
		return 0;
	}

}
