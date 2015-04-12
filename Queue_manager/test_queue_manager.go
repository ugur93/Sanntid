package Queue_manager

import "../Network"
import "../Types"
import "../driver"
import "time"



func Test_queue_manager_init(stop_chan chan int){
	//Queue_chan:=make(chan int);
	driver.Driver_init()
	//Stian:
	//Når det skjer en oppdatering på Queue matrise send gjennom new_message channel
	//Når det skjer en oppdatering på Queue fra andre heiser kommer det innpå Queue chan
	Broadcast_order:=make(chan Network.Message,1)
	Queue_lock_chan:=make(chan int,1) //Bruk det når det leses på Queue
	Queue_Network_lock_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Order_update:=make(chan Network.Message,1)
	Queue_lock_chan<-1
	go Network.Network_Manager_init("20020",Broadcast_order,stop_chan,Order_update,Queue_Network_lock_chan)
	msg:=Network.Message{MessageType:"New order"}
	var Queue Types.Queue_type
	pressed_UP:=[]int{0,0,0,0}
	pressed_DOWN:=[]int{0,0,0,0}
	pressed:=[]int{0,0,0,0}
	updated:=0
	for {
	
			for i:=1; i<=4; i++ {
					if i!=4 {
						if driver.Get_button_signal(0,i)==0 {
							pressed_UP[i-1]=0;
						}
						if driver.Get_button_signal(0,i)==1 && pressed_UP[i-1]==0 {
							pressed_UP[i-1]=1;
							if Queue[i+4-1]==1 {
								Queue[i+4-1]=0
								driver.Set_button_lamp(0,i,0)
							}else {
								Queue[i+4-1]=1;
								driver.Set_button_lamp(0,i,1)
							}
							updated=1
					
						}
					}
					if i!=1 {
						if driver.Get_button_signal(1,i)==0 {
							pressed_DOWN[i-1]=0;
						}
						if driver.Get_button_signal(1,i)==1 && pressed_DOWN[i-1]==0 {
							pressed_DOWN[i-1]=1;
							if Queue[i+7-1]==1 {
								Queue[i+7-1]=0
								driver.Set_button_lamp(1,i,0)
							}else {
								Queue[i+7-1]=1;
								driver.Set_button_lamp(1,i,1)
							}
							updated=1
					
						}
					}		


				}
				for i:=1; i<=4; i++ {
					if driver.Get_button_signal(2,i)==0 {
						pressed[i-1]=0;
					}
					if driver.Get_button_signal(2,i)==1 && pressed[i-1]==0 {
						pressed[i-1]=1;
						if Queue[i-1]==1 {
							Queue[i-1]=0
							driver.Set_button_lamp(2,i,0)
						}else {
							Queue[i-1]=1;
							driver.Set_button_lamp(2,i,1)
						}
						updated=1
						//fmt.Print("Button 2, ",i," is pressed!\r")	
					}
				}
				if updated==1 {
					msg.Data=Queue;
					Broadcast_order<-msg
					updated=0

				}
			//Viktig med delay på alt av tilstandsjekk!!!!!!!!
			time.Sleep(100*time.Millisecond) 
		}
	
	
	
	
	}
	


