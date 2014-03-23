package driver




/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
#include "channels.h"
*/
import "C"

const N_FLOORS = 4;
const M_BUTTONS = 3;


type Elev_button_type_t int

const (
    BUTTON_CALL_UP Elev_button_type_t = iota
    BUTTON_CALL_DOWN
    BUTTON_COMMAND
)

func Elev_init() int {
    return int(C.elev_init())
}

func Elev_set_speed(speed int) {
    C.elev_set_speed(C.int(speed))
}

func Elev_get_floor_sensor_signal() int {
    return int(C.elev_get_floor_sensor_signal())
}

func Elev_get_button_signal(button Elev_button_type_t, floor int) int {
    return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Elev_get_stop_signal() int {
    return int(C.elev_get_stop_signal())
}

func Elev_get_obstruction_signal() int {
    return int(C.elev_get_obstruction_signal())
}

func Elev_set_floor_indicator(floor int) {
    C.elev_set_floor_indicator(C.int(floor))
}

func Elev_set_button_lamp(button Elev_button_type_t, floor int, value int) {
    C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Elev_set_stop_lamp(value int) {
    C.elev_set_stop_lamp(C.int(value))
}

func Elev_set_door_open_lamp(value int) {
    C.elev_set_door_open_lamp(C.int(value))
}

func Elev_direction()int{
   return int(C.elev_direction());
}
