package main

import(
	"fmt"
	"net"
	"time"
	"os"
)
func MakeConnection(endroutine chan int){
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "129.241.187.145:20017")
	l,_:=net.ListenTCP("tcp",tcpAddr)
	conn,_:=l.AcceptTCP()
	//endroutine<-1
	handleRequest(conn,endroutine)



}

func main(){	
	endroutine :=make(chan int)
	
	//serverAddr:="localhost:33546"
	fmt.Println("hello");

	tcpAddr, err := net.ResolveTCPAddr("tcp", "129.241.187.136:33546")
	if err != nil {
        	fmt.Println("ResolveTCP failed: ",err.Error())
		os.Exit(1)

	}

	conn,err:=net.DialTCP("tcp",nil,tcpAddr);

	if err!=nil {
		fmt.Println("Dial failed: ",err.Error());
		os.Exit(1)
	}

	go MakeConnection(endroutine)
	_,err = conn.Write([]byte("Connect to: 129.241.187.145:20017\x00"))
	 if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	} 
	reply := make([]byte, 1024)
	conn.Read(reply)
	fmt.Println("Reply form server: ",string(reply))	 

	<-endroutine
    




}

func handleRequest(conn net.Conn,endroutine chan int){

	for{

		strEcho :="Fak ye\x00"
		_,err:=conn.Write([]byte(strEcho))
		 if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		 } 
		fmt.Println("Write to server: ",strEcho)
		reply := make([]byte, 1024)
		conn.Read(reply)
		 
		fmt.Println("Reply form server: ",string(reply))
		time.Sleep(time.Second)
	}			
	endroutine<-1


}
