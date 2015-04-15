package Queue_manager

import "../driver"
import "../Network"
import "../Types"
import "time"
//import "../driver" //HOW DOES ONE DO THIS

var local_queue Types.Order_queue;

func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	current_floor:=driver.Driver_init()
	local_queue.LastFloor=current_floor
	//Stian:
	//Når det skjer en oppdatering på Queue matrise send gjennom new_message channel
	//Når det skjer en oppdatering på Queue fra andre heiser kommer det innpå Queue chan
	Broadcast_order:=make(chan Network.Message,1)
	Queue_lock_chan:=make(chan int,1) //Bruk det når det leses på Queue
	Queue_Network_lock_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Received_order_ch:=make(chan Network.Message,1)
	Broadcast_buffer:=make(chan Network.Message,500)
	Queue_chan<-1
	go Network.Network_Manager("20020",Broadcast_order,stop_chan,Received_order_ch,Queue_Network_lock_chan)
	go Get_orders(Order_update)
	go Send_orders(Broadcast_order,Broadcast_buffer)
	go Run_elevator(Queue_Network_lock_chan)
	stop_chan<-1
}

func Update_lights(External_order Types.Order_queue){
		for i:=0; i<3; i++ {
			if External_order.Mask.outside_order_up[i]==1 {
				state:=External_order.Data.outside_order_up[i]
				driver.Set_button_lamp(0,i,state)
			}
		}
		for i:=1; i<4; i++ {
			if External_order.Mask.outside_order_down[i]==1
				state:=External_order.Data.outside_order_down[i]
				driver.Set_button_lamp(1,i,state)
		}
		for i:=0; i<4; i++ {
			if External_order.Mask.inside_order[i]==1 {
				state:=External_order.Data.inside_order[i]
				driver.Set_button_lamp(2,i,state)				

			}

		}
}
func Update_outside_order(External_order Types.Order_queue) {
		for i:=0; i<3; i++ {
			if External_order.Mask.outside_order_up[i]==1 {
				local_queue.outside_order_up[i]=1
			}
		}
		for i:=1; i<4; i++ {
			if External_order.Mask.outside_order_down[i]==1 {
				local_queue.outside_order_down[i]
			}
		}
}
func Update_inside_order(floor,state){
		local_queue.inside_order[floor]=state
}

func Run_elevator(Queue_Network_lock_chan chan int){
	


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
		}
		case External_order:=<-Received_order_ch:
			if External_order.message_type=="Update" {
				Update_lights(External_order)
			}else if External_order.message_type=="New_order" {
				Update_outside_order(External_order)
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
			if Queue.Data.outside_order_up[i]==1 {
				amount=amount+1
			}
		}
		for i:=1; i<4; i++ {
			if Queue.Data.outside_order_down[i]==1 {
				amount=amount+1
			}
		}
		for i:=0; i<4; i++ {
			if Queue.Data.inside_order[i]==1 {
				amount=amount+1				

			}
		}
		return amount
}
func calculate_order_cost(floor int, order_type int,Queue_n Types.Order_queue)(int){
		Direction:=Queue_n.Data.direction
		LastFloor:=Queue_n.Data.lastfloor
		Moving:=Queue_n.Data.moving
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
	lowest_cost int
	Local bool
	cost:=0
	lowest_cost_ipAddr string
	Queue_update Types.Order_queue


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
		Queue_update.Mask.outside_order_up[floor]
	}else{
		Queue_update.Mask.outside_order_down[floor]
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
		if floor!=3&&Queue_n.outside_order_up[floor]==1 {
			return true //Not finished
		}else if floor!=0&&Queue_n.outside_order_down[floor]==1 {
			return true
		}else if Queue_n.inside_order[floor]==1 {
			return true
		}
	}
	return false

}


func Redistribute_orders_in_new_network(){

}
