package main

import "fmt"
import "./driver"



func main() {
	driver.Driver_init()
	fmt.Println("test",driver.Get_stop_signal())
	driver.Set_motor_direction(1)
	for{
		if driver.Get_floor_sensor_signal()==3 {
			driver.Set_motor_direction(0)
		}else if driver.Get_floor_sensor_signal()==1 {
			driver.Set_motor_direction(1)
		}
	}
}
