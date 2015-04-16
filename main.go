package main
	
import "fmt"
import "strconv"

import "strings"
//	import "github.com/ugur93/Sanntid"



func main() {
	

	var Queue_Network = map[string][12]int{}
	var x [12]int
	fmt.Print("\033c")
	Queue_Network["192.168.0.155:50553"]=x
	Queue_Network["192.168.0.101:50500"]=x
	fmt.Println("------------------------List of Elevator Queues-------------------------------")
	fmt.Println("ipAddr         | O1 | O2 | O3 | N1 | N2 | N3 | 1  | 2  | 3  | 4  | Dir | F |")
	fmt.Println("______________________________________________________________________________")
	for ipAddr,_:=range Queue_Network {
			    s := strings.Split(ipAddr,":")
			    i,_:=strconv.Atoi(s[1])// string to int
			    fmt.Println(i*2,s[1])
			   /* i, err := strconv.Atoi(s)
			    if err != nil {
			        // handle error
			        fmt.Println(err)
			        //os.Exit(2)
			    }*/
			    //fmt.Println(ipAddr, i)		
		/*
		fmt.Print(ipAddr+"    ")
		for i:=range Queue {
			
			fmt.Print(i)
			fmt.Print("    ")
		}
		fmt.Println("|")
		fmt.Println("---------------------------------------------------------------------")	
		*/
	}
}
