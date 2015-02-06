package main
	
import "fmt"
import "Projects/Sanntid/UDP";



func main() {
	send_ch:=make(chan message,1024)
	receive_ch:=make(chan message,1024)
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