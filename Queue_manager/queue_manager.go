package Queue_manager

import "../driver"
import "../Network"
import "../Types"
//import "../driver" //HOW DOES ONE DO THIS

var Queue Types.Queue_type;

func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	new_message:=make(chan Network.Message,1)
	Queue_chan:=make(chan int,1)
	elev_chan:=make(chan int,1)
	Order_update:=make(chan string,1)
	Queue_chan<-1
	go Network.Network_Manager("20020",Queue_chan,new_message,stop_chan,Order_update,elev_chan)
	pressed:=[]int{0,0,0,0}
	driver.Driver_init()
	msg:=Network.Message{MessageType: "Update",Data: Queue}
	for{
		select{

			case ip:=<-Order_update:
				<-elev_chan
				new_Queue:=Network.Queue_Network[ip];
				elev_chan<-1
				go Update_lights(new_Queue)
			default:
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
						msg.Data=Queue;
						new_message<-msg
						//fmt.Print("Button 2, ",i," is pressed!\r")	
					}
			}
		}	



		}
		stop_chan<-1			//fmt.Println("On default")
}
func Update_lights(Queue_ip Types.Queue_type){
		for i:=0; i<4; i++ {
			if Queue_ip[i]==true  {
				driver.Set_button_lamp(2,i,1)
			}else {
				driver.Set_button_lamp(2,i,0)
			}
		
		}

}
/*
func Get_orders(){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)
	var outside_order_array [2]int
	inside_order := int
	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
		case outside_order_array = <- outside_order_ch:
			if "order already in queue"{ //YES NOT DONE
				continue
			}
			else{
				Send_order_to_network()
				"Add order to queue"
			}
		case inside_order = <- inside_order_ch:
			if "order already in queue"{ //YES NOT DONE
				continue
			}
			else{
				"Add order to queue"
			}
		}
	}
}

func Send_order_to_network(){

}

func Receive_order_from_network(){

}*/
