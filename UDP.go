package main

import(
	"fmt"
	"net"
	"time"
	"runtime"
)




func UDP_send(addr string,endroutine chan int){

	/*BROADCAST_IPv4 := net.IPv4(255, 255, 255, 255)
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 30000,
	})*/
	num:=0
	for i:=0; i<10000; i++ {
		num++;
		con,_:=net.Dial("udp4",addr);
		buff := []byte("Dude, where's my car?")
		_,err:=con.Write(buff)
		if(err!=nil){
			fmt.Println(err)
		}
		time.Sleep(100*time.Millisecond)
	}	
	endroutine<-1



}
func UDP_receive(port string,endroutine chan int){

	addr,err:=net.ResolveUDPAddr("udp",port)
	sock,_:=net.ListenUDP("udp",addr)
	//sock.SetReadBuffer(1048576)
	if err!=nil 
		fmt.Println(err)
	}
	i:=0
	for{
		i++
		//fmt.Println(i)
		buf:=make([]byte,1024)
		//fmt.Println(i)
		//sock.SetDeadline(time.Millisecond)
		rlen,_,err:=sock.ReadFromUDP(buf)
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println(string(buf[0:rlen]))
		//fmt.Println(i)
		
		time.Sleep(100*time.Millisecond)
	}
	endroutine<-1


}
func UDP_init(){




}
func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	endroutine :=make(chan int)
	
	go UDP_send("255.255.255.255:20017",endroutine)
	go UDP_receive(":20020",endroutine)
	fmt.Println(<-endroutine)
	fmt.Println(<-endroutine)
	
	


}

