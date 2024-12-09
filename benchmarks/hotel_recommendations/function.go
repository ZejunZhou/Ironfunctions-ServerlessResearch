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

type HotelRecommendation struct {
	Require string  `json:"require"` // dis, rate, price
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	// HLat    float64
	// HLon    float64
}

// should connect to mongodb
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
			hotel := &HotelRecommendation{}
			if err := json.Unmarshal(p, hotel); err != nil {
				res.StatusCode = 400
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			}
			fmt.Fprintf(&buf, "HotelSearch: \nhotel.Require: %s\\nhotel.Lat: %f\nhotel.Lon: %f\n", hotel.Require, hotel.Lat, hotel.Lon)

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
