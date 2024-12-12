package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string
}

func main() {
	p := &Person{Name: "World2"}
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("Hello %v!", p.Name)
	fmt.Printf("OS GET ENV: %v", os.Getenv("FUNC_ID"))
}
