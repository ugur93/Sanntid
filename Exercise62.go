package main

import "fmt"
import "io/ioutil"
//import "encoding/binary"
//import "bytes"	
import ("encoding/json"
	"time"
	"os")


type data struct {
	Counter int
}

func main(){
	//i :=0;
	var run bool = false;
	var CounterOld int =0; 
	msg:=data{}
	
	
	buffer,_:=ioutil.ReadFile("test")
	msgfrom:=data{}
	json.Unmarshal(buffer,&msgfrom)
	CounterOld=msgfrom.Counter;
	 
	for{
		if run==true {
			
			file,err:=os.OpenFile("test",	os.O_APPEND|os.O_WRONLY,0777);
			 if err != nil {
				 panic(err)
			     }
			     
			msg.Counter=msg.Counter+1;
			 sendMsg,_:=json.Marshal(msg)
			 _,_=file.Write(sendMsg)
			/*
			ioutil.WriteFile("test",sendMsg,0777);
			msg.Counter=msg.Counter+1;
			fmt.Println("Counter is: ",msg.Counter);
			*/
			
			file.Close()
			time.Sleep(time.Second)
			buffer:=make([]byte,1000);
			file,err=os.OpenFile("test",os.O_RDONLY,0666);
			 if err != nil {
				 panic(err)
			     }
			n,err:=file.Read(buffer)		
			 if err != nil {
					 panic(err)
				     }
			msgfrom:=data{}
			json.Unmarshal(buffer[n-11:n],&msgfrom)	
			fmt.Println("len: ",string(buffer[n-10]));
		}else {
			time.Sleep(time.Second);
			buffer,_:=ioutil.ReadFile("test")
			msgfrom:=data{}
			json.Unmarshal(buffer,&msgfrom)
			if CounterOld==msgfrom.Counter {
				run=true;
				msg.Counter=msgfrom.Counter;
			
			}else {
				CounterOld=msgfrom.Counter;
			
			}
		}
	}
	
}
