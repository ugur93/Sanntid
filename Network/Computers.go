package Network

import(
	"fmt"
	"time"
	"../Types"


)

type ConnectionStatus struct {
	State string
	IpAddr string
	Queue Types.Order_queue
}

var network_TimeStamp=map[string]time.Time{} //private
var Queue_Network = map[string]Types.Order_queue{}
var numberOfElevators int
var ackFinished int

//Ip:: 118, 155, 145
func Network_Manager_init(Port string,Broadcast_message_ch chan Message,stop_chan chan int,New_Order chan Message,Queue_Network_lock_chan chan int){ 


	//initialize channels
	send_ch :=make(chan Message,1024)
	receive_ch :=make(chan Message,1024)
	time_chan:=make(chan int,1)
	ack_chan:=make(chan Message,1024)
	time_chan<-1
	Queue_Network_lock_chan<-1
	
	
	//Initialize variables
	ackFinished=1
	numberOfElevators=0
	
	//Calling functions
	UDP_init(Port,send_ch,receive_ch);
	go checkConnectionStatus(time_chan,Queue_Network_lock_chan)
	
	WaitForMessages(Queue_Network_lock_chan,time_chan,ack_chan,send_ch,receive_ch,Broadcast_message_ch)
	

}

func WaitForMessages(Queue_Network_lock_chan chan int, time_chan chan int, ack_chan,send_ch,receive_ch,Broadcast_message_ch chan Message) {

	for{
		select{
			case msg:=<-receive_ch:
				go HandleNewMessage(msg,Queue_Network_lock_chan,time_chan,ack_chan,send_ch,receive_ch)
			case msg:=<-Broadcast_message_ch:
				//Update This Computer (More for printing)
				Queue_Network["This Computer"]=msg.Data
				go printAllOrders(Queue_Network)
				
				go BroadcastMessage(msg,ack_chan,send_ch)
			}
			

	}



}

func HandleNewMessage(msg Message,Queue_Network_lock_chan chan int, time_chan chan int, ack_chan chan Message,send_ch, receive_ch chan Message) {
				
				ipAddr:=msg.RemoteAddr
				ack:=Message{MessageType:"ack"}
				
				//Update timestamp for the ipaddress
				<-time_chan
				network_TimeStamp[ipAddr]=time.Now()
				time_chan<-1
				
				if _,ok:=Queue_Network[ipAddr]; ok==false {
					
					//New Computer connected, put in map
					<-Queue_Network_lock_chan
					Queue_Network[ipAddr]=msg.Data
					Queue_Network_lock_chan<-1
					
					// Notify QueueManager with the Queue
					//Connection:=ConnectionStatus{State:"Connected",IpAddr:ipAddr,Queue:msg.Data}
					fmt.Println(ipAddr,"Connected to the network")
					numberOfElevators=numberOfElevators+1
				
				}else if msg.MessageType == "Update" {
				
						//Update Queue map
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						go printAllOrders(Queue_Network)
						Queue_Network_lock_chan<-1
						
						//Send ack
						ack.AckAddr=ipAddr
						send_ch<-ack
						//Notify QueueManager		
						//Send to QueueManager
						
				}else if msg.MessageType == "New order" {
				
						//Update Queue map
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						go printAllOrders(Queue_Network)
						Queue_Network_lock_chan<-1
						
						//New_Order<-msg
						
						//Send ack
						ack.AckAddr=ipAddr
						send_ch<-ack
						
				}else if msg.MessageType == "ack" {
						//Ack recieved, notify ack manager
						ack_chan<-msg

				}


}
func BroadcastMessage(msg Message,ack_chan,send_ch chan Message){
	
	
	//Wait for other acks to finish
	for {
		if ackFinished == 1 {
			ackFinished = 0
			break
		}
	}
	
	//Sending the message
	send_ch<-msg
	fmt.Println("Message sended,Waiting for ack")
	
	//Initalize variables
	TTR:=time.Now()
	Timeouts:=0
	finished:=0
	ackcounter:=0
	
	if numberOfElevators>0 {
		fmt.Println("In for loop",numberOfElevators)
		for {
				select {
					case ackmsg:=<-ack_chan:
						//Må være sikker på at vi mottar riktig ack
						if ackmsg.AckAddr==localAddr {
							ackcounter=ackcounter+1
							if ackcounter==numberOfElevators {
								fmt.Println("Recieved ack from all")
								finished=1
								break
						}
					}
					default:
						if time.Now().Sub(TTR)>200*time.Millisecond {
							Timeouts=Timeouts+1
							if Timeouts==10 {
								fmt.Println("Timed out too many times, something is wrong")
								finished=1
								break
							}
							//Resend the message
							fmt.Println("Ack timed out,resending")
							TTR=time.Now()
							send_ch<-msg
						}	
				}
				if finished==1 {
					break
				}
			}
		}
		ackFinished=1
}

func checkConnectionStatus(time_chan chan int,Queue_Network_lock_chan chan int){
	for{
			<-time_chan
			temp_TimeStamp:=network_TimeStamp
			time_chan<-1
			timeEnd:=time.Now()
			//Sjekker for alle heisene
			
			for ipAddr,timeStart:=range temp_TimeStamp {		
				if timeEnd.Sub(timeStart)>=300*time.Millisecond {
					//Computer Disconnected from network or not responding (loop?)
					//Comp_chan in use only if another computer is disconnected
					
					//Delete The computer from map
					<-Queue_Network_lock_chan
					//QueueData:=Queue_Network[ipAddr]
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)	
					Queue_Network_lock_chan<-1
					
					
					//Slett fra Computers arrayet
					//Connection:=ConnectionStatus{State:"Disconnected",Queue:QueueData}
					fmt.Println(ipAddr,"Disconnected from the network")
					numberOfElevators=numberOfElevators-1
				}
			}
			
			
		time.Sleep(100*time.Millisecond)

	}


}

func printAllOrders(temp_Queue_Network map[string]Types.Order_queue){
	fmt.Print("\033c")
	fmt.Println("------------------------List of Elevator Queues-------------------------------")
	fmt.Println("ipAddr                 | O1 | O2 | O3 | N1 | N2 | N3 | 1  | 2  | 3  | 4  | Dir | F |")
	fmt.Println("--------------------------------------------------------------------------------------")
	for ipAddr,Queue:=range temp_Queue_Network {
		fmt.Print(ipAddr+"    ")
		for _,k:=range Queue {
			fmt.Print(k)
			fmt.Print("    ")
		}
		fmt.Println("|")
		fmt.Println("-------------------------------------------------------------------------------------")		
	}

}
