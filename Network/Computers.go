package Network

import(
	"fmt"
	"time"
	"../Types"
	"strconv"
	"strings"


)

var numberOfElevators int
//var my_queue Types.Order_queue
//Ip:: 118, 155, 145.146,154,141
func Network_Manager_init(Port string,BroadcastOrderCh,ReceiveOrderCh chan Types.Message,stop_chan chan int,Queue_Network_chan chan map[string]Types.Order_queue){ 

	//Initialize variables
	numberOfElevators=0
	var Queue_Network = map[string]Types.Order_queue{}
	var network_TimeStamp=map[string]time.Time{}

	//initialize channels
	send_ch :=make(chan Types.Message,1024)
	receive_ch :=make(chan Types.Message,1024)
	timeStamp_chan:=make(chan map[string]time.Time,1)
	ack_chan:=make(chan Types.Message,1024)
	ack_lock_chan:=make(chan int,1)

	timeStamp_chan<-network_TimeStamp
	Queue_Network_chan<-Queue_Network
	ack_lock_chan<-1
	
	

	
	//Calling functions
	UDP_init(Port,send_ch,receive_ch);
	go checkConnectionStatus(timeStamp_chan,Queue_Network_chan,ReceiveOrderCh)
	
	WaitForMessages(Queue_Network_chan,timeStamp_chan,ack_chan,send_ch,receive_ch,ReceiveOrderCh,BroadcastOrderCh,ack_lock_chan)
	

}

func WaitForMessages(Queue_Network_chan chan map[string]Types.Order_queue,timeStamp_chan chan map[string]time.Time,ack_chan,send_ch,receive_ch,ReceiveOrderCh,BroadcastOrderCh chan Types.Message,ack_lock_chan chan int) {
	
	for{
		select{
			case msg:=<-receive_ch:
				go HandleNewMessage(msg,Queue_Network_chan,timeStamp_chan,ack_chan,send_ch,receive_ch,ReceiveOrderCh)
			case msg:=<-BroadcastOrderCh:
				//Update This Computer (More for printing)
				if numberOfElevators!=0 {
					go BroadcastMessage(msg,ack_chan,send_ch,ack_lock_chan)
				}
			}

			

	}



}

