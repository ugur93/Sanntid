package driver

import(
	"fmt"
	"time"
)

const N_BUTTONS=3
const N_FLOORS = 4

var button_array =[][] int{
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},	
	/*{BUTTON_UP1,BUTTON_UP2,BUTTON_UP3,BUTTON_UP4},
	{BUTTON_DOWN1,BUTTON_DOWN2,BUTTON_DOWN3,BUTTON_DOWN4},
	{BUTTON_COMMAND1,BUTTON_COMMAND2,BUTTON_COMMAND3,BUTTON_COMMAND4},*/
}

var lamp_array =[][]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	/*{LIGHT_UP1,LIGHT_UP2,LIGHT_UP3,LIGHT_UP4},
	{LIGHT_DOWN1,LIGHT_DOWN2,LIGHT_DOWN3,LIGHT_DOWN4},
	{LIGHT_COMMAND1,LIGHT_COMMAND2,LIGHT_COMMAND3,LIGHT_COMMAND4},*/
}

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
func Driver_init()(int){
	if !IO_init() {
		fmt.Println("Could not initialize IO module")
	}
	current_floor:=Get_to_defined_state()
	for i:=0; i<4; i++ {
		if i < 3 {
			Set_button_lamp(0,i,0)
		}
		if i!=0 {
			Set_button_lamp(1,i,0)
		}
		Set_button_lamp(2,i,0)
	}
	Set_stop_lamp(0)
	Set_door_lamp(0)
	Set_floor_indicator(0)
	return current_floor
	//current_floor = Get_to_defined_state()
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
	}else{
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
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

	if floor == 0x02 {
		IO_set_bit(LIGHT_FLOOR_IND1)
	} else {
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor == 0x01 {
		IO_set_bit(LIGHT_FLOOR_IND2)
	}else {
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}
}
func Set_button_lamp(BUTTON_TYPE int, floor int, value int){
	//fmt.Println(lamp_array[0][0])
	if value==1 {
		IO_set_bit(lamp_array[floor][BUTTON_TYPE])
	}else{
		IO_clear_bit(lamp_array[floor][BUTTON_TYPE])
	}
}
func Get_button_signal(BUTTON_TYPE int,floor int) int{


	if IO_read_bit(button_array[floor][BUTTON_TYPE])==1 {
		return 1;
	} else {
		return 0;
	}
}

func Check_for_outside_order(outside_order_ch chan [2]int) {

	order_array := [2]int{}
	
	pressed_UP:=[]int{0,0,0,0}
	pressed_DOWN:=[]int{0,0,0,0}
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
		
		
			if Get_button_signal(BUTTON_CALL_DOWN, floor)==0 && pressed_DOWN[floor]==1 {
					pressed_DOWN[floor]=0;
			
			} else if Get_button_signal(BUTTON_CALL_UP, floor)==0 && pressed_UP[floor]==1 {
					pressed_UP[floor]=0
			}
			if floor!=0&&Get_button_signal(BUTTON_CALL_DOWN, floor)==1&&pressed_DOWN[floor]==0 {
				pressed_DOWN[floor]=1
				order_array[1] = BUTTON_CALL_DOWN
				order_array[0] = floor
				outside_order_ch <- order_array
			} else if floor!=3&&Get_button_signal(BUTTON_CALL_UP, floor)==1&&pressed_UP[floor]==0 {
				pressed_UP[floor]=1
				order_array[1] = BUTTON_CALL_UP
				order_array[0] = floor
				outside_order_ch <- order_array
			}
		}
		time.Sleep(100*time.Millisecond) 
	}
}
func Check_for_inside_order(inside_order_ch chan int){
	var order int
	pressed:=[]int{0,0,0,0}
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if Get_button_signal(BUTTON_COMMAND, floor) == 0 && pressed[floor] ==1 {
				pressed[floor] = 0
			}
			if Get_button_signal(BUTTON_COMMAND, floor) == 1 && pressed[floor] == 0 {
				pressed[floor] = 1
				order = floor
				inside_order_ch <- order
			}
		}
		time.Sleep(100*time.Millisecond)
	}
}
func Get_to_defined_state()(current_floor int){
	if Get_floor_sensor_signal() != -1 {
		return Get_floor_sensor_signal()
	} else {
		Set_motor_direction(DIRN_DOWN)
		for {
			if Get_floor_sensor_signal() != -1 {
					Set_motor_direction(DIRN_STOP)
					return Get_floor_sensor_signal()
			}
		}
	}
}
