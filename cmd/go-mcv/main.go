package main

import (
	"fmt"

	"github.com/ya0201/go-mcv/pkg/nozzle"
)

func main() {
	c, _ := nozzle.HelloNozzle().Pump()
	for c := range c {
		fmt.Println(c.Msg)
	}
}
