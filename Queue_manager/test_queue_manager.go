package Queue_manager

import "../Network"
import "../Types"
import ("fmt"
		//"os"
		//"bufio"
		)


func Test_queue_manager_init(stop_chan chan int){
	new_message:=make(chan Network.Message,1)
	Queue_chan:=make(chan int,1) //Bruk det når det leses på Queue
	elev_chan:=make(chan int,1) //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Order_update:=make(chan Network.Message,1)
	Queue_chan<-1
	go Network.Network_Manager("1",new_message,stop_chan,Order_update,elev_chan)
	msg:=Network.Message{MessageType:"New order",Data: Queue}
	var Mask Types.Queue_type
	for {

		select {
			case data:=<-Order_update:
				fmt.Println(data.Mask)
				fmt.Println(data.Data)
			default: 
				var text int
				fmt.Scanf("%d",&text)
				switch text {
						case 1://fmt.Println(text)
							fmt.Println("sending: 1")
							Queue[0]=1
							msg.Data=Queue
							Mask[0]=1
							msg.Mask=Mask
							new_message<-msg
							fmt.Println("sent")
						case 2:
							fmt.Println("sending: 2")
							Queue[1]=1
							msg.Data=Queue
							Mask[1]=1
							msg.Mask=Mask
							new_message<-msg
							fmt.Println("sent")
						case 3:
							fmt.Println("sending: 3")
							Queue[2]=1
							msg.Data=Queue
							Mask[2]=1
							msg.Mask=Mask
							new_message<-msg
							fmt.Println("sent")
						case 4:
							fmt.Println("sending: 4")
							Queue[3]=1
							msg.Data=Queue
							Mask[3]=1
							msg.Mask=Mask
							new_message<-msg
							fmt.Println("sent")
						}
		}

	}

	stop_chan<-1

}