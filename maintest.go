package main

import "fmt"
import "./driver"



func main() {
	driver.Driver_init()
	fmt.Println("test",driver.Get_stop_signal())
	driver.Set_stop_lamp(0)
	driver.Set_motor_direction(0)
}
