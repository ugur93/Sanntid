package Queue_manager

import "../driver"
import "../Network"
import "../Types"
//import "time"
import "fmt"



func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	
	current_floor:=driver.Driver_init()

	//Initialize variables

	Broadcast_order:=make(chan Types.Message,1)
	local_queue_chan:=make(chan Types.Order_queue,1)
	Queue_Network_chan:=make(chan map[string]Types.Order_queue,1) //Bruk det når det leses på Queue
	Received_order_ch:=make(chan Types.Message,1)
	Broadcast_buffer:=make(chan Types.Message,500)

	var local_queue Types.Order_queue;
	local_queue.LastFloor=current_floor
	local_queue.Moving=false
	local_queue_chan<-local_queue
	//Queue_lock_chan<-1

	go Network.Network_Manager_init("20020",Broadcast_order,Received_order_ch,stop_chan,Queue_Network_chan)
	go Get_orders(Received_order_ch,Broadcast_buffer,Queue_Network_chan,local_queue_chan)
	go Send_orders(Broadcast_order,Broadcast_buffer)
	go Run_elevator(Broadcast_buffer,current_floor,local_queue_chan)

	stop_chan<-1
}



func Update_outside_order(External_order Types.Message, Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue) {
		local_queue:=<-local_queue_chan
		for i:=0; i<Types.N_FLOORS-1; i++ {
			if External_order.Mask.Outside_order_up[i]==1 {
				local_queue.Outside_order_up[i]=1
			}
		}
		for i:=1; i<Types.N_FLOORS; i++ {
			if External_order.Mask.Outside_order_down[i]==1 {
				local_queue.Outside_order_down[i]=1
			}
		}

		External_order.Data=local_queue
		local_queue_chan<-local_queue

		Update_outside_lights(External_order)
		//Generate Update message Fiks dette med egen funksjon
		var msg Types.Message
		msg.MessageType=Types.MT_update
		msg.Data=local_queue
		msg.Mask=External_order.Mask
		Broadcast_buffer<-msg
}
func Update_inside_order(floor int,state int,Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue){
		
		local_queue:=<-local_queue_chan
		local_queue.Inside_order[floor]=state
		Update_lights(2,floor,state)
		local_queue_chan<-local_queue
		
		var msg Types.Message
		msg.MessageType=Types.MT_update
		msg.Data=local_queue
		Broadcast_buffer<-msg
}

