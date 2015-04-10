package Network

import(
	"fmt"
	"time"
	"../Types"


)
var network_TimeStamp=map[string]time.Time{} //private
var Queue_Network = map[string]Types.Queue_type{}
var numberOfElevators int
var int ackFinished

type ConnectionStatus struct {
	State string
	IpAddr string
	Queue Types.Queue_type
}
func Network_Manager(Port string,new_message chan Message,stop_chan chan int,New_Order chan Message,Queue_Network_lock_chan chan int){ //,Ip_chan chan string,Comp_chan chan map[string]int){


	//initialize channels
	send_ch :=make(chan Message,1024)
	receive_ch :=make(chan Message,1024)
	time_chan:=make(chan int,1)
	ack_chan:=make(chan Message,1024)
	time_chan<-1

	//Initialize variables
	ack:=Message{MessageType:"ack"}
	ackFinished=1
	numberOfElevators=0
	var ipAddr string
	
	//Calling functions
	UDP_init(Port,send_ch,receive_ch);
	go check_ComputerConnection(time_chan,elev_chan)

	for{
		select{
			case msg:=<-receive_ch:
				ipAddr=msg.RemoteAddr
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
					Connection:={State:"Connected",IpAddr:ipAddr,Queue:msg.Data}
					fmt.Println(ipAddr,"Connected to the network")
					numberOfElevators=numberOfElevators+1
				}else if msg.MessageType == "Update" {
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						go printAllOrders(Queue_Network)
						Queue_Network_lock_chan<-1
						ack.AckAddr=ipAddr
						send_ch<-ack
						//Notify QueueManager		
					//Send to QueueManager
				}else if msg.MessageType == "New order" {
						<-Queue_Network_lock_chan
						Queue_Network[ipAddr]=msg.Data
						go printAllOrders(Queue_Network)
						Queue_Network_lock_chan<-1
						New_Order<-msg
						ack.AckAddr=ipAddr
						send_ch<-ack
				}else if msg.MessageType == "ack" {
						ack_chan<-msg

				}
	
			case msg:=<-new_message:
				send_ch<-msg
				go waitforack(msg,ack_chan,send_ch)
			default:
			}

	}


}
func waitforack(msg Message,ack_chan,send_ch chan Message){
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
				ackmsg:=<-ack_chan
				//Må være sikker på at vi mottar riktig ack
				if ackmsg.AckAddr==LocalAddr {
					ackcounter=ackcounter+1
					if ackcounter==numberOfElevators {
						break
					}else if time.Now().Sub(TTR)>200*time.Millisecond {
						fmt.Println("Ack timed out")
						TTR=time.Now()
						send_ch<-msg
					}
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
			//Sjekker for alle heisene
			for ipAddr,timeStart:=range temp_TimeStamp {		
				if timeEnd.Sub(timeStart)>=300*time.Millisecond {
					//Computer Disconnected from network or not responding (loop?)
					//Comp_chan in use only if another computer is disconnected
					<-Queue_Network_lock_chan
					QueueData:=Queue_Network[ipAddr]
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)	
					Queue_Network_lock_chan<-1
					//Slett fra Computers arrayet
					Connection:={State:"Disconnected",Queue:QueueData}
					fmt.Println(ipAddr,"Disconnected from the network")
					numberOfElevators=numberOfElevators-1
				}
			}
			
			
		time.Sleep(100*time.Millisecond)

	}


}

func printAllOrders(temp_Queue_Network map[string]Types.Queue_type){
	fmt.Println("------------------------List of Elevator Queues------------------------------------")
	fmt.Println("ipAddr        | O1 | O2 | O3 | N1 | N2 | N3 | 1  | 2  | 3  | 4  | Dir | F |")
	fmt.Println("---------------------------------------------------------------------")
	for ipAddr,Queue:=range temp_Queue_Network {
		fmt.Print(ipAddr+"    ")
		for i:=range Queue {
			fmt.Print(i)
			fmt.Print("    ")
		}
		fmt.Println("|")
		fmt.Println("---------------------------------------------------------------------")		
	

}
 

//Maybe a message handler module????
//func Broadcast_new_order(
