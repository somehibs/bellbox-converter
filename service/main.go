package main

import (
	"fmt"

	"git.circuitco.de/self/bellbox-converter"
)

func main() {
	m, err := converter.LoadTranslations()
	fmt.Printf("Map: %+v Err: %+v\n", m, err)
	c, err := converter.LoadConfig()
	fmt.Printf("Config: %+v Err: %+v\n", c, err)
	fmt.Printf("Run: %+v\n", converter.New())
}
