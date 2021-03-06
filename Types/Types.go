package Types



const (
	MT_update="Update"
	MT_new_order="New_order"
	MT_disconnected="Disconnected"
	MT_ack="Ack"
	MT_new_elevator="New_elevator"
	MT_out="Out"

)
const N_FLOORS int = 4

const(
	BUTTON_CALL_UP int =0
	BUTTON_CALL_DOWN int =1
	BUTTON_COMMAND int = 2

)
const (
	DIRN_DOWN int=-1
	DIRN_STOP int=0
	DIRN_UP int=1

)

type Order_queue struct {
		Outside_order_down [(N_FLOORS)]int
		Outside_order_up [(N_FLOORS)]int
		Inside_order [N_FLOORS]int
		Last_floor int 
		Moving_direction int
		Moving bool
}

type Message struct {
	Message_type string
	Data Order_queue
	Mask Order_queue
	Remote_addr string //SenderAddr
 	Recipient_addr string
	Ack_addr string
}


