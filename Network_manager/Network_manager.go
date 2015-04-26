package Network_manager

import(
	"fmt"
	"time"
	"../Types"
	"strconv"
	"strings"


)

var number_of_elevators int
//Ip:: 118, 155, 145.146,154,141
func Network_Manager_init(Port string,Broadcast_order_ch,Receive_order_ch chan Types.Message,stop_chan chan int,Queue_Network_chan chan map[string]Types.Order_queue){ 

	//Initialize variables
	number_of_elevators=0
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
	go check_connection_status(timeStamp_chan,Queue_Network_chan,Receive_order_ch)
	Wait_for_messages(Queue_Network_chan,timeStamp_chan,ack_chan,send_ch,receive_ch,Receive_order_ch,Broadcast_order_ch,ack_lock_chan)
	

}

func Wait_for_messages(Queue_Network_chan chan map[string]Types.Order_queue,timeStamp_chan chan map[string]time.Time,ack_chan,send_ch,receive_ch,Receive_order_ch,Broadcast_order_ch chan Types.Message,ack_lock_chan chan int) {
	
	for{
		select{
			case msg:=<-receive_ch:
				go Handle_new_message(msg,Queue_Network_chan,timeStamp_chan,ack_chan,send_ch,receive_ch,Receive_order_ch)
			case msg:=<-Broadcast_order_ch:
				//Update This Computer (More for printing)
				if number_of_elevators!=0 {
					go Broadcast_Message(msg,ack_chan,send_ch,ack_lock_chan)
				}
			}

			

	}



}

func Handle_new_message(msg Types.Message,Queue_Network_chan chan map[string]Types.Order_queue, timeStamp_chan chan map[string]time.Time, ack_chan chan Types.Message,send_ch, receive_ch,Receive_order_ch chan Types.Message) {
				
				ipAddr:=msg.Remote_addr
				ack:=Types.Message{Message_type:Types.MT_ack}
				NewElevator:=Types.Message{Message_type:Types.MT_new_elevator}
				//Update timestamp for the ipaddress
				temp_network_TimeStamp:=<-timeStamp_chan
				timeStamp_chan<-temp_network_TimeStamp				
				
				if _,ok:=temp_network_TimeStamp[ipAddr]; ok==false {
					
					//New Computer connected, put in map
					Queue_Network:=<-Queue_Network_chan
					Queue_Network[ipAddr]=msg.Data
					Queue_Network_chan<-Queue_Network
					
					//
					NewElevator.Data=msg.Data
					NewElevator.Mask=msg.Data
					Receive_order_ch<-NewElevator
					// Notify QueueManager with the Queue
					fmt.Println(ipAddr,"Connected to the network")
					number_of_elevators=number_of_elevators+1
				
				}else if msg.Message_type ==Types.MT_update {
				
		
						Queue_Network:=<-Queue_Network_chan
						Queue_Network[ipAddr]=msg.Data
						go print_all_orders(Queue_Network)
						Queue_Network_chan<-Queue_Network
						
						//Send ack
						Receive_order_ch<-msg
						ack.Ack_addr=ipAddr
						send_ch<-ack
						//Notify QueueManager		
						//Send to QueueManager
						
				}else if msg.Message_type == Types.MT_new_order {
						if msg.Recipient_addr==local_addr {
							Receive_order_ch<-msg
						}
						
						//Send ack
						ack.Ack_addr=ipAddr
						send_ch<-ack
						
				}else if msg.Message_type == Types.MT_ack {
						//Ack recieved, notify ack manager
						ack_chan<-msg

				}else if msg.Message_type ==Types.MT_out {
					
					msg:=Types.Message{Message_type: Types.MT_disconnected}
					Queue_Network:=<-Queue_Network_chan
					msg.Data=Queue_Network[ipAddr]
					msg.Mask=Queue_Network[ipAddr]
					delete(Queue_Network,ipAddr)
					Queue_Network_chan<-Queue_Network
					if is_master(Queue_Network_chan) {
						Receive_order_ch<-msg
					}
					fmt.Println(ipAddr," have some problems, taking over orders")
				}
				network_TimeStamp:=<-timeStamp_chan
				network_TimeStamp[ipAddr]=time.Now()
				timeStamp_chan<-network_TimeStamp
}
func Broadcast_Message(msg Types.Message,ack_chan,send_ch chan Types.Message,ack_lock_chan chan int){
	
	
	//Wait for other acks to finish
	<-ack_lock_chan
	//Sending the message
	send_ch<-msg

	fmt.Println("Message sent, waiting for ack")
	
	//Initalize variables
	TTR:=time.Now()
	Timeouts:=0
	ack_finished:=0
	ackcounter:=0
	
	if number_of_elevators>0 {
		for {
			
				select {
					case ackmsg:=<-ack_chan:
						//Må være sikker på at vi mottar riktig ack
						if ackmsg.Ack_addr==local_addr {
							ackcounter=ackcounter+1
							if ackcounter==number_of_elevators {
								fmt.Println("Recieved ack from all")
								ack_finished=1
								break
						}
					}
					default:
						if time.Now().Sub(TTR)>200*time.Millisecond {
							Timeouts=Timeouts+1
							if Timeouts==5 {
								fmt.Println("Timed out too many times, something is wrong")
								ack_finished=1
								break
							}
							//Resend the message
							fmt.Println("Ack timed out,resending")
							TTR=time.Now()
							send_ch<-msg
						}	
				}
				if ack_finished==1 {
					break
				}
			}
		}
		ack_lock_chan<-1
}

