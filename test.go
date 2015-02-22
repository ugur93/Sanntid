package main


import (
	"fmt"
	"time"
)





func main(){
	Computers:=map[string]time.Time{};
	v=Computers
	fmt.Println(v,k)
	timeComputers:=make(chan map[string]time.Time,1)
	timeComputers<-Computers
	temp_timeComputers:=<-timeComputers
	temp_timeComputers["test2"]=time.Now()
	temp_timeComputers["test2"]=time.Now()
	temp_timeComputers["test2"]=time.Now()
	temp_timeComputers["test2"]=time.Now()
	timeComputers<-temp_timeComputers

	test:=<-timeComputers
	for k,v:=range test {
		fmt.Println(k,v)
	}
	

}