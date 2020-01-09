package main

import (
	"fmt"
	"reflect"

	"github.com/influxdata/telegraf/plugins/inputs"
	_ "github.com/influxdata/telegraf/plugins/inputs/all"
)

func printPlugin(pad string, t reflect.Type) {
	n := pad + " "

	fmt.Printf("%s-type %s (kind %s)\n", pad, t.Name(), t.Kind())

	switch t.Kind() {
	case reflect.Ptr:
		printPlugin(n, t.Elem())
	case reflect.Map:
		//for key and elem
		printPlugin(pad+" ", t.Key())
		printPlugin(pad+" ", t.Elem())
	case reflect.Struct:
		//for each field, print name and tags then recurse
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			fmt.Printf("%sfield %s, tags: %s\n", pad, sf.Name, sf.Tag)
			printPlugin(pad+" ", sf.Type)
		}
	default:
		//fmt.Printf("unhandled kind: %s\n", v.Kind())
		fmt.Printf("%sleaf\n", pad)
	}
}

func main() {
	for k, v := range inputs.Inputs {
		fmt.Printf("plugins map key: %s\n", k)
		s := v()
		printPlugin("", reflect.TypeOf(s))
	}
}
