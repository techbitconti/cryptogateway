package logic

import (
	"fmt"
	"reflect"
)

type Writer struct {
	Status int8        `json:"result"`
	Api    string      `json:"api"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func reflectString(value interface{}) bool {

	r := reflect.ValueOf(value)

	fmt.Println("reflect : ", r.Kind(), r)

	if r.Kind() == reflect.String {
		return true
	}

	return false
}

func reflectSlice(value interface{}) bool {

	r := reflect.ValueOf(value)

	fmt.Println("reflect : ", r.Kind(), r)

	if r.Kind() == reflect.Slice {
		return true
	}

	return false
}
