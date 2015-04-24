package main


import "./Queue_manager"
import "os/signal"
import "os"
import "fmt"
import "syscall"
import "time"
const N_FLOORS int = 4

func main() {
	stop_chan:=make(chan int,1);
	Queue_manager.Queue_manager_init(stop_chan);
    
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    go func() {
        <-c
        Queue_manager.Cleanup()
        stop_chan<-1
        os.Exit(1)
    }()

    for {
        fmt.Println("Working...")
        time.Sleep(10 * time.Second) 
    }
	
	
}

