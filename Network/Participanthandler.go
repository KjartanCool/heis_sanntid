package Network

import (
	. "../Driver"
	. "../Variables"
	"fmt"
	"time"
)

// Updates status_struct, removes dead elevators and puts its orders on channel
func Update_participants(stat Status_struct, dead_orders chan Status_struct,light_chan chan [4][3]int) {
	// fmt.Println(Participant_status)
	Check_Participants_Alive(dead_orders)
	Update_participant_info(stat)
	new_global_orders := Update_global_orders()
	Listen_for_ext_lights(new_global_orders, light_chan)
	Global_orders = new_global_orders
}

// HELPFUNCTIONS //
func Update_participant_info(status_struct Status_struct){

	fmt.Println(Participant_status)
	for i := 0; i < len(Participant_status); i++ {

		if Participant_status[i].Ip_tag == status_struct.Ip_tag {
			Participant_status[i] = status_struct
			return
		}
	}
	Participant_status = append(Participant_status, status_struct)
	return
}

func Check_Participants_Alive(dead_orders chan Status_struct) {
	Timestamp_now := time.Now().UnixNano()
	for i := 0; i < len(Participant_status); i++ {
		Time_Difference := Timestamp_now - Participant_status[i].Timestamp
		lokal_ip,_:=Get_NonLoopBack_Ip();
		if Time_Difference > 3000000000 && Participant_status[i].Ip_tag != lokal_ip.String() {
			if i == (len(Participant_status) - 1){
				slice1 := Participant_status[0:i]
				slice2 := Participant_status
				Participant_status = slice1
				dead_orders <- slice2[i]
			} else if i == 0 {
				slice1 := Participant_status[i+1 : len(Participant_status)]
				slice2 := Participant_status
				Participant_status = slice1
				dead_orders <- slice2[i]
				i--

			} else {
				slice1 := Participant_status[0:i]
				slice2 := Participant_status[i+1 : len(Participant_status)]
				slice1 = append(slice1, slice2...)
				slice3 := Participant_status
				Participant_status = slice1
				dead_orders <- slice3[i]
				i--
			}
		} else {
			continue
		}
	}
}

func Update_global_orders() [4][3]int{
	var new_global_orders [4][3]int
	for i:=0; i<len(Participant_status);i++{
		for j:=0; j<N_FLOORS;j++{
			for k :=0; k<M_BUTTONS;k++{
				new_global_orders[j][k] |= Participant_status[i].Work_array[j][k]
			}
		}
	}
	return new_global_orders
}