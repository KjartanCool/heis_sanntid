package Elevator

import (
	. "../Network"
	. "../Driver"
	. "../Variables"
	"fmt"
	"time"
	"os"
	"bufio"
)

// Elevator module
func Elevator(got_order chan Order, external_order chan Order, internal_order chan Order, 
				job chan Order, is_dead chan bool, door_closed_chan chan bool, 
				floor_sensor_chan chan int, light_chan chan [4][3]int,quit_chan chan bool,
				same_floor_chan chan bool, dead_orders chan Status_struct) {
	for {
		select {
		case b := <-external_order:
			Broadcast_order(b)
		case c := <-internal_order:
			Write_to_file(Order_matrix)
			if Has_Orders() == false && DOOR == false{ 
				fmt.Println("innejobb")
				Add_Order_Matrix(c)
				if c.Floor == LAST_FLOOR {
					fmt.Println("innejobb i samme etasje")
					Stop_at_order(door_closed_chan,quit_chan,same_floor_chan)
				} else {
					fmt.Println("Setter speed")
					Set_speed()
				}
			} else {
				fmt.Println("else")
				Add_Order_Matrix(c)
				Stop_at_order(door_closed_chan,quit_chan,same_floor_chan)
			}
			Set_int_lights()
		case d := <-job:
			if Has_Orders() == false && DOOR == false {
				fmt.Println("ute i lastfloor")
				Add_Order_Matrix(d)
				if d.Floor == LAST_FLOOR {
					fmt.Println("inne i lastfloor")
					Stop_at_order(door_closed_chan,quit_chan,same_floor_chan)
				} else {
					Set_speed()
				}
			} else {
				Add_Order_Matrix(d)
				Stop_at_order(door_closed_chan,quit_chan,same_floor_chan)
			}
		case g := <-is_dead:
			fmt.Println(g)
			Delete_outside_orders()
			Init_network()
			Emergency_stop()
			
		case  _ =<-floor_sensor_chan:
			Update_last_floor()
			Update_floor_ligth()
			Stop_at_order(door_closed_chan,quit_chan,same_floor_chan)
			Set_int_lights()
			select{
				case a := <- same_floor_chan:
					fmt.Println(a,"Flushed same_floor_chan")
				default:
					// do nothing
			}
		case  <-door_closed_chan:
			Elev_set_door_open_lamp(0)
			if Has_Orders() {
				Set_speed()
			}
		case k := <-light_chan:
			Set_ext_lights(k)
		case f := <-dead_orders:
			fmt.Println(f, "Received a dead elevators orders")
			go Get_dead_elevators_orders(f, got_order)
		case <-time.After(10 * time.Millisecond):
			continue
		}
	}
}

// HELPFUNCTIONS // 

// Checks if the elevator has more orders
func Has_Orders() (answer bool) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < M_BUTTONS; j++ {
			if Order_matrix[i][j] == 1 {
				return true
			}
		}
	}
	return false
}

// Deletes all outside orders
func Delete_outside_orders() {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < M_BUTTONS-1; j++ {
			Order_matrix[i][j] = 0
		}
	}
}

// Stops the elvator when a floot is reached
func Emergency_stop() {
	for {
		if Elev_get_floor_sensor_signal() >= 0 {
			Elev_set_speed(0)
			MOVING = false
			break
		}
	}
}

