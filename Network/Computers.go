package Network

import(
	"fmt"
	"time"
	//"net"


)
var Network_elevators=map[string]int{}
var network_TimeStamp=map[string]time.Time{} //private
func Network_Manager(Port string,Queue_chan chan int){ //,Ip_chan chan string,Comp_chan chan map[string]int){
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
	msgAlive:=Message{MessageType: "I am alive",Data: 0,RemoteAddr: " "}

	
	go check_ComputerConnection(time_chan,elev_chan)

	for{
		select{
			case msg:=<-receive_ch:
				ipAddr=msg.RemoteAddr
				if msg.MessageType != "I am alive" {
					<-time_chan
					network_TimeStamp[ipAddr]=time.Now()
					time_chan<-1
					//Send to QueueManager
				}else if _,ok:=Network_elevators[ipAddr]; ok==false {
					//New Computer connected, put in map
					<-elev_chan
					Network_elevators[ipAddr]=1
					elev_chan<-1
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
			case send_ch<-msgAlive:
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
func check_ComputerConnection(time_chan chan int,elev_chan chan int){
	for{

		<-time_chan

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
				Network_elevators[ipAddr]=0
				delete(Network_elevators,ipAddr)
				delete(network_TimeStamp,ipAddr)
				elev_chan<-1
				//Slett fra Computers arrayet
				fmt.Println(ipAddr,"Disconnected from the network")
			}
			//fmt.Println(ipAddr)
			

		}
		
		time_chan<-1

		//fmt.Println("waiting again")
		//time.Sleep(2*time.Second)





	}


}

//Maybe a message handler module????
//func Broadcast_new_order(