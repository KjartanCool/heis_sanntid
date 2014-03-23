package main

import (
	. "./Network"
	. "./Elevator"
	. "./Variables"
	. "./Driver"
	"fmt"
	"time"
)

func main() {
	

	fmt.Println("Elevator starting up...")
	
	internal_order := make(chan Order)
	external_order := make(chan Order)
	got_order := make(chan Order)
	dead_orders := make(chan Status_struct)
	job := make(chan Order)
	participant_info := make(chan Status_struct)
	is_dead := make(chan bool)
	light_chan := make(chan [N_FLOORS][M_BUTTONS]int)
	door_closed_chan := make(chan bool)
	floor_sensor_chan := make(chan int)
	quit_chan := make(chan bool)
	same_floor_chan := make(chan bool)

	Heis_init(door_closed_chan,quit_chan)
	
	go Elevator(got_order,external_order,internal_order,job,is_dead,door_closed_chan,floor_sensor_chan,light_chan,quit_chan, same_floor_chan,dead_orders)
	go Network(got_order,participant_info,job,light_chan, dead_orders)
	go Get_internal_signal(internal_order)
	go Get_external_signal(external_order)
	go Read_floor_indicator(floor_sensor_chan,)
	go Read_same_floor(same_floor_chan)
	go Handle_order(got_order)
	go Broadcast_status()
	go Listen_status(participant_info,Participant_status)
	go Init_network2(is_dead)
	
	fmt.Println("Elevator initialized")
	
	for {
		time.Sleep(100000000 * time.Second)
	}
}	
