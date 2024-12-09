package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Reservation struct {
	InDate       string  `json:"inDate"`
	OutDate      string  `json:"outDate"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	HotelId      string  `json:"hotelId"`
	CustomerName string  `json:"customerName"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Number       int     `json:"number"`
}

func main() {

	for {
		res := http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: 200,
			Status:     "OK",
		}

		r := bufio.NewReader(os.Stdin)
		req, err := http.ReadRequest(r)
		var buf bytes.Buffer
		if err != nil {
			res.StatusCode = 500
			res.Status = http.StatusText(res.StatusCode)
			fmt.Fprintln(&buf, err)
		} else {
			l, _ := strconv.Atoi(req.Header.Get("Content-Length"))
			p := make([]byte, l)
			r.Read(p)
			hotel := &Reservation{}
			if err := json.Unmarshal(p, hotel); err != nil {
				res.StatusCode = 400
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			}
			fmt.Fprintf(&buf, "Reservation: \nhotel.InDate: %s\nhotel.OutDate: %s\nhotel.Lat: %f\nhotel.Lon: %f\nhotel.HotelId: %s\nhotel.CustomerName: %s\nhotel.Username: %s\nhotel.Password: %s\nhotel.Number: %d\n", hotel.InDate, hotel.OutDate, hotel.Lat, hotel.Lon, hotel.HotelId, hotel.CustomerName, hotel.Username, hotel.Password, hotel.Number)

			// for k, vs := range req.Header {
			// 	fmt.Fprintf(&buf, "ENV: %s %#v\n", k, vs)
			// }
		}
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
	}

	// 1. check user grpc
	// 2. make reservation grpc

	// 2.1. mongodb reservation-db reservation

	// 2.2. mongodb reservation-db number

}
