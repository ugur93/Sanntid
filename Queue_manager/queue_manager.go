package Queue_manager

import "../driver"
import "../Network"
import "../Types"
import "time"
//import "../driver" //HOW DOES ONE DO THIS

var Queue Types.Queue_type;

func Queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);

	//Stian:
	//Når det skjer en oppdatering på Queue matrise send gjennom new_message channel
	//Når det skjer en oppdatering på Queue fra andre heiser kommer det innpå Queue chan
	new_message:=make(chan Network.Message,1)
	Queue_chan:=make(chan int,1) //Bruk det når det leses på Queue
	elev_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Order_update:=make(chan string,1)
	Queue_chan<-1
	go Network.Network_Manager("20020",new_message,stop_chan,Order_update,elev_chan)
	pressed:=[]int{0,0,0,0}
	pressed_UP:=[]int{0,0,0,0}
	pressed_DOWN:=[]int{0,0,0,0}
	updated:=0
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
		for i:=4; i<7; i++ {
			if Queue_ip[i]==true  {
				driver.Set_button_lamp(0,i-4,1)
			}else {
				driver.Set_button_lamp(0,i-4,0)
			}
		}
		for i:=7; i<9; i++ {
			if Queue_ip[i]==true  {
				driver.Set_button_lamp(1,i-7,1)
			}else {
				driver.Set_button_lamp(1,i-7,0)
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
