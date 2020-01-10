package main

import (
	"fmt"
	"reflect"

	"github.com/influxdata/telegraf/plugins/inputs"
	_ "github.com/influxdata/telegraf/plugins/inputs/all"
)

func printPlugin(pad string, ts []reflect.Type) {
	//we're working on the last type in the slice provided
	last := len(ts) - 1
	t := ts[last]

	//recursion helper
	nextPad := pad + " "
	r := func(next reflect.Type) {
		nts := make([]reflect.Type, len(ts)+1) //next type slice
		copy(nts, ts)
		nts[len(ts)] = next
		printPlugin(nextPad, nts)
	}

	fmt.Printf("%s%s", pad, t.Kind())
	name := t.Name()
	if name != "" {
		fmt.Printf(" name: %s", name)
	}

	//detect cycles in the plugin type graph
	for _, d := range ts[:last] {
		if d == t {
			fmt.Printf(" CYCLE\n")
			return
		}
	}

	fmt.Printf("\n") //end first line

	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		r(t.Elem())
	case reflect.Map:
		r(t.Key())
		r(t.Elem())
	case reflect.Struct:
		//for each field, print name and tags then recurse
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if sf.PkgPath != "" {
				//unexported (lowercase) field.  ignore
				continue
			}
			fmt.Printf("%sfield %s", pad, sf.Name)
			if sf.Tag != "" {
				fmt.Printf(" tags: %s", sf.Tag)
			}
			fmt.Printf("\n")
			r(sf.Type)
		}
	}
}

func main() {
	for k, v := range inputs.Inputs {
		fmt.Printf("plugins map key: %s\n", k)
		s := v()
		printPlugin("", []reflect.Type{reflect.TypeOf(s)})
	}
}
