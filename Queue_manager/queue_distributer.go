package Queue_manager

import "../driver"
import "../Network"
import "../Types"
import "../Elevator"
import "time"
import "fmt"
import "math"



func Queue_manager_init(stop_chan chan int){
	
	current_floor:=driver.Driver_init()
	//Initialize variables
	fmt.Println("___________---Initialized---__________")

	Broadcast_order:=make(chan Types.Message,1)
	local_queue_chan:=make(chan Types.Order_queue,1)
	Queue_Network_chan:=make(chan map[string]Types.Order_queue,1) 
	Received_order_ch:=make(chan Types.Message,1)
	Broadcast_buffer:=make(chan Types.Message,500)
	Delete_order_ch:=make(chan Types.Order_queue,1)
	Activity_timer_ch:=make(chan time.Time,1)
	
	var local_queue Types.Order_queue;
	Activity_timer:=time.Now()
	Activity_timer_ch<-Activity_timer
	local_queue.Last_floor=current_floor
	local_queue.Moving=false
	local_queue_chan<-local_queue

	go Network.Network_Manager_init("20020",Broadcast_order,Received_order_ch,stop_chan,Queue_Network_chan)
	go Elevator.Elevator_init(Broadcast_buffer,current_floor,local_queue_chan,Delete_order_ch)
	go Get_orders(Received_order_ch,Broadcast_buffer,Queue_Network_chan,local_queue_chan,stop_chan,Delete_order_ch,Activity_timer_ch)
	go Send_orders(Broadcast_order,Broadcast_buffer)
	go Activity_monitor(Activity_timer_ch,local_queue_chan,Broadcast_buffer)
	
}
func Cleanup(){
	fmt.Println(" ")
	fmt.Println("Exiting elevator")
	driver.Driver_init()
	fmt.Println("Elevator in defined state")
	fmt.Println(".............................")
	fmt.Println("Goood byee")
	
}

func Update_lights(inside_order bool,Queue_mask Types.Order_queue,Queue_data Types.Order_queue,local_queue_chan chan Types.Order_queue){
		local_queue:=<-local_queue_chan
		local_queue_chan<-local_queue
		for i:=0; i<Types.N_FLOORS-1; i++ {
			if Queue_mask.Outside_order_up[i]==1 {
				state:=Queue_data.Outside_order_up[i]
				if (state==0 && local_queue.Outside_order_up[i]==1)==false {
					driver.Set_button_lamp(0,i,state)
				}
			}
		}
		for i:=1; i<Types.N_FLOORS; i++ {
			if Queue_mask.Outside_order_down[i]==1 {
				state:=Queue_data.Outside_order_down[i]
				if (state==0 && local_queue.Outside_order_down[i]==1)==false {
					driver.Set_button_lamp(1,i,state)
				}
			}
		}
		if inside_order==true {
			for i:=0; i<Types.N_FLOORS; i++ {
				if Queue_mask.Inside_order[i]==1 {
					state:=Queue_data.Inside_order[i]
					driver.Set_button_lamp(2,i,state)
				}
			}
		}		
}
func Update_queue(inside_order bool,state int,Order_mask Types.Order_queue, Broadcast_buffer chan Types.Message,local_queue_chan chan Types.Order_queue) {
		local_queue:=<-local_queue_chan
		if inside_order==true {
			for i:=0; i<Types.N_FLOORS; i++ {
				if Order_mask.Inside_order[i]==1 {
					local_queue.Inside_order[i]=state
				}
			}
		}
		for i:=0; i<Types.N_FLOORS-1; i++ {
			if Order_mask.Outside_order_up[i]==1 {
				local_queue.Outside_order_up[i]=state
			}
		}
		for i:=1; i<Types.N_FLOORS; i++ {
			if Order_mask.Outside_order_down[i]==1 {
				local_queue.Outside_order_down[i]=state
			}
		}

		temp_local_queue:=local_queue
		local_queue_chan<-local_queue

		Update_lights(inside_order,Order_mask,temp_local_queue,local_queue_chan)
		
		var msg Types.Message
		msg.Message_type=Types.MT_update
		msg.Data=local_queue
		msg.Mask=Order_mask
		Broadcast_buffer<-msg
		
}

