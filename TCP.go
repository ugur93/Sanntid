package main

import(
	"fmt"
	"net"
	"bufio"

)


func main(){
	fmt.Println("hello");
	conn, err := net.Dial("tcp", ":33546")
	if err != nil {
        	fmt.Println(err)

	}
	connbuf := bufio.NewReader(conn)
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	for{
	    str, err := connbuf.ReadString('\n')
	    if len(str)>0 {
		fmt.Println(str)
	    }
	    if err!= nil {
		break
	    }
	}
    




}
