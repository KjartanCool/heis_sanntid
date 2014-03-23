package Variables

import(
	"encoding/json"
	"fmt"
)

// Global variables
var Order_matrix [4][3]int
var DIRECTION = 1
var LAST_FLOOR = 0
var MOVING = false
var DOOR = false
var Global_orders [4][3]int

// Structs
type Order struct {
	Floor     int
	Direction string
}
type Status_struct struct {
	Work_array [4][3]int
	DIRECTION  int
	Last_floor int
	Ip_tag     string
	Timestamp  int64
}
type Participant_score struct {
	Tag   string
	Score int
}

// Declarations
var Participant_status []Status_struct
var Bestilling_decode Order
var Status_decode Status_struct
type Participant_scores []Participant_score

// Participants sort functions
func (slice Participant_scores) Len() int {
	return len(slice)
}

func (slice Participant_scores) Less(i, j int) bool {
	return slice[i].Score < slice[j].Score
}

func (slice Participant_scores) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// JSON functions
func Encode_order(order_encode Order) []byte {
	order_info, err := json.Marshal(order_encode)
	Error_Check(err)

	return order_info
}

func Encode_status_struct(order_encode Status_struct) []byte {

	order_info, err := json.Marshal(order_encode)
	Error_Check(err)

	return order_info
}

func Decode_status_info(stat_decode []byte) {
	err := json.Unmarshal(stat_decode, &Status_decode)
	Error_Check(err)
}

func Decode_order_info(order_decode []byte) {
	err := json.Unmarshal(order_decode, &Bestilling_decode)
	Error_Check(err)
}

// Error handling
func Error_Check(err error) {
	if err != nil {
		fmt.Println("failed to solve problem: %s \n", err)
	}
}