func check_connection_status(timeStamp_chan chan map[string]time.Time,Queue_Network_chan chan map[string]Types.Order_queue,Receive_order_ch chan Types.Message){
	var msg Types.Message
	var Queue_buffer []Types.Message
	var Empty_buffer []Types.Message
	msg.Message_type=Types.MT_disconnected

	Time_lost_elevator:=time.Now()
	for{
			temp_timeStamp:=<-timeStamp_chan
			timeStamp_chan<-temp_timeStamp
			timeEnd:=time.Now()
			//Sjekker for alle heisene
			
			for ipAddr,timeStart:=range temp_timeStamp {		
				if timeEnd.Sub(timeStart)>=300*time.Millisecond {
					
					
					Queue_Network:=<-Queue_Network_chan
					network_TimeStamp:=<-timeStamp_chan
				
					msg.Data=Queue_Network[ipAddr]
					msg.Mask=Queue_Network[ipAddr]
					delete(Queue_Network,ipAddr)
					delete(network_TimeStamp,ipAddr)

					timeStamp_chan<-network_TimeStamp	
					Queue_Network_chan<-Queue_Network
					
					fmt.Println(ipAddr,"Disconnected from the network")
					number_of_elevators=number_of_elevators-1

					if number_of_elevators==0 {
						fmt.Println("No more elevators in the system, going single-mode")
					}

					Queue_buffer=append(Queue_buffer,msg)
					Time_lost_elevator=time.Now()
					
					
				}
		
			}
			if time.Now().Sub(Time_lost_elevator)>=300*time.Millisecond && time.Now().Sub(Time_lost_elevator)<1*time.Second {
				for _,msg:=range Queue_buffer {
					if is_master(Queue_Network_chan) {
						Receive_order_ch<-msg
					}					
			
				}
				Queue_buffer=Empty_buffer
			} 
			
		time.Sleep(100*time.Millisecond)

	}


}
func is_master(Queue_Network_chan chan map[string]Types.Order_queue)(bool){
	s:= strings.Split(local_addr,":")
	s=strings.Split(s[0],"129.241.187.")
	myip,_:=strconv.Atoi(s[1])// string to int
	if number_of_elevators>0 {
		temp_Queue_Network:=<-Queue_Network_chan
		Queue_Network_chan<-temp_Queue_Network
		for ipAddr,_:=range temp_Queue_Network {
			s:=strings.Split(ipAddr,":")
			s=strings.Split(s[0],"129.241.187.")
			nextip,err:=strconv.Atoi(s[1])
			if err!=nil {fmt.Println(err)}
			if myip>nextip {
				return false
			}
		}
	}
	return true
}

func print_all_orders(temp_Queue_Network map[string]Types.Order_queue){
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
