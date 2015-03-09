package Queue_manager

import "../Network"
import "../Types"

func Test_queue_manager_init(){
	new_message:=make(chan Network.Message,1)
	Queue_chan:=make(chan int,1) //Bruk det når det leses på Queue
	elev_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Order_update:=make(chan string,1)
	Queue_chan<-1
	go Network.Network_Manager("20020",new_message,stop_chan,Order_update,elev_chan)
	
	
}