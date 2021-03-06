package Network

import(
	"fmt"
	"net"
	"time"
	//"runtime"
	"encoding/json"
	"../Types"
)

//const BroadcastAddr="localhost"//"192.168.0.255"
const BroadcastAddr="129.241.187.255"


type Message struct {
	MessageType string
	Data Types.Order_queue
	Mask Types.Order_queue
	RemoteAddr string //SenderAddr
 	RecipientAddr string
	AckAddr string
}
var localAddr string

func UDP_init(Port string,send_ch,receive_ch chan Message){
	
	go UDP_send(BroadcastAddr+":"+Port,send_ch)
	go UDP_receive(":"+Port,receive_ch)
}

func UDP_send(addr string,send_ch chan Message){
	con,err:=net.Dial("udp",addr);
	if err!=nil {
		fmt.Println("Error Dial",err)
	}
	localAddr=con.LocalAddr().String()
	fmt.Println("My ip adress is: ",con.LocalAddr())
	msgAlive:=Message{MessageType: "I am alive"}
	//time.Sleep(10000*time.Second)
	for{
		select{
			case msg:=<-send_ch:
				msg.RemoteAddr=localAddr
				message,err:=json.Marshal(msg)
				if err!=nil {
					fmt.Println("Error with Marshal: ",err)
				}
				_,err=con.Write([]byte(message))	
				if(err!=nil){
					fmt.Println("Error with Write: ",err)
				}
			default:
				message,err:=json.Marshal(msgAlive)
				if err!=nil {
					fmt.Println("Error with Marshal: ",err)
				}
				_,err=con.Write([]byte(message))	
				if(err!=nil){
					fmt.Println("Error with Write: ",err)
				}
				
			}
		time.Sleep(100*time.Millisecond)
	}	
}

func UDP_receive(port string,receive_ch chan Message){
	addr,err:=net.ResolveUDPAddr("udp",port)
	sock,err:=net.ListenUDP("udp",addr)
	if err!=nil {
		panic(err)
		fmt.Println(err)
	}
	//timeStart:=time.Now();
	msg:=Message{}
	buffer:=make([]byte,1024)
	for{	
		
		
		n,Raddr,err:=sock.ReadFromUDP(buffer)
		if err!=nil{
			fmt.Println(err)
		}
		err=json.Unmarshal(buffer[:n],&msg)
		if err!=nil {
			fmt.Println(err);
		}
		//Check if not reading own message
		if Raddr.String()!=localAddr{
			msg.RemoteAddr=Raddr.String();
			receive_ch<-msg
		}
		
		
		
	}
}



