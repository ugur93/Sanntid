package Queue_manager

import "../driver"
import "../Network"
import "../Types"
import "time"
import "fmt"

func Update_outside_lights(External_order Types.Message){
		for i:=0; i<Types.N_FLOORS-1; i++ {
			if External_order.Mask.Outside_order_up[i]==1 {
				state:=External_order.Data.Outside_order_up[i]
				driver.Set_button_lamp(0,i,state)
			}
		}
		for i:=1; i<Types.N_FLOORS; i++ {
			if External_order.Mask.Outside_order_down[i]==1 {
				state:=External_order.Data.Outside_order_down[i]
				driver.Set_button_lamp(1,i,state)
			}
		}
}
func Update_lights(button_type int, floor int,state int){
	driver.Set_button_lamp(button_type,floor,state)
}

func is_queue_empty(local_queue_chan chan Types.Order_queue)(bool){
	for i:=0; i<Types.N_FLOORS; i++ {
		if is_order_in_this_floor(i,local_queue_chan)==true {
			return false
		}
	}
	return true
}
func is_order_in_this_floor(floor int,local_queue_chan chan Types.Order_queue)(bool) {
	local_queue:=<-local_queue_chan
	local_queue_chan<-local_queue
	if local_queue.Outside_order_up[floor]==1 {
		return true
	}else if local_queue.Inside_order[floor]==1{
		return true
	}else if local_queue.Outside_order_down[floor]==1{
		return true
	}
	return false
}
func is_order_in_same_direction(direction int,current_floor int,local_queue_chan chan Types.Order_queue)(bool) {
	
	if direction==Types.DIRN_DOWN {
		for floor := current_floor-1; floor >= 0; floor-- {
			if is_order_in_this_floor(floor,local_queue_chan)==true {
				return true
			}
		}
	}else if direction==Types.DIRN_UP {
		for floor := current_floor+1; floor < Types.N_FLOORS; floor++ {
			if is_order_in_this_floor(floor,local_queue_chan)==true {
				return true
			}
		}
	}
	return false
}

func is_order_in_current_floor(direction,current_floor int,local_queue_chan chan Types.Order_queue)(bool) {
	local_queue:=<-local_queue_chan
	local_queue_chan<-local_queue
	if direction==Types.DIRN_UP {
		if local_queue.Outside_order_up[current_floor]==1 || local_queue.Inside_order[current_floor]==1 {
			return true
		}
	
	}
	if direction==Types.DIRN_DOWN {
		if local_queue.Outside_order_down[current_floor]==1 || local_queue.Inside_order[current_floor]==1 {
			return true
		}
	}
	return false
}

func stop_routine(current_direction,current_floor int,Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue) {
	driver.Set_motor_direction(0)
	driver.Set_door_lamp(1)
	go Update_queue(current_direction,current_floor,0,Broadcast_buffer,local_queue_chan)
	fmt.Println("Welcome to the floor: ",current_floor,",We are in direction(-1D_1U): ",current_direction)
	time.Sleep(3*time.Second)
	driver.Set_door_lamp(0)
}
func Update_queue(current_direction,current_floor int,state int,Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue){
	
	var new_msg Types.Message
	new_msg.MessageType=Types.MT_update

	local_queue:=<-local_queue_chan

	local_queue.LastFloor=current_floor
	local_queue.Moving_direction=current_direction
	if current_direction==Types.DIRN_UP&&current_floor!=3{

		local_queue.Outside_order_up[current_floor]=state
		new_msg.Mask.Outside_order_up[current_floor]=1
		Update_lights(Types.BUTTON_CALL_UP,current_floor,0)

		local_queue.Inside_order[current_floor]=state
		new_msg.Mask.Inside_order[current_floor]=1
		Update_lights(Types.BUTTON_COMMAND,current_floor,0)

	}else if current_direction==Types.DIRN_DOWN&&current_floor!=0 {
		
		local_queue.Outside_order_down[current_floor]=state
		new_msg.Mask.Outside_order_down[current_floor]=1
		Update_lights(Types.BUTTON_CALL_DOWN,current_floor,0)

		local_queue.Inside_order[current_floor]=state
		new_msg.Mask.Inside_order[current_floor]=1
		Update_lights(Types.BUTTON_COMMAND,current_floor,0)		
		
	}else {
		local_queue.Inside_order[current_floor]=state
		new_msg.Mask.Inside_order[current_floor]=1
		Update_lights(Types.BUTTON_COMMAND,current_floor,0)
	}	
	new_msg.Data=local_queue

	local_queue_chan<-local_queue
	Broadcast_buffer<-new_msg
}
func change_direction(current_direction int)(int){
	if current_direction==Types.DIRN_UP {
		return Types.DIRN_DOWN
	}	
	return Types.DIRN_UP
}
func Run_elevator(Broadcast_buffer chan Types.Message,current_floor int,local_queue_chan chan Types.Order_queue){
	var current_direction int
	var temp_current_floor int
	moving:=false
	
	current_direction=Types.DIRN_UP
	direction_changed:=false
	for {
		if is_queue_empty(local_queue_chan)==false {
				for {
				
					if is_order_in_same_direction(current_direction,current_floor,local_queue_chan)==false {
						current_direction=change_direction(current_direction)
						direction_changed=true
					}
					temp_current_floor=driver.Get_floor_sensor_signal()
					if temp_current_floor!=-1 {
						current_floor=temp_current_floor
						driver.Set_floor_indicator(current_floor)
					}

					if temp_current_floor!=-1 && (is_order_in_current_floor(current_direction,current_floor,local_queue_chan)==true)  {
						stop_routine(current_direction,current_floor,Broadcast_buffer,local_queue_chan)
						moving=false
						break
					}else if moving==false && direction_changed==false {
						driver.Set_motor_direction(current_direction)
						moving=true
					}else if direction_changed==true {
						direction_changed=false
					}
					
					time.Sleep(100*time.Millisecond)
				}

			}
		}

}

