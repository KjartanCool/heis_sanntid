package Network

import (
	. "../Driver"
	. "../Variables"
	"fmt"
	. "net"
	"sort"
	. "math"
)

// Deciding which elevator gets an order
func Cost_function(participant_status []Status_struct, order Order, job chan Order) {

	//fmt.Println(participant_status)

	var score_status_array Participant_scores
	var score_status Participant_score
	var score = 0

	for i := 0; i < len(participant_status); i++ {
		if order.Direction == "up" {
			if participant_status[i].Work_array[order.Floor][0] == 1 { // sjekker om jobben allerede er i heisens work_array
				return // hvis den er det, driter man i hele orderen
			}
		} else if order.Direction == "down" {
			if participant_status[i].Work_array[order.Floor][1] == 1 { // sjekker om jobben allerede er i heisens work_array
				return // hvis den er det, driter man i hele orderen
			}
		}
	}

	for i := 0; i < len(participant_status); i++ {
		score = 0

		if participant_status[i].Last_floor <= order.Floor && participant_status[i].DIRECTION == 0 { //jobben er over heisen && DIR er ned
			score = order.Floor - participant_status[i].Last_floor // oppdaterer score utifra avstand

			for j := participant_status[i].Last_floor; j >= 0; j-- { //iterer fra last floor og kjoerer nedover
				for k := 0; k < M_BUTTONS; k++ {
					if participant_status[i].Work_array[j][k] == 1 {
						score = score + (order.Floor - j)
					}
				}

			}
			score_status.Tag = participant_status[i].Ip_tag
			score_status.Score = score
			score_status_array = append(score_status_array, score_status)
			fmt.Println(score)
			fmt.Println(score_status_array)

		} else if participant_status[i].Last_floor >= order.Floor && participant_status[i].DIRECTION == 1 { //jobben er under heisen &&  DIR er opp

			score = participant_status[i].Last_floor - order.Floor

			for j := participant_status[i].Last_floor; j < N_FLOORS; j++ {
				for k := 0; k < M_BUTTONS; k++ {

					if participant_status[i].Work_array[j][k] == 1 {

						score = score + (j - order.Floor)
					} else {
					}

				}

			}

			score_status.Tag = participant_status[i].Ip_tag
			score_status.Score = score
			score_status_array = append(score_status_array, score_status)
			fmt.Println(score)
			fmt.Println(score_status_array)

		} else {

			if participant_status[i].Last_floor < order.Floor { //jobben er over og DIR er opp
				score = order.Floor - participant_status[i].Last_floor
				for j := participant_status[i].Last_floor; j <= 0; j++ {
					for k := 0; k < M_BUTTONS; k++ {
						if participant_status[i].Work_array[j][k] == 1 { //iterer oppover og finner jobber imellom last_floor og order
							score = score + int(Abs(float64(order.Floor - j)))

						}

					}
				}

				score_status.Tag = participant_status[i].Ip_tag
				score_status.Score = score
				score_status_array = append(score_status_array, score_status)
				fmt.Println(score)
				fmt.Println(score_status_array)

			} else if participant_status[i].Last_floor > order.Floor { //jobben er under og DIR er ned
				score = participant_status[i].Last_floor - order.Floor
				for j := participant_status[i].Last_floor; j >= 0; j-- { //iterer nedover og finner jobber imellom last_floor og order
					for k := 0; k < M_BUTTONS; k++ {
						if participant_status[i].Work_array[j][k] == 1 {

							score = score + int(Abs(float64(j - order.Floor)))
						}

					}

				}

				score_status.Tag = participant_status[i].Ip_tag
				score_status.Score = score
				score_status_array = append(score_status_array, score_status)
				fmt.Println(score)
				fmt.Println(score_status_array)

			}
		}
	}

	sort.Sort(score_status_array)
	low := "HORE"
	var slice2 Participant_scores
	for i := 0; i < len(score_status_array); i++ {
		if score_status_array[i].Score == score_status_array[0].Score {
			slice2 = score_status_array[0 : i+1]
			fmt.Println(slice2, "DETTE ER SLICE 2 I FORLØKKEN")
			low = slice2[0].Tag
		} else {
			break
		}
	}
	for j := 0; j < len(slice2); j++ {
		if slice2[j].Tag < low {
			low = slice2[j].Tag
		}
	}
	fmt.Println(low, "DETTE ER LOW UTENFOR FORLØKKEN, RETT FØR VI BRUKERN")
	fmt.Println(slice2, "DETTE ER SLICE 2")
	Local_Tag, _ := Get_NonLoopBack_Ip()
	if low == Local_Tag.String() {
		fmt.Println("EOERHGOERHGEROGHERØOGH")
		job <- order

	} else {
		 /*for i:=0;i<len(participant_status);i++{
		 	if participant_status[i].Ip_tag==low{
		 		for j:=0;j<M_BUTTONS-1; j++{
		 			if order.Floor==j{
		 				if order.Direction=="up"{
		 					Participant_status[i].Work_array[j][0]==1
		 					
		 				}else{
		 					Participant_status[i].Work_array[j][1]==1
		 				} 
		 			}
		 		
		 		}
		 	}
		 }*/
	}
}