func Get_orders(Received_order_ch,Broadcast_buffer chan Types.Message,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue,stop_chan chan int,Delete_order_ch chan Types.Order_queue,Activity_timer_ch chan time.Time){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)

	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
			case outside_order_array := <- outside_order_ch:
				floor := outside_order_array[0]
				order_type := outside_order_array[1]
				if is_already_in_queue(floor, order_type,Queue_Network_chan,local_queue_chan)==false {
						Activity_timer:=<-Activity_timer_ch
						Activity_timer=time.Now()
						Activity_timer_ch<-Activity_timer						
						Local_order,New_Queue:=assign_order(floor,order_type,Queue_Network_chan,local_queue_chan)
						if Local_order{
							Update_queue(false,1,New_Queue.Mask,Broadcast_buffer,local_queue_chan)		
						}else {
							fmt.Println("Sending order to",New_Queue.Recipient_addr)
							Broadcast_buffer<-New_Queue
						}
				}
			case inside_order := <- inside_order_ch:
				var Order_mask Types.Order_queue
				floor := inside_order
				order_type := 2
				Order_mask.Inside_order[floor]=1
				if is_already_in_queue(floor,order_type,Queue_Network_chan,local_queue_chan)==false {
						Activity_timer:=<-Activity_timer_ch
						Activity_timer=time.Now()
						Activity_timer_ch<-Activity_timer
						Update_queue(true,1,Order_mask,Broadcast_buffer,local_queue_chan)
				}
		
			case External_order:=<-Received_order_ch:
				if External_order.Message_type==Types.MT_update { //Update
					Update_lights(false,External_order.Mask,External_order.Data,local_queue_chan)
				}else if External_order.Message_type==Types.MT_new_order { //New order {
					Update_queue(false,1,External_order.Mask,Broadcast_buffer,local_queue_chan)
				}else if External_order.Message_type==Types.MT_disconnected { //Disconnected {
					Redistribute_orders(External_order,Broadcast_buffer,Queue_Network_chan,local_queue_chan)
				}else if External_order.Message_type==Types.MT_new_elevator { //New elevator {
					msg:=Types.Message{Message_type:Types.MT_update}
					local_queue:=<-local_queue_chan
					local_queue_chan<-local_queue
					msg.Data=local_queue
					msg.Mask=local_queue
					Broadcast_buffer<-msg							
				}else if External_order.Message_type==Types.MT_out {
					Redistribute_orders(External_order,Broadcast_buffer,Queue_Network_chan,local_queue_chan)
					
				}
			case <-stop_chan:
				<-local_queue_chan
			case Mask:=<-Delete_order_ch:
				Activity_timer:=<-Activity_timer_ch
				Activity_timer=time.Now()
				Activity_timer_ch<-Activity_timer
				Update_queue(true,0,Mask,Broadcast_buffer,local_queue_chan)
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
		Last_floor:=Queue_n.Last_floor
		cost:=0
		if Last_floor==floor && Direction==order_type  {
			cost=-2
		}
		
		if order_type==Types.BUTTON_CALL_UP&&Direction==Types.DIRN_UP {
			if Last_floor>floor { 
				cost=cost+5
			}else {
				cost=cost+floor-Last_floor
			
			}
		}else if order_type==Types.BUTTON_CALL_DOWN&&Direction==Types.DIRN_DOWN {
			if Last_floor<floor {
				cost=cost+5
			}else {
				cost=cost+Last_floor-floor
			}
		}
		if order_type==Types.BUTTON_CALL_UP&&Direction==Types.DIRN_DOWN {
			cost=cost+(Last_floor)-floor
		}else if order_type==Types.BUTTON_CALL_DOWN&&Direction==Types.DIRN_UP {
			cost=cost+Last_floor-floor
		}
		cost=cost+Calculate_orders_amount(Queue_n)+int(math.Abs(float64(Last_floor-floor)))
		return cost

}
func assign_order(floor int,order_type int,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue)(bool,Types.Message){
	temp_Queue_Network:=<-Queue_Network_chan
	Queue_Network_chan<-temp_Queue_Network
	var lowest_cost,cost int
	var Local bool
	var lowest_cost_ipAddr string
	var Queue_update Types.Message
	
	cost=0

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
		Queue_update.Message_type=Types.MT_update
	}else {
		Queue_update.Message_type=Types.MT_new_order
		Queue_update.Recipient_addr=lowest_cost_ipAddr
	}
	
	return Local,Queue_update
}


func is_already_in_queue(floor int, order_type int,Queue_Network_chan chan map[string]Types.Order_queue,local_queue_chan chan Types.Order_queue) (bool) {
	temp_Queue_Network:=<-Queue_Network_chan
	Queue_Network_chan<-temp_Queue_Network
	for _,Queue_n := range temp_Queue_Network {
		if order_type==Types.BUTTON_CALL_UP&&floor!=3&&Queue_n.Outside_order_up[floor]==1 {
			return true 
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
				Update_queue(false,1,New_Queue.Mask,Broadcast_buffer,local_queue_chan)
			}else {
				Broadcast_buffer<-New_Queue
			}
		}
	}
	for floor:=1; floor<Types.N_FLOORS; floor++ {
		if Queue.Outside_order_down[floor]==1 {
			Local_order,New_Queue:=assign_order(floor,Types.BUTTON_CALL_DOWN,Queue_Network_chan,local_queue_chan)
			if Local_order {
				Update_queue(false,1,New_Queue.Mask,Broadcast_buffer,local_queue_chan)
			}else {
				Broadcast_buffer<-New_Queue
			}
		}
	}
}
func Activity_monitor(Activity_timer_ch chan time.Time,local_queue_chan chan Types.Order_queue,Broadcast_buffer_ch chan Types.Message){
   msg:=Types.Message{Message_type: Types.MT_out}
   for {
   		Activity_timer:=<-Activity_timer_ch
   		Activity_timer_ch<-Activity_timer
   		if time.Now().Sub(Activity_timer)>10*time.Second && Elevator.Is_queue_empty(local_queue_chan)==false {
   			//Send message
	   		Activity_timer:=<-Activity_timer_ch
			Activity_timer=time.Now()
	   		Activity_timer_ch<-Activity_timer
   			Broadcast_buffer_ch<-msg
			fmt.Println("I have some problems")
   		}
   		time.Sleep(5*time.Second)
   }

}
