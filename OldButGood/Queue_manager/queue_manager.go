package Queue_manager
/*
import "../driver"
import "../Network"
import "../Types"
import "time"
//import "../driver" //HOW DOES ONE DO THIS

var Queue Types.Order_queue;
var Global_Queue Types.Order_queue;
func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	driver.Driver_init()
	//Stian:
	//Når det skjer en oppdatering på Queue matrise send gjennom new_message channel
	//Når det skjer en oppdatering på Queue fra andre heiser kommer det innpå Queue chan
	Broadcast_order:=make(chan Network.Message,1)
	Queue_lock_chan:=make(chan int,1) //Bruk det når det leses på Queue
	Queue_Network_lock_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Order_update:=make(chan Network.Message,1)
	Queue_chan<-1
	go Network.Network_Manager("20020",Broadcast_order,stop_chan,Order_update,Queue_Network_lock_chan)
	go Get_orders(Queue_Network_lock_chan,Broadcast_order,Order_update)
	go Run_elevator(Queue_Network_lock_chan)
	stop_chan<-1
}

func Update_lights(Queue_ip Types.Order_queue){
		for i:=0; i<4; i++ {
			if Queue_ip[i]==1  {
				driver.Set_button_lamp(2,i,1)
			}else {
				driver.Set_button_lamp(2,i,0)
			}
		
		}
		for i:=4; i<7; i++ {
			if Queue_ip[i]==1  {
				driver.Set_button_lamp(0,i-4,1)
			}else {
				driver.Set_button_lamp(0,i-4,0)
			}
		}
		for i:=7; i<9; i++ {
			if Queue_ip[i]==1 {
				driver.Set_b utton_lamp(1,i-7,1)
			}else {
				driver.Set_button_lamp(1,i-7,0)
			}
		}

}
func Run_elevator(Queue_Network_lock_chan chan int){
	


}
//Denne sjekker etter bestillinger på lokal heis, og sjekker bestillinger fra nettverket. 
func Get_orders(Queue_Network_lock_chan chan int, Broadcast_order chan Network.Message,Order_external chan Network.Message){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)
	//var outside_order_array [2]int
	//inside_order := int
	msg:=Network.Message{MessageType: "New order"}
	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
		case outside_order_array := <- outside_order_ch:
			floor := outside_order_array[0]
			order_type := outside_order_array[1]
			if isAlreadyinQueue(floor, order_type, Queue_Network_lock_chan) {
				continue
			}
			else{
				msg.Mask[(floor-1) + 3*order_type)]=1
				msg.MessageType="New order"
				Broadcast_order<-msg
				calculate_order_cost(Queue_Network_lock_chan, )
			}
		case inside_order := <- inside_order_ch:
			order_array := []int{inside_order,2}
			if isAlreadyinQueue(order_array,Queue_Network_lock_chan) {
				continue
			}
			else{
				floor := inside_order
				order_type := 2
				<-Queue_Network_lock_chan
				Queue[(floor-1) + 3*order_type)]=1
				Queue_Network_lock_chan<-1
				msg.Data=Queue
				msg.MessageType="Update"
				Broadcast_order<-msg
			}
		}
		case External_order:=<Order_external:
			Calculate_order_cost_queue(External_order) //Queue=External_order.Data Mask=External_order.Mask
		case something:=<-NetworkStateUpdate:
			Redistribute_orders_in_new_network(something)
	}
}

func calculate_order_cost_single_queue(floor int, queue Types.Order_queue){
	//First check if it can take order on this floor now
	//USES MANY IMAGINARY VARIABLES
	cost = 1 + calculate_order_cost_local_queue(next_floor(queue), queue)
	return cost
}

func next_floor(Queue_N Types.Order_queue)(next_floor int, Queue_N Types.Order_queue){
	if direction == Driver.DIRN_UP {
		for floor := current_floor ; floor < Types.N_FLOORS ; floor++ {
			if floor!=Types.N_FLOORS && Queue_N[floor-1]==1 {
				return floor 
			}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
			   	return floor
			}
		} 
	}
	else if direction == Driver.DIRN_DOWN {
		for floor:=current_floor; floor>0; floor-- {
			if floor!=Types.N_FLOORS && Queue_N[floor-1]==1 {
			   	return floor
			}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
			   	return floor
			}	
		}
		for floor:=0; floor<order_floor; floor++ {
			if floor!=1 && Queue_N[Types.N_FLOORS-1+floor-1]==1 {
				return floor
			}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
				return floor
			}	
		}	
	}
	else if order_type==Driver.BUTTON_CALL_DOWN {
		if Direction ==Driver.DIRN_DOWN {
			for floor:=Floor; floor<order_floor; floor-- {
				if floor!=1&&Queue_N[Types.N_FLOORS-1+floor-1]==1 {
					return floor
				}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
					return floor
				}
			}
		}
		else if Direction==Driver.DIRN_UP {
			for floor:=Floor; floor<Types.N_FLOOR; floor++ {
				if floor!=Types.N_FLOORS&&Queue_N[floor-1]==1 {
				   	return floor
				}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
				  	return floor
				}
			}
			for floor:=Types.N_FLOOR; floor>order_floor; floor-- {
				if floor!=Types.N_FLOORS&&Queue_N[floor-1]==1 {
					return floor
				}else if Queue_N[2*(Types.N_FLOORS-1)+floor-1]==1 {
					return floor
				}   	
			}	
			
		}
	}
}

func calculate_order_cost(Queue_lock_chan chan int,order_array chan [2]int)(string int) {
	temp_Queue_Network:=Network.GetQueueNetwork()
	order_floor := outside_order_array[0]
	order_type := outside_order_array[1]
	lowest_cost:=0
	cost:=0
	my_cost:=0
	lowest_cost_ipAddr string
	for ipAddr,Queue_N := range temp_Queue_Network {
		Direction:=Queue_N[Types.N_FLOORS+2*(Types.N_FLOORS-1)]
		Floor:=Queue_N[Types.N_FLOORS+2*(Types.N_FLOORS-1)+1]
		cost=0
		if order_type==Driver.BUTTON_CALL_UP {
			if DIRECTION == Driver.DIRN_UP {
			   for i:=Floor; i<order_floor; i++ {
			   			if i!=Types.N_FLOORS&&Queue_N[i-1]==1 {
			   				cost=cost+1
			   			}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
			   				cost=cost+1
			   			}
			   }
			}else if Direction==Driver.DIRN_DOWN {
				for i:=Floor; i>0; i-- {
					if i!=Types.N_FLOORS&&Queue_N[i-1]==1 {
			   				cost=cost+1
			   		}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
			   				cost=cost+1
			   		}	
			   	}
			   	for i:=0; i<order_floor; i++ {
			   		if i!=1&&Queue_N[Types.N_FLOORS-1+i-1]==1 {
			   				cost=cost+1
			   		}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
			   				cost=cost+1
			   		}	
			   	}	
			}
		
		}
		else if order_type==Driver.BUTTON_CALL_DOWN {
		
			if Direction ==Driver.DIRN_DOWN {
				for i:=Floor; i<order_floor; i-- {
				   			if i!=1&&Queue_N[Types.N_FLOORS-1+i-1]==1 {
				   				cost=cost+1
				   			}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
				   				cost=cost+1
				   			}
				   }
				}
			else if Direction==Driver.DIRN_UP {
					for i:=Floor; i<Types.N_FLOOR; i++ {
						if i!=Types.N_FLOORS&&Queue_N[i-1]==1 {
				   				cost=cost+1
				   		}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
				   				cost=cost+1
				   		}
				   	}
				   	for i:=Types.N_FLOOR; i>order_floor; i-- {
				   		if i!=Types.N_FLOORS&&Queue_N[i-1]==1 {
				   				cost=cost+1
				   		}else if Queue_N[2*(Types.N_FLOORS-1)+i-1]==1 {
				   				cost=cost+1
				   		}
				   	
				   	}	
			
			}
		}
		if cost < lowest_cost {
			lowest_cost = cost
			lowest_cost_ipAddr = ipAddr
		}
		
	
	}
	return lowest_cost_ipAddr,lowest_cost


}
//Denne endres til å ta inn floor og order_type istedenfor order_array
//Føler at låsen her er ubrukelig
func isAlreadyinQueue(floor int, order_type int) (bool) {
	temp_Queue_Network:=Network.Queue_Network
	for ipAddr,Queue_N := range temp_Queue_Network {
		if Queue_N.Outside_order_up[floor] {
			return true //Not finished
		}
	
	}

}


func Redistribute_orders_in_new_network(){

}
//func Receive_order_from_network(){

	for i:=0; i<4; i++ {
					if i!=3 {
						if driver.Get_button_signal(0,i)==0 {
							pressed_UP[i]=0;
						}
						if driver.Get_button_signal(0,i)==1 && pressed_UP[i]==0 {
							pressed_UP[i]=1;
							if Queue[i+4]==true {
								Queue[i+4]=false
								driver.Set_button_lamp(0,i,0)
							}else {
								Queue[i+4]=true;
								driver.Set_button_lamp(0,i,1)
							}
							updated=1
					
						}
					}
					if i!=0 {
						if driver.Get_button_signal(1,i)==0 {
							pressed_DOWN[i]=0;
						}
						if driver.Get_button_signal(1,i)==1 && pressed_DOWN[i]==0 {
							pressed_DOWN[i]=1;
							if Queue[i+7]==true {
								Queue[i+7]=false
								driver.Set_button_lamp(1,i,0)
							}else {
								Queue[i+7]=true;
								driver.Set_button_lamp(1,i,1)
							}
							updated=1
					
						}
					}		


				}
				for i:=0; i<4; i++ {
					if driver.Get_button_signal(2,i)==0 {
						pressed[i]=0;
					}
					if driver.Get_button_signal(2,i)==1 && pressed[i]==0 {
						pressed[i]=1;
						if Queue[i]==true {
							Queue[i]=false
							driver.Set_button_lamp(2,i,0)
						}else {
							Queue[i]=true;
							driver.Set_button_lamp(2,i,1)
						}
						updated=1
						//fmt.Print("Button 2, ",i," is pressed!\r")	
					}
				}
				if updated==1 {
					msg.Data=Queue;
					new_message<-msg
					updated=0

				}
			//Viktig med delay på alt av tilstandsjekk!!!!!!!!
			time.Sleep(100*time.Millisecond) 
		}	
//}*/
