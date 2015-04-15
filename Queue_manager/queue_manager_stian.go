package Queue_manager

import "../driver"
import "../Network"
import "../Types"
//import "time"
//import "../driver" //HOW DOES ONE DO THIS

var local_queue Types.Order_queue;

func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	
	current_floor:=driver.Driver_init()
	local_queue.LastFloor=current_floor
	local_queue.Moving=false
	//Stian:
	//Når det skjer en oppdatering på Queue matrise send gjennom new_message channel
	//Når det skjer en oppdatering på Queue fra andre heiser kommer det innpå Queue chan
	Broadcast_order:=make(chan Network.Message,1)
	Queue_lock_chan:=make(chan int,1) //Bruk det når det leses på Queue
	Queue_Network_lock_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Received_order_ch:=make(chan Network.Message,1)
	Broadcast_buffer:=make(chan Network.Message,500)
	Queue_lock_chan<-1
	go Network.Network_Manager_init("20020",Broadcast_order,stop_chan,Received_order_ch,Queue_Network_lock_chan)
	go Get_orders(Received_order_ch)
	go Send_orders(Broadcast_order,Broadcast_buffer)
	//go Run_elevator(Broadcast_order)
	stop_chan<-1
}

func go_to_lowest_up_order(Broadcast_buffer chan Network.Message){
	var lowest_up_order int
	for floor := 0; floor < Types.N_FLOORS; floor++{
		if local_queue.Outside_order_up[floor] == 1{
			lowest_up_order = floor
			break
		}
	}
	for {
		if driver.Get_floor_sensor_signal() == lowest_up_order {
			stop_at_floor("up",Broadcast_buffer)
			driver.Set_motor_direction(1)
			local_queue.Moving_direction = 1
		}
	}	
}

func stop_at_all_up_orders(Broadcast_buffer chan Network.Message){
	var current_floor int
	var floor int
	var Queue_update Network.Message
	Queue_update.MessageType="Update"
	for {
		if driver.Get_floor_sensor_signal() != 1{
			current_floor = driver.Get_floor_sensor_signal()
			if local_queue.Outside_order_up[current_floor]==1 || local_queue.Inside_order[current_floor]==1 {
				stop_at_floor("up",Broadcast_buffer)
				for floor = current_floor+1; floor < Types.N_FLOORS-1; floor++ {
					if local_queue.Outside_order_up[floor]==1 || local_queue.Inside_order[floor]==1 {
						driver.Set_motor_direction(1)
						local_queue.Moving = true
						Queue_update.Data=local_queue
						Broadcast_buffer <- Queue_update
						break
					}
				}
				if (floor == (Types.N_FLOORS-1)) && (local_queue.Moving == false){
					break
				}
			}
		}
	}
}
func go_to_highest_down_order(Broadcast_buffer chan Network.Message){
	var highest_down_order int
	for floor := Types.N_FLOORS-1; floor > -1; floor--{
		if local_queue.Outside_order_down[floor] == 1{
			highest_down_order = floor
			break
		}
	}
	for {
		if driver.Get_floor_sensor_signal() == highest_down_order {
			stop_at_floor("down",Broadcast_buffer)
			driver.Set_motor_direction(-1)
			local_queue.Moving_direction = -1
		}
	}	
}

func stop_at_all_down_orders(Broadcast_buffer chan Network.Message){
	var Queue_update Network.Message
	Queue_update.MessageType="Update"
	var current_floor int
	var floor int
	for {
		if driver.Get_floor_sensor_signal() != 1{
			current_floor = driver.Get_floor_sensor_signal()
			if local_queue.Outside_order_down[current_floor]==1 || local_queue.Inside_order[current_floor]==1 {
				stop_at_floor("down",Broadcast_buffer,current_floor)
				for floor = current_floor-1; floor > -1; floor-- {
					if local_queue.Outside_order_down[floor]==1 || local_queue.Inside_order[floor]==1 {
						driver.Set_motor_direction(-1)
						local_queue.Moving = true
						Queue_update.Data=local_queue
						Broadcast_buffer <- Queue_update
						break
					}
				}
				if (floor == 0) && (local_queue.Moving == false){
					break
				}
			}
		}
	}
}