// Stops the elevator if there is an order in that floor
func Stop_at_order(door_closed_chan chan bool,quit_chan chan bool,same_floor_chan chan bool) {
	order_in_dir := Check_if_more_orders_in_direction(DIRECTION)
	fmt.Println(order_in_dir)
	if Elev_get_floor_sensor_signal() != -1 {
		if DIRECTION == 1 {
			if Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] == 1 {
				Door_handler(door_closed_chan,quit_chan,same_floor_chan)
			} else if Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] == 1 {
				if order_in_dir != 1 {
					Door_handler(door_closed_chan,quit_chan,same_floor_chan)
				}
			} else if Order_matrix[LAST_FLOOR][BUTTON_COMMAND] == 1 {
				Door_handler(door_closed_chan,quit_chan,same_floor_chan)
			}
		} else if DIRECTION == 0 {
			if Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] == 1 {
				Door_handler(door_closed_chan,quit_chan,same_floor_chan)
			} else if Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] == 1 {
				if order_in_dir != 1 {
					Door_handler(door_closed_chan,quit_chan,same_floor_chan)
				} else if Order_matrix[LAST_FLOOR][BUTTON_COMMAND] == 1 {
					Door_handler(door_closed_chan,quit_chan,same_floor_chan)
				}
			} else if Order_matrix[LAST_FLOOR][BUTTON_COMMAND] == 1 {
				Door_handler(door_closed_chan,quit_chan,same_floor_chan)
			}
		}
	}
}

// Checks if there is more orders in a given direction
func Check_if_more_orders_in_direction(direction int) (answer int) {
	if direction == 1 {
		if LAST_FLOOR == N_FLOORS {
			return 0
		}
		for i := LAST_FLOOR + 1; i < N_FLOORS; i++ {
			for j := 0; j < M_BUTTONS; j++ {
				if Order_matrix[i][j] == 1 {
					return 1
				} else {
					continue
				}
			}
		}
	} else if direction == 0 {
		if LAST_FLOOR == 0 {
			return 0
		}
		for i := LAST_FLOOR - 1; i >= 0; i-- {
			for j := 0; j < M_BUTTONS; j++ {
				if Order_matrix[i][j] == 1 {
					return 1
				} else {
					continue
				}
			}
		}
	}
	return 0
}

// Adds an order to the order matrix
func Add_Order_Matrix(Next_Order Order) {
	if Next_Order.Direction == "up" {
		Order_matrix[Next_Order.Floor][0] = 1
	} else if Next_Order.Direction == "down" {
		Order_matrix[Next_Order.Floor][1] = 1
	} else {
		Order_matrix[Next_Order.Floor][2] = 1
		fmt.Println(Order_matrix)
	}
}

// Remove an order from the order matrix and writes to file
func Remove_order() {
	order_in_dir := Check_if_more_orders_in_direction(DIRECTION)
	if DIRECTION == 1 {
		if Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] == 1 {
			if order_in_dir != 1 {
				Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] = 0
			} else {
				Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] = 0
			}
		} else if Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] == 1 {
			Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
			Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] = 0
		} else if Order_matrix[LAST_FLOOR][BUTTON_COMMAND] == 1 {
			Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
		}
	} else if DIRECTION == 0 {
		if Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] == 1 {
			if order_in_dir != 1 {
				Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_UP] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] = 0
			} else {
				Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
				Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] = 0
			}
		} else if Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] == 1 {
			Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
			Order_matrix[LAST_FLOOR][BUTTON_CALL_DOWN] = 0
		} else if Order_matrix[LAST_FLOOR][BUTTON_COMMAND] == 1 {
			Order_matrix[LAST_FLOOR][BUTTON_COMMAND] = 0
		}
	}
	Write_to_file(Order_matrix)
	//fmt.Println(Order_matrix);
}

// Creates/overwrites file and writes to it
func Write_to_file(order_matrix [4][3]int){
	file, err := os.Create("output.txt")
    if err != nil { panic(err) }
    // close file on exit and check for its returned error
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()
	buf := make([]byte,4)
	for i := 0; i < N_FLOORS; i++{
		buf[i] = byte(order_matrix[i][2])
	}
    if _, err := file.Write(buf[:N_FLOORS]); err != nil {
        panic(err)
    }
}

// Reads from file and adds to order matrix
func Read_from_file(){
	file, err := os.Open("output.txt")
    if err != nil { panic(err) }
    // close file on exit and check for its returned error
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()
    r:= bufio.NewReader(file)
    buf :=make([]byte,4)
	for i := 0; i < N_FLOORS; i++{
		b, _ :=r.ReadByte()
		buf[i] = b
		Order_matrix[i][M_BUTTONS-1] = int(b)
	}
	fmt.Println(buf)
}