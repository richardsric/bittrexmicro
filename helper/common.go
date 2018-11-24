package helper

import "reflect"

//IsNil Checks for Nil Interface and recovers
func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
