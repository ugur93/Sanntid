package Network

import(
	"fmt"
	"time"
	//"../Queue_manager"
	//"net"
	//"builtin"
	"../Types"


)
var network_TimeStamp=map[string]time.Time{} //private
var Queue_Network = map[string]Types.Queue_type{}
var numberOfElevators int
var ackFinished int
func Network_Manager(Port string,new_message chan Message,stop_chan chan int,Order_update chan Message,Queue_Network_lock_chan chan int){ //,Ip_chan chan string,Comp_chan chan map[string]int){


	
	
	
	//fmt.Println(Queue)
	//Comp_chan:=make(chan map[string]int,1)
	//Ip_chan:=make(chan string,1)
	//Initialize UDP connection
	send_ch :=make(chan Message,1024)
	receive_ch :=make(chan Message,1024)
	UDP_init(Port,send_ch,receive_ch);

	//initialize channels
	time_chan:=make(chan int,1)
	ack_chan:=make(chan Message,1024)
	time_chan<-1
	elev_chan<-1
	ack:=Message{MessageType:"ack"}

	ackFinished=1
	numberOfElevators=0
	//Initialize variables
	var ipAddr string
	//lastRecieved:=time.Now()
	go check_ComputerConnection(time_chan,elev_chan)

	for{
		select{
			case msg:=<-receive_ch:
				//fmt.Println("her")
				ipAddr=msg.RemoteAddr
				<-time_chan
				network_TimeStamp[ipAddr]=time.Now()
				time_chan<-1
				if _,ok:=Queue_Network[ipAddr]; ok==false {
					//New Computer connected, put in map
					<-Queue_Network_lock_chan
					Queue_Network[ipAddr]=msg.Data
					Queue_Network_lock_chan<-1
					// what if elevator have order already????
					fmt.Println(ipAddr,"Connected to the network")
					numberOfElevators=numberOfElevators+1
				}else if msg.MessageType == "Update" {
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						Queue_Network_lock_chan<-1
						fmt.Println("Updated")
						send_ch<-ack
						//Notify QueueManager		
					//Send to QueueManager
				}else if msg.MessageType == "New order" {
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						Queue_Network_lock_chan<-1
						Order_update<-msg
						send_ch<-ack

				}else if msg.MessageType == "ack" {
						ack_chan<-msg

				}
	
			case msg:=<-new_message:
				send_ch<-msg
				go waitforack(msg,ack_chan,send_ch)
			default:
			}
		
		//time.Sleep(100*time.Millisecond)
		
		


	}


}
func waitforack(msg Message,ackchan,send_ch chan Message){
	ackcounter:=0
	fmt.Println("Waiting for ack")
	
	for {
		if ackFinished == 1 {
			ackFinished = 0
			break
		}
	}
	fmt.Println("In for loop",numberOfElevators)
	TTR:=time.Now()
	if numberOfElevators>0 {
		for {
				msg:=<-ackchan
				fmt.Println("waiting")
				ackcounter=ackcounter+1
				if ackcounter==numberOfElevators {
					break
				}else if time.Now().Sub(TTR)>100*time.Millisecond {
					fmt.Println("Ack timed out")
					TTR=time.Now()
					send_ch<-msg
				}
			}
	}
	ackFinished=1
	
	fmt.Println("Recieved ack from all")
}

func check_ComputerConnection(time_chan chan int,Queue_Network_lock_chan chan int){
	for{

		
		

			<-time_chan
			temp_TimeStamp:=network_TimeStamp
			time_chan<-1
			timeEnd:=time.Now()
			for ipAddr,timeStart:=range temp_TimeStamp {
				//fmt.Println("in loop for check")
				//fmt.Println(temp_timeComputers)
				
				if timeEnd.Sub(timeStart)>200*time.Millisecond {
					fmt.Println(ipAddr,timeEnd.Sub(timeStart))
				}
				
				if timeEnd.Sub(timeStart)>=300*time.Millisecond {
					//Computer Disconnected from network or not responding (loop?)
					//Comp_chan in use only if another computer is disconnected
					<-Queue_Network_lock_chan
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)
					//fmt.Println("Lenght is: ",builtin.len(Queue_Network))
					//fmt.Println(Queue_Network);	
					Queue_Network_lock_chan<-1
					//Slett fra Computers arrayet

					fmt.Println(ipAddr,"Disconnected from the network")
					numberOfElevators=numberOfElevators-1
				}
				
				//fmt.Println(ipAddr)
				

			}
			
			
		time.Sleep(50*time.Millisecond)
		//fmt.Println("waiting again")
		//time.Sleep(2*time.Second)





	}


}

//Maybe a message handler module????
//func Broadcast_new_order(
