package Types

//import "../driver"

const (
	MT_update="Update"
	MT_new_order="NewOrder"
	MT_disconnected="Disconnected"
	MT_ack="ack"

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
		LastFloor int 
		Moving_direction int
		Moving bool
}

type Message struct {
	MessageType string
	Data Order_queue
	Mask Order_queue
	RemoteAddr string //SenderAddr
 	RecipientAddr string
	AckAddr string
}

//[N_FLOORS+2*(N_FLOORS-1)+2]int
