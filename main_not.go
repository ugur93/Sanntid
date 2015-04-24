package main
	
import "fmt"
import "strconv"
import "./Network"
import "./Types"
import "strings"
import "time"
//	import "github.com/ugur93/Sanntid"



func main() {
	
	stop_chan:=make(chan int)
	Broadcast_order:=make(chan Network.Message,1)
	//Queue_lock_chan:=make(chan int,1)
	Queue_Network_chan:=make(chan map[string]Types.Order_queue,1) //Bruk det når det leses på Queue

	
	 //Bruk dette når du leser på Queue_Network array (andre heisenes køer)
	Received_order_ch:=make(chan Network.Message,1)
	go Network.Network_Manager_init("20020",Broadcast_order,Received_order_ch,stop_chan,Queue_Network_chan)
	msg:=Network.Message{MessageType:"Update"}
	for {
		Broadcast_order<-msg
		time.Sleep(100*time.Millisecond)
	}
	<-stop_chan
	var Queue_Network = map[string][12]int{}
	var x [12]int
	fmt.Print("\033c")
	Queue_Network["192.168.0.155:50553"]=x
	Queue_Network["192.168.0.101:50500"]=x
	fmt.Println("------------------------List of Elevator Queues-------------------------------")
	fmt.Println("ipAddr         | O1 | O2 | O3 | N1 | N2 | N3 | 1  | 2  | 3  | 4  | Dir | F |")
	fmt.Println("______________________________________________________________________________")
	for ipAddr,_:=range Queue_Network {
			    s := strings.Split(ipAddr,":")
			    i,_:=strconv.Atoi(s[1])// string to int
			    fmt.Println(i*2,s[1])
			   /* i, err := strconv.Atoi(s)
			    if err != nil {
			        // handle error
			        fmt.Println(err)
			        //os.Exit(2)
			    }*/
			    //fmt.Println(ipAddr, i)		
		/*
		fmt.Print(ipAddr+"    ")
		for i:=range Queue {
			
			fmt.Print(i)
			fmt.Print("    ")
		}
		fmt.Println("|")
		fmt.Println("---------------------------------------------------------------------")	
		*/
	}
}
