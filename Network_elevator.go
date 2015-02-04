package main

import(
	"fmt"
	"net"
	"time"
	"strconv"
	//"os"
	"encoding/json"
)

		

const broadcastListenPort=20017

type message struct {
	messageType string
	Data byte
	id int
	


}
func UDP_init(message_size int, send_ch, receive_ch chan message){

	baddr,err:=net.ResolveUDPAddr("udp4","255.255.255.255:"+strconv.Itoa(broadcastListenPort));
	fmt.Println(baddr)
	if err!=nil {
		fmt.Printf("Could not resolve UDP addr: ")
		//os.exit()
	}
	//writeConn,err:=net.DialUDP("udp4",nil,baddr)
	//if err!=nil {
	//	fmt.Printf("Could not make socket");
	//}
	
	//Listen Broadcast
	//readConn,err:=net.ListenUDP("udp4",baddr)
	//if err!=nil {
	//	fmt.Printf("Could not connect")
	//}

	//go UDP_receive(message_size,readConn,receive_ch)
	//go UDP_transmit(message_size,writeConn,send_ch)
	
	
	

	
	
		

}

func UDP_receive(message_size int,receive_ch chan message){
	
	baddr,_:=net.ResolveUDPAddr("udp","129.241.187.255:20017");
	readConn,err:=net.ListenUDP("udp",baddr)
	if err!=nil {
		fmt.Printf("Could not resolve UDP addr: ")
		//os.exit()
	}	
	for{
		buffer:=make([]byte,1024)
		msg:=message{}
		fmt.Println("Reading...")
		n,_,err := readConn.ReadFromUDP(buffer)	
		if err!=nil {
			fmt.Println("error reading")
		
		}	
		fmt.Println(buffer)
		err=json.Unmarshal(buffer[:n],&msg)
		if err!=nil {
			fmt.Printf("Errror")
		}
		fmt.Println(buffer)
		receive_ch<-msg

	}


}

func UDP_transmit(message_size int, send_ch chan message){
	
	//baddr,err := net.ResolveUDPAddr("udp","129.241.187.255:"+strconv.Itoa(broadcastListenPort))
	
	//if err!=nil {
	//	fmt.Printf("Could not resolve UDP addr: ")
	///	//os.exit()
	//}
	writeConn,_:=net.Dial("udp","129.241.187.255:20017")
	for{
		
		msg :=	<-send_ch
		b,_:=json.Marshal(msg)
		writeConn.Write([]byte(b))
		time.Sleep(time.Millisecond*100)
		
	}




}


func main(){
	send_ch:=make(chan message,1024)
	receive_ch:=make(chan message,1024)
	message_size:=1024;
	send_msg:=message{messageType: "I am alive",Data: 2,id: 5};
	//receive_msg:=message{}
	//UDP_init(message_size,send_ch,receive_ch)
	fmt.Println("Test1")
	send_ch<-send_msg;
	go UDP_receive(message_size,receive_ch)
	go UDP_transmit(message_size,receive_ch)
	
	for{
		send_ch<-send_msg;
		receive_msg:=<-receive_ch;
		//fmt.Println(receive_msg)
		
		fmt.Println(receive_msg)
		time.Sleep(time.Millisecond*1000)
	}




}
