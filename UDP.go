package main

import(
	"fmt"
	"net"
	"time"
	"runtime"
	"encoding/json"
)


const BroadcastAddr="78.91.75.205"
const SendPort="20020"

const ReadPort="20020"


type message struct {
	MessageType string
	MessageId int //Maybe viktighet??
	Data byte
	BroadcastPort string
	RemoteAddr *net.UDPAddr
}
func UDP_send(addr string,send_ch chan message){
	con,err:=net.Dial("udp4",addr);
	if err!=nil {
		fmt.Println("Error dialing",err)
	}

	for{
		msg:=<-send_ch
		message,err:=json.Marshal(msg)
		if err!=nil {
			fmt.Println("error with Marshal")

		}
		_,err=con.Write([]byte(message))	
		if(err!=nil){
			fmt.Println(err)
		}
		time.Sleep(1000*time.Millisecond)
	}	



}
func UDP_receive(port string,receive_ch chan message){

	addr,err:=net.ResolveUDPAddr("udp",port)
	sock,_:=net.ListenUDP("udp",addr)
	if err!=nil {
		fmt.Println(err)
	}
	for{	
		msg:=message{}
		buffer:=make([]byte,1024)
		n,Raddr,err:=sock.ReadFromUDP(buffer)
		//fmt.Println(string(buffer))
		if err!=nil{
			fmt.Println(err)
		}
		err=json.Unmarshal(buffer[:n],&msg)
		if err!=nil {
			fmt.Println(err);

		}
		msg.RemoteAddr=Raddr
		
		receive_ch<-msg
		
		time.Sleep(1000*time.Millisecond)
	}


}
func UDP_init(port,broadcastAddr int, send_ch,receive_ch chan message){
	
	go UDP_send(BroadcastAddr+":"+SendPort,send_ch)
	go UDP_receive(":"+ReadPort,receive_ch)
	

}
func main(){
	//runtime.GOMAXPROCS(runtime.NumCPU())
	send_ch :=make(chan message,1024)
	receive_ch :=make(chan message,1024)
	UDP_init(send_ch,receive_ch)
	msg:=message{MessageType: "I am alive",MessageId: 1,Data: 2,BroadcastPort: SendPort }
	for{
		send_ch<-msg;
		melding:=<-receive_ch
		
		if melding.MessageType=="I am alive" {
			fmt.Println("Message type: ",melding.MessageType)
			fmt.Println("Message ID: ",melding.MessageId)
			fmt.Println("Data: ",melding.Data)
			fmt.Println("BroadcastPort: ",melding.BroadcastPort)
			fmt.Println("Remote Address: ",melding.RemoteAddr)
			fmt.Println("-------------------------------------------")
		}
	}

	
	


}