func HandleNewMessage(msg Types.Message,Queue_Network_chan chan map[string]Types.Order_queue, timeStamp_chan chan map[string]time.Time, ack_chan chan Types.Message,send_ch, receive_ch,ReceiveOrderCh chan Types.Message) {
				
				ipAddr:=msg.RemoteAddr
				ack:=Types.Message{MessageType:Types.MT_ack}
				NewElevator:=Types.Message{MessageType:Types.MT_new_elevator}
				//Update timestamp for the ipaddress
				network_TimeStamp:=<-timeStamp_chan
				network_TimeStamp[ipAddr]=time.Now()
				timeStamp_chan<-network_TimeStamp

				temp_Queue_Network:=<-Queue_Network_chan
				Queue_Network_chan<-temp_Queue_Network

				if _,ok:=temp_Queue_Network[ipAddr]; ok==false {
					
					//New Computer connected, put in map
					Queue_Network:=<-Queue_Network_chan
					Queue_Network[ipAddr]=msg.Data
					Queue_Network_chan<-Queue_Network
					
					//
					NewElevator.Data=msg.Data
					NewElevator.Mask=msg.Mask
					ReceiveOrderCh<-NewElevator
					// Notify QueueManager with the Queue
					//Connection:=ConnectionStatus{State:"Connected",IpAddr:ipAddr,Queue:msg.Data}
					fmt.Println(ipAddr,"Connected to the network")
					numberOfElevators=numberOfElevators+1
				
				}else if msg.MessageType ==Types.MT_update {
				
						//Update Queue map
						Queue_Network:=<-Queue_Network_chan
						Queue_Network[ipAddr]=msg.Data
						go printAllOrders(Queue_Network)
						Queue_Network_chan<-Queue_Network
						
						//Send ack
						ReceiveOrderCh<-msg
						ack.AckAddr=ipAddr
						send_ch<-ack
						//Notify QueueManager		
						//Send to QueueManager
						
				}else if msg.MessageType == Types.MT_new_order {
						if msg.RecipientAddr==localAddr {
							ReceiveOrderCh<-msg
						}
						
						//Send ack
						ack.AckAddr=ipAddr
						send_ch<-ack
						
				}else if msg.MessageType == Types.MT_ack {
						//Ack recieved, notify ack manager
						ack_chan<-msg

				}


}
func BroadcastMessage(msg Types.Message,ack_chan,send_ch chan Types.Message,ack_lock_chan chan int){
	
	
	//Wait for other acks to finish
	<-ack_lock_chan
	//Sending the message
	send_ch<-msg

	fmt.Println("Message sended,Waiting for ack")
	
	//Initalize variables
	TTR:=time.Now()
	Timeouts:=0
	finished:=0
	ackcounter:=0
	
	if numberOfElevators>0 {
		//fmt.Println("In for loop",numberOfElevators)
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
							if Timeouts==5 {
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
		ack_lock_chan<-1
}

func checkConnectionStatus(timeStamp_chan chan map[string]time.Time,Queue_Network_chan chan map[string]Types.Order_queue,ReceiveOrderCh chan Types.Message){
	var msg Types.Message
	msg.MessageType=Types.MT_disconnected
	for{
			temp_timeStamp:=<-timeStamp_chan
			timeStamp_chan<-temp_timeStamp
			timeEnd:=time.Now()
			//Sjekker for alle heisene
			
			for ipAddr,timeStart:=range temp_timeStamp {		
				if timeEnd.Sub(timeStart)>=300*time.Millisecond {
					//Computer Disconnected from network or not responding (loop?)
					//Comp_chan in use only if another computer is disconnected
					
					//Delete The computer from map
					Queue_Network:=<-Queue_Network_chan
					network_TimeStamp:=<-timeStamp_chan
					//Elevator died, redestribute orders
					msg.Data=Queue_Network[ipAddr]
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)
					
					timeStamp_chan<-network_TimeStamp	
					Queue_Network_chan<-Queue_Network
					numberOfElevators=numberOfElevators-1
					if numberOfElevators==0 {
						fmt.Println("No more elevators in the system, going single-mode")
					}
					if isMaster(Queue_Network_chan) {
						ReceiveOrderCh<-msg
					}
					
					//Slett fra Computers arrayet
					//Connection:=ConnectionStatus{State:"Disconnected",Queue:QueueData}
					fmt.Println(ipAddr,"Disconnected from the network")
					
				}
			}
			
			
		time.Sleep(100*time.Millisecond)

	}


}
func isMaster(Queue_Network_chan chan map[string]Types.Order_queue)(bool){
	s:= strings.Split(localAddr,":")
	myip,_:=strconv.Atoi(s[1])// string to int
	if numberOfElevators>0 {
		temp_Queue_Network:=<-Queue_Network_chan
		Queue_Network_chan<-temp_Queue_Network
		for ipAddr,_:=range temp_Queue_Network {
			s:=strings.Split(ipAddr,":")
			nextip,err:=strconv.Atoi(s[1])
			if err!=nil {fmt.Println(err)}
			if myip>nextip {
				return false
			}
		}
	}
	return true
}

func printAllOrders(temp_Queue_Network map[string]Types.Order_queue){
	fmt.Print("\033c")
	fmt.Println("------------------------List of Elevator Queues-------------------------------")
	fmt.Println("ipAddr                 | O1 | O2 | O3 | N1 | N2 | N3 | 1  | 2  | 3  | 4  |")
	fmt.Println("--------------------------------------------------------------------------------------")
	for ipAddr,Queue:=range temp_Queue_Network {
		fmt.Print(ipAddr+"    ")
		for i:=0; i<3; i++ {
			fmt.Print(Queue.Outside_order_up[i])
			fmt.Print("    ")
		}
		for i:=1; i<4; i++ {
			fmt.Print(Queue.Outside_order_down[i])
			fmt.Print("    ")			
		}
		for _,k:=range Queue.Inside_order {
			fmt.Print(k)
			fmt.Print("    ")
		}
		fmt.Println("|")
		fmt.Println("-------------------------------------------------------------------------------------")		
	}

}
