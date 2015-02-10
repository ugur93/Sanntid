package main

import "fmt"
import "./driver"



func main() {
	fmt.Println("test")
	driver.Set_stop_lamp(1);
	driver.Set_motor_direction(1);
}