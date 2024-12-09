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

type HotelSearch struct {
	InDate  string
	OutDate string
	Lat     float64
	Lon     float64
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
			// fmt.Fprintf(&buf, "Hello %s\n", p)
			hotel := &HotelSearch{}
			if err := json.Unmarshal(p, hotel); err != nil {
				res.StatusCode = 400
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			}
			fmt.Fprintf(&buf, "HotelSearch: \nhotel.InDate: %s\nhotel.OutDate: %s\nhotel.Lat: %f\nhotel.Lon: %f\n", hotel.InDate, hotel.OutDate, hotel.Lat, hotel.Lon)

			// fmt.Printf("HotelSearch: %+v\n", hotel)
			// fmt.Printf("inDate: %s\n", hotel.InDate)
			// fmt.Printf("outDate: %s\n", hotel.OutDate)
			// fmt.Printf("lat: %f\n", hotel.Lat)
			// fmt.Printf("lon: %f\n", hotel.Lon)
			// fmt.Printf("lat: %f\n", hotel.Lat)
			// fmt.Printf("lon: %f\n", hotel.Lon)

			// for k, vs := range req.Header {
			// 	fmt.Fprintf(&buf, "ENV: %s %#v\n", k, vs)
			// }
		}
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
	}

	// fmt.Printf("Hotel Recommendation: %+v\n", hotel)

}
