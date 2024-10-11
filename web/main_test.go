package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAdhoc(t *testing.T) {
	var form1 Form = Form1{}
	ProcessForm(form1)
}

func ProcessForm(form any) {
	v := reflect.TypeOf(form)
	fmt.Printf("%s\n", v.Name())
	switch f := form.(type) {
	default:
		fmt.Printf("%v\n", f)
	}
}
