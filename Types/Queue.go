package Types

const N_FLOORS int = 4;
type Order_queue struct {
		Outside_order_down [(N_FLOORS-1)]int
		Outside_order_up [(N_FLOORS-1)]int
		Inside_order [N_FLOORS]int
		LastFloor int
		Direction int
}

//[N_FLOORS+2*(N_FLOORS-1)+2]int
