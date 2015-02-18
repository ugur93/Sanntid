package main

import "fmt"
import "io/ioutil"
//import "encoding/binary"
//import "bytes"	
import ("encoding/json"
	"time"
	"os/exec")


type data struct {
	Counter int
}

func main(){
	//i :=0;
	var run bool = false;
	var CounterOld int =0; 
	msg:=data{}
	var cmd string = "clear"
	exec.Command("sh","-c",cmd).Output()
	buffer,_:=ioutil.ReadFile("test")
	msgfrom:=data{}
	//times:=time.Now();
	json.Unmarshal(buffer,&msgfrom)
	CounterOld=msgfrom.Counter;
	//time.Sleep(time.Second)
	/*timee:=time.Now();
	if time.Duration(timee.Sub(times))<100*time.Millisecond {
		fmt.Println("Time diff: ",time.Duration(timee.Sub(times)));
	}*/
	//os.Command("clear")
	for{
		if run==true {
	
			sendMsg,_:=json.Marshal(msg)
			ioutil.WriteFile("test",sendMsg,0777);
			fmt.Print("\r")
			fmt.Print("Counter is: ",msg.Counter);
			msg.Counter=msg.Counter+1;
			time.Sleep(time.Millisecond*100)			
		}else {
			time.Sleep(time.Millisecond*200);
			buffer,_:=ioutil.ReadFile("test")
			msgfrom:=data{}
			json.Unmarshal(buffer,&msgfrom)
			if CounterOld==msgfrom.Counter {
				run=true;
				msg.Counter=msgfrom.Counter;
				fmt.Println("\nStarting value: ",msgfrom.Counter);
			
			}else {
				CounterOld=msgfrom.Counter;
			
			}
			time.Sleep(time.Millisecond*200);
			fmt.Print("\r")
			fmt.Print("Waiting .   ")
			time.Sleep(time.Millisecond*200);
			fmt.Print("\r")
			fmt.Print("Waiting ..  ");
			time.Sleep(time.Millisecond*200);
			fmt.Print("\r")
			fmt.Print("Waiting ... ");
			time.Sleep(time.Millisecond*200);
			fmt.Print("\r")
			fmt.Print("Waiting ....");
			
			
		}
	}
	
}
