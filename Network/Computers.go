package Network

import(
	"fmt"
	"time"
	//"net"


)
var Computers=map[string]int{}
var timeComputers=map[string]time.Time{}
func handle_ComputerNetwork(Ip_chan chan string,Comp_chan chan map[string]int){
	Comp_chan<-Computers
	time_chan:=make(chan map[string]time.Time,1)
	time_chan<-timeComputers
	time_wait_for_update:=make(chan int)
	temp_Computers:=Computers
	var ipAddr string;
	go check_ComputerConnection(time_chan,Comp_chan,time_wait_for_update)
	for{
		select{
		case ipAddr=<-Ip_chan:
			if _,ok:=temp_Computers[ipAddr]; ok==false {
			//New Computer connected, put in map
			temp_Computers=<-Comp_chan
			temp_Computers[ipAddr]=1
			Comp_chan<-temp_Computers
			fmt.Println(ipAddr,"Connected to the network")
			//fmt.Println(temp_Computers)

			temp_timeComputers:=<-time_chan
			temp_timeComputers[ipAddr]=time.Now()
			time_chan<-temp_timeComputers
			//fmt.Println("TimeComputers updated")
		}else{
			//Computer already connected, update timestamp
			//fmt.Println("Time update starting")
			
			temp_timeComputers:=<-time_chan
			temp_timeComputers[ipAddr]=time.Now()
			time_chan<-temp_timeComputers

			//fmt.Println("updated")
			
			fmt.Println("Time updated")
			//fmt.Println("Time Updated")
			//fmt.Println(temp_timeComputers)

		}
		case <-time_wait_for_update:
		
		default:
		}

		//fmt.Println("Wroking")
		//ipAddr:=<-Ip_chan
		
		
		
		


	}


}
func check_ComputerConnection(time_chan chan map[string]time.Time,Comp_chan chan map[string]int,time_wait_for_update chan int){
	for{
		//fmt.Println("waiting")
		

		temp_timeComputers:=<-time_chan
		//fmt.Println(temp_timeComputers)
		
		//fmt.Println("before loop")
		timeEnd:=time.Now()
		for ipAddr,timeStart:=range temp_timeComputers {
			//fmt.Println("in loop for check")
			//fmt.Println(temp_timeComputers)
			
			if timeEnd.Sub(timeStart)>200*time.Millisecond {
				fmt.Println(timeEnd.Sub(timeStart))
			}
			
			if timeEnd.Sub(timeStart)>=300*time.Millisecond {
				//Computer Disconnected from network or not responding (loop?)
				//Comp_chan in use only if another computer is disconnected
				temp_Computers:=<-Comp_chan
				temp_Computers[ipAddr]=0
				delete(temp_Computers,ipAddr)
				delete(temp_timeComputers,ipAddr)
				Comp_chan<-temp_Computers
				//Slett fra Computers arrayet
				fmt.Println(ipAddr,"Disconnected from the network")
			}
			//fmt.Println(ipAddr)
			

		}
		
		time_chan<-temp_timeComputers
		time_wait_for_update<-1

		//fmt.Println("waiting again")
		//time.Sleep(2*time.Second)





	}


}
