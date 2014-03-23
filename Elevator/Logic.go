package Elevator

import (
	. "../Driver"
	. "../Variables"
	"fmt"
	"time"
)

// Initalizes the elevator
func Heis_init(door_closed_chan chan bool,quit_chan chan bool) {
	Read_from_file()
	Elev_init()
	Elev_set_speed(300)
	for{
		if Elev_get_floor_sensor_signal() >=0{
			go Keep_door_open(door_closed_chan,quit_chan)
			fmt.Println(Order_matrix)
			Elev_set_speed(0)
			break
		}
	}
}

// Reads from floor sensors and announces if a new floor is reached
func Read_floor_indicator(floor_sensor_chan chan int) {
		last := -1
		for{
			if Elev_get_floor_sensor_signal() != last && Elev_get_floor_sensor_signal() >= 0 {
				last = Elev_get_floor_sensor_signal()
				floor_sensor_chan <- last
			}
			time.Sleep(10*time.Millisecond)
		}	
}

// Same as above, uses another channel
func Read_same_floor(same_floor_chan chan bool) {
		last := -1
		for{
			if Elev_get_floor_sensor_signal() != last && Elev_get_floor_sensor_signal() >= 0 {
				last = Elev_get_floor_sensor_signal()
				same_floor_chan <- true
			}
			time.Sleep(10*time.Millisecond)
		}	
}

// Updates last floor
func Update_last_floor() {
	if Elev_get_floor_sensor_signal() == -1 {
	} else {
		LAST_FLOOR = Elev_get_floor_sensor_signal()
	}
}

// Sets elevator speed
func Set_speed() {
	if DIRECTION == 1 {
		if Check_if_more_orders_in_direction(1) == 1 {
			DIRECTION = 1
		} else {
			DIRECTION = 0
		}
	} else {
		if Check_if_more_orders_in_direction(0) == 1 {
			DIRECTION = 0
		} else {
			DIRECTION = 1
		}
	}
	if DIRECTION == 1 {
		Elev_set_speed(300)
	} else {
		Elev_set_speed(-300)
	}
	MOVING = true
}

// Stops elevator and controls the door
func Door_handler(door_closed_chan chan bool,quit_chan chan bool, same_floor_chan chan bool) {
	MOVING = false
	select{
		case <- same_floor_chan:
			if DIRECTION == 0 {
				Elev_set_speed(200)
				time.Sleep(10 * time.Millisecond)
				Elev_set_speed(-200)
			} else {
				Elev_set_speed(-200)
				time.Sleep(10 * time.Millisecond)
				Elev_set_speed(200)
			}
			break
		default:
			break
	}
	Elev_set_speed(0)
	Elev_set_door_open_lamp(1)
	if DOOR == false{
		go Keep_door_open(door_closed_chan,quit_chan)
	}else{
		quit_chan <- true
		go Keep_door_open(door_closed_chan,quit_chan)
	}
	Remove_order()
}

// Keeps the door open for three seconds
func Keep_door_open(door_closed_chan chan bool, quit_chan chan bool){
	DOOR = true
	timer := time.Now().Unix()
	for {
		select{
			case _ =<- quit_chan:
				return
			default:
				if time.Now().Unix()-timer < 3 {
					time.Sleep(10 * time.Millisecond)
				} else {
					door_closed_chan <- true
					DOOR = false
					return
				}
		}
	}
}

// Reads internal buttons
func Get_internal_signal(internal_order chan Order) {
	var recieved_order Order
	var last [N_FLOORS]int
	recieved_order.Direction = "indre"
	for {
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) != last[i] {
				last[i] = Elev_get_button_signal(BUTTON_COMMAND, i) 
				if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 { 
					recieved_order.Floor = i
					fmt.Println(recieved_order)
					internal_order <- recieved_order
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Reads external buttons
func Get_external_signal(external_order chan Order) {
	var recieved_order Order
	var last [2][N_FLOORS]int
	for {	
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_CALL_UP, i) != last[0][i] {
				last[0][i] = Elev_get_button_signal(BUTTON_CALL_UP, i)
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					recieved_order.Floor = i
					recieved_order.Direction = "up"
					fmt.Println(recieved_order)
					external_order <- recieved_order
					//fmt.Println("kommet igjennom ext")

				}
			}
			if Elev_get_button_signal(BUTTON_CALL_DOWN, i) != last[1][i] {
				last[1][i] = Elev_get_button_signal(BUTTON_CALL_DOWN, i)
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					recieved_order.Floor = i
					recieved_order.Direction = "down"
					external_order <- recieved_order
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// LIGHTS
func Set_ext_lights(light_array [4][3]int) {
	fmt.Println(light_array)
	for j := 0; j < N_FLOORS; j++ { //gaar gjennom etasjer
		for k := 0; k < M_BUTTONS-1; k++ { //gaar gjennom buttons
			if light_array[j][k] == 1 && k == 0 {
				Elev_set_button_lamp(BUTTON_CALL_UP, j, 1)
			} else if light_array[j][k] == 1 && k == 1 {
				Elev_set_button_lamp(BUTTON_CALL_DOWN, j, 1)
			} else if light_array[j][k] == 0 && k == 0 {
				Elev_set_button_lamp(BUTTON_CALL_UP, j, 10)
			} else if light_array[j][k] == 0 && k == 1 {

				Elev_set_button_lamp(BUTTON_CALL_DOWN, j, 10)
			}
		}
	}
}
func Set_int_lights() {

	for i := 0; i < N_FLOORS; i++ {
		if Order_matrix[i][M_BUTTONS-1] == 1 {
			Elev_set_button_lamp(BUTTON_COMMAND, i, 1)

		} else {
			Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
		}

	}
}
func Update_floor_ligth() {
	Elev_set_floor_indicator(LAST_FLOOR)
}