func stop_at_floor(direction string,Broadcast_buffer chan Network.Message,current_floor int){
	var mask Types.Order_queue
	var msg Network.Message
	msg.MessageType="Update"
	driver.Set_motor_direction(0)
	local_queue.Moving = false
	if direction == "up" {
		local_queue.Outside_order_up[current_floor] = 0
		mask.Outside_order_up[current_floor] = 1
	}else if direction=="down" {
		local_queue.Outside_order_down[current_floor] = 0
		mask.Outside_order_down[current_floor] = 1
	}
	local_queue.Inside_order[current_floor] = 0
	mask.Outside_inside_order[current_floor] = 1
	driver.Set_door_lamp(1)
	msg.Mask=mask
	Broadcast_buffer <- mask
	time.Sleep(3000*time.Millisecond)
	driver.Set_door_lamp(0)
}

func unhandled_order_in_queue() {
	for floor := 0; floor < Types.N_FLOORS-1; floor++{
		if local_queue.Outside_order_up[floor] ==1 {
			return true
		}
	}
}

func Run_elevator(Broadcast_buffer chan Network.Message){
	//Bør kanskje ha en guard som passer på at heisen ikke kjøre for langt (ved topp eller bunn)
	//Skal kjøre så heisen så lenge det finnes bestillinger i egen kø
	//Sender update til Network når en bestilling er utført og fjernet fra kø
	//Skrur ikke av lys på etasjepanel
	idle = true
	var handling_direction int
	var lowest_up_order int
	var highest_down_order int
	for {
		for {if !idle {break}
			//if unhandled_order_in_queue(){
				idle = false
			//	start_moving_to_new_order()
			//}
		}
		if handling_direction == 1 && local_queue.moving_direction = -1 {
			handling_direction = 1
			go_to_lowest_up_order(Broadcast_buffer)
		}else if handling_direction == 1 && local_queue.moving_direction = 1 {
			stop_at_all_up_orders(Broadcast_buffer)
			idle = true
		}else if handling_direction == -1 && local_queue.moving_direction = 1 {
			handling_direction = -1
			go_to_highest_down_order(Broadcast_buffer)
		}else if handling_direction == -1 && local_queue.moving_direction = -1 {
			stop_at_all_down_orders(Broadcast_buffer)
			idle = true
		}
	}
}
func Update_lights(External_order Types.Order_queue){
		for i:=0; i<3; i++ {
			if External_order.Mask.Outside_order_up[i]==1 {
				state:=External_order.Data.Outside_order_up[i]
				driver.Set_button_lamp(0,i,state)
			}
		}
		for i:=1; i<4; i++ {
			if External_order.Mask.Outside_order_down[i]==1 {
				state:=External_order.Data.Outside_order_down[i]
				driver.Set_button_lamp(1,i,state)
			}
		}
		for i:=0; i<4; i++ {
			if External_order.Mask.Inside_order[i]==1 {
				state:=External_order.Data.Inside_order[i]
				driver.Set_button_lamp(2,i,state)				
			}

		}
}
func Update_outside_order(External_order Types.Order_queue) {
		for i:=0; i<3; i++ {
			if External_order.Mask.Outside_order_up[i]==1 {
				local_queue.Outside_order_up[i]=1
			}
		}
		for i:=1; i<4; i++ {
			if External_order.Mask.Outside_order_down[i]==1 {
				local_queue.Outside_order_down[i]
			}
		}
}
func Update_inside_order(floor int,state int){
		local_queue.Inside_order[floor]=state
}

