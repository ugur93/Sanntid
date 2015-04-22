package Elevator

import "../driver"
import "../Types"
import "time"
import "fmt"

func Elevator_init(Broadcast_buffer chan Types.Message,current_floor int, local_queue_chan chan Types.Order_queue,Delete_order_ch chan Types.Order_queue ) {

	Run_elevator(Broadcast_buffer,current_floor,local_queue_chan,Delete_order_ch)

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
func are_there_orders_in_this_direction(direction int,current_floor int,local_queue_chan chan Types.Order_queue)(bool) {
	
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

func stop_routine(current_direction,current_floor int,Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue,Delete_order_ch chan Types.Order_queue) {
	driver.Set_motor_direction(0)
	driver.Set_door_lamp(1)
	go Delete_order(current_direction,current_floor,Broadcast_buffer,local_queue_chan,Delete_order_ch)
	fmt.Println("Welcome to the floor: ",current_floor,",We are in direction(-1D_1U): ",current_direction)
	time.Sleep(3*time.Second)
	driver.Set_door_lamp(0)
}
func Delete_order(current_direction,current_floor int,Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue,Delete_order_ch chan Types.Order_queue){
	
	var new_msg Types.Message
	new_msg.Message_type=Types.MT_update
	
	Update_state_information(current_direction,current_floor,local_queue_chan)
	
	if current_direction==Types.DIRN_UP || current_floor==0 {
		new_msg.Mask.Outside_order_up[current_floor]=1
		new_msg.Mask.Inside_order[current_floor]=1
	}
	if current_direction==Types.DIRN_DOWN || current_floor==Types.N_FLOORS-1 {	
		new_msg.Mask.Outside_order_down[current_floor]=1
		new_msg.Mask.Inside_order[current_floor]=1
	}else {
		new_msg.Mask.Inside_order[current_floor]=1
	}
	Delete_order_ch<-new_msg.Mask	
	//Queue_manager.Update_queue(true,0,new_msg.Mask,Broadcast_buffer,local_queue_chan)
}
func Update_state_information(current_direction,current_floor int,local_queue_chan chan Types.Order_queue){
	
	local_queue:=<-local_queue_chan
	local_queue.Last_floor=current_floor
	local_queue.Moving_direction=current_direction
	local_queue_chan<-local_queue

}
func change_direction(current_direction int)(int){
	if current_direction==Types.DIRN_UP {
		return Types.DIRN_DOWN
	}	
	return Types.DIRN_UP
}
func Run_elevator(Broadcast_buffer chan Types.Message,current_floor int,local_queue_chan chan Types.Order_queue,Delete_order_ch chan Types.Order_queue){
	var current_direction int
	var temp_current_floor int
	moving:=false
	
	current_direction=Types.DIRN_UP
	direction_changed:=false
	for {
		if moving==false && driver.Get_floor_sensor_signal()==-1 {
				current_floor=driver.Get_to_defined_state()
		
		}
		if is_queue_empty(local_queue_chan)==false {
				for {
				
					if are_there_orders_in_this_direction(current_direction,current_floor,local_queue_chan)==false {
						if (is_order_in_current_floor(current_direction,current_floor,local_queue_chan)==true){
							stop_routine(current_direction,current_floor,Broadcast_buffer,local_queue_chan,Delete_order_ch)
							moving=false
							break
						}
						current_direction=change_direction(current_direction)
						direction_changed=true
					}
					temp_current_floor=driver.Get_floor_sensor_signal()
					if temp_current_floor!=-1 {
						current_floor=temp_current_floor
					}

					if temp_current_floor!=-1 && (is_order_in_current_floor(current_direction,current_floor,local_queue_chan)==true)  {
						stop_routine(current_direction,current_floor,Broadcast_buffer,local_queue_chan,Delete_order_ch)
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
			time.Sleep(100*time.Millisecond)
		}
		

}

