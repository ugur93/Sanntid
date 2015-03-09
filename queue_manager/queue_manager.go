package Queue_manager

//import "../driver" //HOW DOES ONE DO THIS
const N_FLOORS int =4;
type Queue_type [N_FLOORS+2*(N_FLOORS-1)+2]bool
var Queue Queue_type;

func Queue_manager_init(){
	//Queue_chan:=make(chan int);
	Queue[1]=true;
	Queue[2]=true;
	Queue[3]=true;
}
/*
func Get_orders(){ //STILL IN PRODUCTION, how is the queue going to be
	outside_order_ch := make(chan [2]int)
	inside_order_ch := make(chan int)
	var outside_order_array [2]int
	inside_order := int
	go driver.Check_for_outside_order(outside_order_ch)
	go driver.Check_for_inside_order(inside_order_ch)
	for{
		select{
		case outside_order_array = <- outside_order_ch:
			if "order already in queue"{ //YES NOT DONE
				continue
			}
			else{
				Send_order_to_network()
				"Add order to queue"
			}
		case inside_order = <- inside_order_ch:
			if "order already in queue"{ //YES NOT DONE
				continue
			}
			else{
				"Add order to queue"
			}
		}
	}
}

func Send_order_to_network(){

}

func Receive_order_from_network(){

}*/