//Denne sjekker etter bestillinger på lokal heis, og sjekker bestillinger fra nettverket. 
func Get_orders(Received_order_ch chan Network.Message){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)
	//var outside_order_array [2]int
	//inside_order := int
	msg:=Network.Message{}
	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
		case outside_order_array := <- outside_order_ch:
			//Etasje blir aldri null for ned og 4 for oppbestillinger
			floor := outside_order_array[0]
			order_type := outside_order_array[1]
			if isAlreadyinQueue(floor, order_type)==false {
					Local_order,New_Queue=assign_order(floor,order_type)
					if Local_order {
						Update_outside_order(New_Queue)
						Update_lights(New_Queue)
					}
					Broadcast_buffer<-New_Queue
			}
		case inside_order := <- inside_order_ch:
			floor := inside_order
			order_type := 2
			if isAlreadyinQueue(floor,order_type)==false {
					Update_inside_order(floor,1)
			}
		
		case External_order:=<-Received_order_ch:
			if External_order.message_type=="Update" {
				Update_lights(External_order)
			}else if External_order.message_type=="New_order" {
				Update_outside_order(External_order)
			}
		}
	}
}
func Send_orders(Broadcast_order chan Network.Message,Broadcast_buffer chan Network.Message){
	for {
		select {
			case msg:=<-Broadcast_buffer:
				Broadcast_order<-msg
		}
	}
}
func Calculate_orders_amount(Queue Types.Order_queue)(int){
		amount:=0
		for i:=0; i<3; i++ {
			if Queue.Data.Outside_order_up[i]==1 {
				amount=amount+1
			}
		}
		for i:=1; i<4; i++ {
			if Queue.Data.Outside_order_down[i]==1 {
				amount=amount+1
			}
		}
		for i:=0; i<4; i++ {
			if Queue.Data.Inside_order[i]==1 {
				amount=amount+1				

			}
		}
		return amount
}
func calculate_order_cost(floor int, order_type int,Queue_n Types.Order_queue)(int){
		Direction:=Queue_n.Data.Moving_Direction
		LastFloor:=Queue_n.Data.Lastfloor
		Moving:=Queue_n.Data.Moving
		cost:=0
		if LastFloor==floor&&Moving==false {
			return -1
		}

		if order_type==Driver.BUTTON_CALL_UP&&Direction==DIRN_UP { //Riktig retning{
			if LastFloor>floor { //Kjørt forbi etasje
				cost=cost+5
			}
		}else if order_type==Driver.BUTTON_CALL_DOWN&&Direction==DIRN_DOWN {
			if LastFloor<floor {
				cost=cost+5
			}
		}
		if order_type==Driver.BUTTON_CALL_UP&&Direction==DIRN_DOWN {
			cost=cost+(Types.N_FLOORS-LastFloor)+(Types.N_FLOORS-floor)
		}else if order_type==Driver.BUTTON_CALL_DOWN&&Direction==DIRN_UP {
			cost=cost+(LastFloor)+floor
		}
		cost=cost+Calculate_orders_amount(Queue_n)
		return cost

}
func assign_order(floor int,order_type int)(bool,Types.Order_queue){
	temp_Queue_Network:=Network.Queue_Network
	var lowest_cost int
	var Local bool
	cost:=0
	var lowest_cost_ipAddr string
	var Queue_update Network.Message


	lowest_cost=calculate_order_cost(floor,order_type,local_queue)
	lowest_cost_ipAddr="Local"
	Local=true
	for ipAddr,Queue_n:=range temp_Queue_Network {
		cost=calculate_order_cost(floor,order_type,Queue_n)
		if cost<lowest_cost {
			lowest_cost=cost
			lowest_cost_ipAddr=ipAddr
			Local=false
		}
	}

	if order_type==BUTTON_CALL_UP {
		Queue_update.Mask.Outside_order_up[floor]
	}else{
		Queue_update.Mask.Outside_order_down[floor]
	}

	if lowest_cost_ipAddr=="Local" {
		Queue_update.MessageType="Update"
	}else {
		Queue_update.MessageType="New order"
		Queue_update.RecipientAddr=lowest_cost_ipAddr
	}
	
	return Local,Queue_update
}

//Denne endres til å ta inn floor og order_type istedenfor order_array
//Føler at låsen her er ubrukelig
func isAlreadyinQueue(floor int, order_type int) (bool) {
	temp_Queue_Network:=Network.Queue_Network
	for ipAddr,Queue_n := range temp_Queue_Network {
		if floor!=3&&Queue_n.Outside_order_up[floor]==1 {
			return true //Not finished
		}else if floor!=0&&Queue_n.Outside_order_down[floor]==1 {
			return true
		}else if Queue_n.Inside_order[floor]==1 {
			return true
		}
	}
	return false

}


func Redistribute_orders_in_new_network(){

}*/
