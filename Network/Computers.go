package Network

import(
	"fmt"
	"time"
	"../Queue_manager"
	//"net"


)
var network_TimeStamp=map[string]time.Time{} //private
var Queue_Network = map[string]Queue_manager.Queue_type{}

func Network_Manager(Port string,Queue_chan chan int,new_message chan Message){ //,Ip_chan chan string,Comp_chan chan map[string]int){





	fmt.Println(Queue_manager.Queue)
	//Comp_chan:=make(chan map[string]int,1)
	//Ip_chan:=make(chan string,1)

	//Initialize UDP connection
	send_ch :=make(chan Message,1024)
	receive_ch :=make(chan Message,1024)
	UDP_init(Port,send_ch,receive_ch);

	//initialize channels
	time_chan:=make(chan int,1)
	elev_chan:=make(chan int,1)
	time_chan<-1
	elev_chan<-1

	//Initialize variables
	var ipAddr string
	//lastRecieved:=time.Now()
	<-Queue_chan
	msgAlive:=Message{MessageType: "I am alive",Data: Queue_manager.Queue,RemoteAddr: " "}
	msg:=msgAlive
	Queue_chan<-1
	
	go check_ComputerConnection(time_chan,elev_chan)

	for{
		select{
			case msg:=<-receive_ch:
				ipAddr=msg.RemoteAddr
				if msg.MessageType != "I am alive" {

					if msg.MessageType == "Update" {
						<-elev_chan
						Queue_Network[ipAddr]=msg.Data
						elev_chan<-1
						//Notify QueueManager
					}else if msg.MessageType == "New order" {
							//Notify 
					}
					<-time_chan
					network_TimeStamp[ipAddr]=time.Now()
					time_chan<-1
					//Send to QueueManager
				}else if _,ok:=Queue_Network[ipAddr]; ok==false {
					//New Computer connected, put in map
					<-elev_chan
					Queue_Network[ipAddr]=msg.Data
					elev_chan<-1
					// what if elevator have order already????
					fmt.Println(ipAddr,"Connected to the network")
					//fmt.Println(temp_Computers)
					<-time_chan
					network_TimeStamp[ipAddr]=time.Now()
					time_chan<-1
					//fmt.Println("TimeComputers updated")
				}else{
					//Computer already connected, update timestamp
					<-time_chan
					network_TimeStamp[ipAddr]=time.Now()
					time_chan<-1

			}
			case update_msg:=<-new_message:
				msg=update_msg
				msgAlive.Data=update_msg.Data
				msgAlive.MessageType="I am alive"
				send_ch<-msg
			case send_ch<-msg:
				//Do something
			default:
				/*if lastRecieved.Sub(time.Now()) > 2*time.Second {
					fmt.Println("No Computers on network, going singlemode :/")
				}else {
					lastRecieved=time.Now()
				}
				*/
				//Do something
			}
		
		
		
		


	}


}
func messageHandler(msg Message){


}
func check_ComputerConnection(time_chan chan int,elev_chan chan int){
	for{

		id:=<-time_chan
		if id==2 {
			time_chan<-2
		}else{


			timeEnd:=time.Now()
			for ipAddr,timeStart:=range network_TimeStamp {
				//fmt.Println("in loop for check")
				//fmt.Println(temp_timeComputers)
				
				if timeEnd.Sub(timeStart)>350*time.Millisecond {
					//fmt.Println(timeEnd.Sub(timeStart))
				}
				
				if timeEnd.Sub(timeStart)>=500*time.Millisecond {
					//Computer Disconnected from network or not responding (loop?)
					//Comp_chan in use only if another computer is disconnected
					<-elev_chan
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)
					elev_chan<-1
					//Slett fra Computers arrayet
					fmt.Println(ipAddr,"Disconnected from the network")
				}
				//fmt.Println(ipAddr)
				

			}
			
			time_chan<-2
		}
		//fmt.Println("waiting again")
		//time.Sleep(2*time.Second)





	}


}

//Maybe a message handler module????
//func Broadcast_new_order(
