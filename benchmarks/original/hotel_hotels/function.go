package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type HotelSearch struct {
	InDate  string  `json:"inDate"`
	OutDate string  `json:"outDate"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func main() {
	p := &HotelSearch{}
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("Hello %v %v %v %v!\n", p.InDate, p.OutDate, p.Lat, p.Lon)
}