//Denne sjekker etter bestillinger på lokal heis, og sjekker bestillinger fra nettverket. 
func Get_orders(Received_order_ch,Broadcast_buffer chan Types.Message,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)
	//var outside_order_array [2]int
	//inside_order := int
	//var msg Network.Message
	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
			case outside_order_array := <- outside_order_ch:
				//Etasje blir aldri null for ned og 4 for oppbestillinger
				floor := outside_order_array[0]
				order_type := outside_order_array[1]
				if isAlreadyinQueue(floor, order_type,Queue_Network_chan,local_queue_chan)==false {
						Local_order,New_Queue:=assign_order(floor,order_type,Queue_Network_chan,local_queue_chan)
						if Local_order{
							Update_outside_order(New_Queue,Broadcast_buffer,local_queue_chan)		
						}else {
							fmt.Println("Sending order to",New_Queue.RecipientAddr)
							Broadcast_buffer<-New_Queue
						}
				}
			case inside_order := <- inside_order_ch:
				floor := inside_order
				order_type := 2
				if isAlreadyinQueue(floor,order_type,Queue_Network_chan,local_queue_chan)==false {
						Update_inside_order(floor,1,Broadcast_buffer,local_queue_chan)
				}
		
			case External_order:=<-Received_order_ch:
				if External_order.MessageType==Types.MT_update { //Update 
					Update_outside_lights(External_order)
				}else if External_order.MessageType==Types.MT_new_order { //New order {
					Update_outside_order(External_order,Broadcast_buffer,local_queue_chan)
				}else if External_order.MessageType==Types.MT_disconnected { //Disconnected {
					go Redistribute_orders(External_order,Broadcast_buffer,Queue_Network_chan,local_queue_chan)
				}else if External_order.MessageType==Types.MT_new_elevator { //New elevator {
					Update_outside_lights(External_order)				
				}
			}
	}
}
func Send_orders(Broadcast_order chan Types.Message,Broadcast_buffer chan Types.Message){
	for {
		select {
			case msg:=<-Broadcast_buffer:
				Broadcast_order<-msg
		}
	}
}
func Calculate_orders_amount(Queue Types.Order_queue)(int){
		amount:=0
		for i:=0; i<Types.N_FLOORS-1; i++ {
			if Queue.Outside_order_up[i]==1 {
				amount=amount+1
			}
		}
		for i:=1; i<Types.N_FLOORS; i++ {
			if Queue.Outside_order_down[i]==1 {
				amount=amount+1
			}
		}
		for i:=0; i<Types.N_FLOORS; i++ {
			if Queue.Inside_order[i]==1 {
				amount=amount+1				

			}
		}
		return amount
}
func calculate_order_cost(floor int, order_type int,Queue_n Types.Order_queue)(int){
		Direction:=Queue_n.Moving_direction
		LastFloor:=Queue_n.LastFloor
		//Moving:=Queue_n.Moving
		cost:=0
		if LastFloor==floor  {
			cost=-1
		}

		if order_type==Types.BUTTON_CALL_UP&&Direction==Types.DIRN_UP { //Riktig retning{
			if LastFloor>floor { //Kjørt forbi etasje
				cost=cost+5
			}else {
				cost=cost+floor-LastFloor
			
			}
		}else if order_type==Types.BUTTON_CALL_DOWN&&Direction==Types.DIRN_DOWN {
			if LastFloor<floor {
				cost=cost+5
			}else {
				cost=cost+LastFloor-floor
			}
		}
		if order_type==Types.BUTTON_CALL_UP&&Direction==Types.DIRN_DOWN {
			cost=cost+(Types.N_FLOORS-LastFloor)+floor
		}else if order_type==Types.BUTTON_CALL_DOWN&&Direction==Types.DIRN_UP {
			cost=cost+LastFloor+floor
		}
		cost=cost+Calculate_orders_amount(Queue_n)
		//fmt.Println("cost: ",cost)
		return cost

}
func assign_order(floor int,order_type int,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue)(bool,Types.Message){
	temp_Queue_Network:=<-Queue_Network_chan
	Queue_Network_chan<-temp_Queue_Network
	var lowest_cost int
	var Local bool
	cost:=0
	var lowest_cost_ipAddr string
	var Queue_update Types.Message

	local_queue:=<-local_queue_chan
	local_queue_chan<-local_queue

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

	if order_type==Types.BUTTON_CALL_UP {
		Queue_update.Mask.Outside_order_up[floor]=1
	}else if order_type==Types.BUTTON_CALL_DOWN {
		Queue_update.Mask.Outside_order_down[floor]=1
	}

	if Local {
		Queue_update.MessageType=Types.MT_update
	}else {
		Queue_update.MessageType=Types.MT_new_order
		Queue_update.RecipientAddr=lowest_cost_ipAddr
	}
	
	return Local,Queue_update
}

//Denne endres til å ta inn floor og order_type istedenfor order_array
//Føler at låsen her er ubrukelig
func isAlreadyinQueue(floor int, order_type int,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue) (bool) {
	temp_Queue_Network:=<-Queue_Network_chan
	Queue_Network_chan<-temp_Queue_Network
	for _,Queue_n := range temp_Queue_Network {
		if order_type==Types.BUTTON_CALL_UP&&floor!=3&&Queue_n.Outside_order_up[floor]==1 {
			return true //Not finished
		}else if order_type==Types.BUTTON_CALL_DOWN&&floor!=0&&Queue_n.Outside_order_down[floor]==1 {
			return true
		}
	}
	local_queue:=<-local_queue_chan
	local_queue_chan<-local_queue
	if order_type==Types.BUTTON_COMMAND&&local_queue.Inside_order[floor]==1 {
		return true
	}else if order_type==Types.BUTTON_CALL_UP&&floor!=3&&local_queue.Outside_order_up[floor]==1 {
		return true
	}else if order_type==Types.BUTTON_CALL_DOWN&&floor!=0&&local_queue.Outside_order_down[floor]==1 {
		return true
	}
	return false

}


func Redistribute_orders(External_order Types.Message,Broadcast_buffer chan Types.Message,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue){
	Queue:=External_order.Data
	for floor:=0; floor<Types.N_FLOORS-1; floor++ {
		if Queue.Outside_order_up[floor]==1 {
			Local_order,New_Queue:=assign_order(floor,Types.BUTTON_CALL_UP,Queue_Network_chan,local_queue_chan)
			if Local_order {
				Update_outside_order(New_Queue,Broadcast_buffer,local_queue_chan)
			}else {
				Broadcast_buffer<-New_Queue
			}
		}
	}
	for floor:=1; floor<Types.N_FLOORS; floor++ {
		if Queue.Outside_order_down[floor]==1 {
			Local_order,New_Queue:=assign_order(floor,Types.BUTTON_CALL_DOWN,Queue_Network_chan,local_queue_chan)
			if Local_order {
				Update_outside_order(New_Queue,Broadcast_buffer,local_queue_chan)
			}else {
				Broadcast_buffer<-New_Queue
			}
		}
	}
}
