package Network

import(
	"fmt"
	"net"
	"time"
	//"runtime"
	"encoding/json"
)


const BroadcastAddr="192.168.0.255"
const SendPort="20020"
const ReadPort="20020"


type Message struct {
	MessageType string
	MessageId int 
	Data byte
	BroadcastPort string
	RemoteAddr *net.UDPAddr
}
func UDP_send(addr string,send_ch chan Message){
	con,err:=net.Dial("udp4",addr);
	if err!=nil {
		fmt.Println("Error Dial",err)
	}
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
		time.Sleep(100*time.Millisecond)
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
		msg.RemoteAddr=Raddr
		
		receive_ch<-msg
		
	}
}

func UDP_init(Port string,send_ch,receive_ch chan Message){
	go UDP_send(BroadcastAddr+":"+Port,send_ch)
	go UDP_receive(":"+Port,receive_ch)
}

/*
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
*/
