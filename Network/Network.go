package Network

import (
	. "../Driver"
	. "../Variables"
	"fmt"
	. "net"
	"time"
	"errors"
)

// Network module
func Network(got_order chan Order, participant_info chan Status_struct, job chan Order,
			light_chan chan [4][3]int,dead_orders chan Status_struct) {
	for {
		select {
		case a := <-got_order:
			fmt.Println(a, "got_order")
			Cost_function(Participant_status, a, job)
		case b := <-participant_info:
			//fmt.Println(e)
			Update_participants(b, dead_orders, light_chan)
		case <-time.After(10 * time.Millisecond):
			continue
		}
	}
}

// Listens for other elevators
func Listen_status(participant_info chan Status_struct, Participant_status []Status_struct) {
	statusAdr, _ := ResolveUDPAddr("udp", "129.241.187.255:25014")
	statusConn, _ := ListenUDP("udp", statusAdr)
	for {
		select {
		default:
			data := make([]byte, 1024)
			n, _ := statusConn.Read(data)
			Decode_status_info(data[:n])
			participant_info <- Status_decode
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// Broadcasting own status.
func Broadcast_status() {
	statusAdr, _ := ResolveUDPAddr("udp", "129.241.187.255:25014")
	statusConn, _ := DialUDP("udp", nil, statusAdr)
	for {
			status := Make_status()
			message := Encode_status_struct(status)
			Write_status(message,statusConn)
			time.Sleep(10 * time.Millisecond)
	}
}

// Helpfunction for Broadcast_status
func Write_status(message []byte,statusConn *UDPConn) {
	statusConn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))

	_, err := statusConn.Write(message)
	Error_Check(err)
}

// Receives order and puts it on a channel
func Handle_order(got_order chan Order) {
	orderAdr, _ := ResolveUDPAddr("udp", "129.241.187.255:23014")
	orderConn, _ := ListenUDP("udp", orderAdr)
	for {
		Listen_order(got_order, orderConn)
		time.Sleep(2 * time.Millisecond)
	}
}

// Helpfunction for Handler_order
func Listen_order(got_order chan Order, orderConn *UDPConn) {
	data := make([]byte, 1024)
	n, err := orderConn.Read(data)
	Error_Check(err)
	Decode_order_info(data[:n])
	got_order <- Bestilling_decode
}

// Receives a dead elevators orders and puts it got_orders
func Get_dead_elevators_orders(dead_elevator Status_struct, got_order chan Order) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < M_BUTTONS-1; j++ {
			if dead_elevator.Work_array[i][j] == 1 {
				if j == 0 {
					jobb := Order{i, "up"}
					got_order <- jobb
				} else {
					jobb := Order{i, "down"}
					got_order <- jobb
				}
			}
		}
	}
}

// Updates the elevators status
func Make_status() Status_struct {
	timestamp := time.Now().UnixNano()
	Local_Tag, _ := Get_NonLoopBack_Ip()
	tag := Local_Tag.String()
	return Status_struct{Order_matrix, DIRECTION, LAST_FLOOR, tag, timestamp}
}

// Broadcasts and order
func Broadcast_order(order_send Order) {
	order_info := Encode_order(order_send)
	Write_order(order_info)
}

// Helpfunction for Broadcast_order
func Write_order(message []byte) { 
	orderAdr, _ := ResolveUDPAddr("udp", "129.241.187.255:23014")
	orderConn, _ := DialUDP("udp", nil, orderAdr)
	for i:=0; i<30; i++ {
		orderConn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
		_, err := orderConn.Write(message)
		Error_Check(err)
		time.Sleep(1 * time.Millisecond)
	}
	err2 := orderConn.Close()
	Error_Check(err2)
}

// Checks for external lights
func Listen_for_ext_lights(new_global_order [4][3]int, light_chan chan [4][3]int){
    var UPDATED = false
	for i :=0; i <N_FLOORS; i++{ //gaar gjennom participants
		for k:=0; k<M_BUTTONS-1; k++ { //gaar gjennom buttons
			if new_global_order[i][k] != Global_orders[i][k] {
	           UPDATED = true
	        }
	    }
	}
	if UPDATED {
    	light_chan <- new_global_order
	}
}

// Returns the local IP address
func Get_NonLoopBack_Ip()(IP, error){
	net_interface, err := Interfaces()
        if err != nil {
                return nil, err;
        }
        for _, t := range net_interface {
                aa, err := t.Addrs()
                if err != nil {
                        return nil, err
                }
                for _, a := range aa {
                        ipnet, ok := a.(*IPNet)
                        if !ok {
                                continue;
                        }
                        v4 := ipnet.IP.To4()
                        if v4 == nil || v4[0] == 127 { // loopback address
                                continue
                        }
                        return v4, nil
                }
        }
        return nil, errors.New("cannot find local IP address")
}

// Blocking untill a network connection is detected
func Init_network(){
	for{
	     _, err2 :=DialTimeout("udp","129.241.187.255:21014",30*time.Millisecond)
	     if err2 != nil {
	        fmt.Println(err2)
	        continue
	     }else{
	        break
	     }
	     time.Sleep(10*time.Millisecond)           	
	}
}

// Detecs a bad connection
func Init_network2(is_dead chan bool){
    for{
         _, err2 :=DialTimeout("udp","129.241.187.255:21000",30*time.Millisecond)
         if err2 != nil {
            fmt.Println(err2)
            is_dead <- true

         }else{
         }
         time.Sleep(1*time.Millisecond)           	
    }
}