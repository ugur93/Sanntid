package Network

import(
	"fmt"
	"net"
	"time"
	//"runtime"
	"encoding/json"
)


const BroadcastAddr="78.91.73.255"
//const BroadcastAddr="192.168.0.255"


type Message struct {
	MessageType string
	Data byte
	RemoteAddr string
}
var localAddr string

func UDP_init(Port string,send_ch,receive_ch chan Message){
	
	go UDP_send(BroadcastAddr+":"+Port,send_ch)
	go UDP_receive(":"+Port,receive_ch)
}

func UDP_send(addr string,send_ch chan Message){
	con,err:=net.Dial("udp4",addr);
	if err!=nil {
		fmt.Println("Error Dial",err)
	}
	localAddr=con.LocalAddr().String()
	fmt.Println("My ip adress is: ",con.LocalAddr())
	//time.Sleep(10000*time.Second)
	for{
		msg:=<-send_ch
		message,err:=json.Marshal(msg)
		if err!=nil {
			fmt.Println("Error with Marshal: ",err)
		}
		_,err=con.Write([]byte(message))	
		if(err!=nil){
			fmt.Println("Error with Write: ",err)
		}
		time.Sleep(300*time.Millisecond)
	}	
}

func UDP_receive(port string,receive_ch chan Message){
	addr,err:=net.ResolveUDPAddr("udp",port)
	sock,_:=net.ListenUDP("udp",addr)
	if err!=nil {
		fmt.Println(err)
	}
	//timeStart:=time.Now();
	for{	
		msg:=Message{}
		buffer:=make([]byte,1024)
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



