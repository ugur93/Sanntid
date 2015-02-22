package main

import "fmt"
//import "./driver"
//import "time"
import "./Network"




func main() {
	Network_test();
}

func Network_test(){

	msg:=Network.Message{MessageType: "I am alive",MessageId: 1,Data: 2,BroadcastPort: "2500" }
	send_ch :=make(chan Network.Message,1024)
	receive_ch :=make(chan Network.Message,1024)
	Network.UDP_init("2500",send_ch,receive_ch);
	for {
		select{
		case send_ch<-msg:
		case <-receive_ch:
		default:
			fmt.Println("On default")

		}
		
		/*
		if melding.MessageType=="I am alive" {
			fmt.Println("Message type: ",melding.MessageType)
			fmt.Println("Message ID: ",melding.MessageId)
			fmt.Println("Data: ",melding.Data)
			fmt.Println("BroadcastPort: ",melding.BroadcastPort)
			fmt.Println("Remote Address: ",melding.RemoteAddr)
			fmt.Println("-------------------------------------------")
		}*/

	}
}
/*
func Driver_test(){

	driver.Driver_init()
	fmt.Println("test",driver.Get_stop_signal())
	driver.Set_motor_direction(1)
	for i:=0; i<4; i++ {
		time.Sleep(500*time.Millisecond)
		if i<3 {
			driver.Set_button_lamp(1,i,1)
			driver.Set_button_lamp(0,i,1)
		}
		driver.Set_button_lamp(2,i,1)
	}
	driver.Set_floor_indicator(2)
	pressed:=[]int{0,0,0,0}
	for{
		for i:=0; i<4; i++ {
			if driver.Get_button_signal(2,i)==0 {
				pressed[i]=0;
			}
			if driver.Get_button_signal(2,i)==1 && pressed[i]==0 {
				pressed[i]=1;
				fmt.Print("Button 2, ",i," is pressed!\r")
				if i==3 {
					driver.Set_door_lamp(1)
					driver.Set_stop_lamp(1)
				}
				if i==2 {
					driver.Set_door_lamp(0)
					driver.Set_stop_lamp(0)
				}
			}
		}
		if driver.Get_floor_sensor_signal()==2 {
			driver.Set_motor_direction(-1)
		}else if driver.Get_floor_sensor_signal()==1 {
			driver.Set_motor_direction(1)
		}
	}


}
*/