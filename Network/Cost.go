package Network

import (
	. "../Driver"
	. "../Variables"
	. "net"
	"sort"
	. "math"
)

// Deciding which elevator gets an order
func Cost_function(participant_status []Status_struct, order Order, job chan Order) {
	
	var score_status_array Participant_scores
	var score_status Participant_score
	var score = 0

	for i := 0; i < len(participant_status); i++ {
		if order.Direction == "up" {
			if participant_status[i].Work_array[order.Floor][0] == 1 {
				return 
			}
		} else if order.Direction == "down" {
			if participant_status[i].Work_array[order.Floor][1] == 1 { 
				return 
			}
		}
	}
	for i := 0; i < len(participant_status); i++ {
		score = 0
		if participant_status[i].Last_floor <= order.Floor && participant_status[i].DIRECTION == 0 { 
			score = order.Floor - participant_status[i].Last_floor 
			for j := participant_status[i].Last_floor; j >= 0; j-- { 
				for k := 0; k < M_BUTTONS; k++ {
					if participant_status[i].Work_array[j][k] == 1 {
						score = score + (order.Floor - j)
					}
				}
			}
			score_status.Tag = participant_status[i].Ip_tag
			score_status.Score = score
			score_status_array = append(score_status_array, score_status)
		} else if participant_status[i].Last_floor >= order.Floor && participant_status[i].DIRECTION == 1 { 
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
		} else {
			if participant_status[i].Last_floor < order.Floor { 
				score = order.Floor - participant_status[i].Last_floor
				for j := participant_status[i].Last_floor; j <= 0; j++ {
					for k := 0; k < M_BUTTONS; k++ {
						if participant_status[i].Work_array[j][k] == 1 { 
							score = score + int(Abs(float64(order.Floor - j)))
						}
					}
				}
				score_status.Tag = participant_status[i].Ip_tag
				score_status.Score = score
				score_status_array = append(score_status_array, score_status)

			} else if participant_status[i].Last_floor > order.Floor {
				score = participant_status[i].Last_floor - order.Floor
				for j := participant_status[i].Last_floor; j >= 0; j-- { 
					for k := 0; k < M_BUTTONS; k++ {
						if participant_status[i].Work_array[j][k] == 1 {

							score = score + int(Abs(float64(j - order.Floor)))
						}

					}

				}
				score_status.Tag = participant_status[i].Ip_tag
				score_status.Score = score
				score_status_array = append(score_status_array, score_status)
			}
		}
	}
	sort.Sort(score_status_array)
	low := ""
	var slice2 Participant_scores
	for i := 0; i < len(score_status_array); i++ {
		if score_status_array[i].Score == score_status_array[0].Score {
			slice2 = score_status_array[0 : i+1]
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
	Local_Tag, _ := Get_NonLoopBack_Ip()
	if low == Local_Tag.String() {
		job <- order
	} 
}
