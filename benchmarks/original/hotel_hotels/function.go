package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type HotelSearch struct {
	InDate  string  `json:"inDate"`
	OutDate string  `json:"outDate"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func getHotelFromSharedMemory() (*HotelSearch, error) {
	filePath := "/dev/shm/hotels"
	err := os.Chown("/dev/shm/hotels", os.Getuid(), os.Getgid())
	if err != nil {
		return nil, fmt.Errorf("failed to change owner of /dev/shm/hotels: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var hotels []HotelSearch
	if err := json.Unmarshal(data, &hotels); err != nil {
		return nil, err
	}

	if len(hotels) == 0 {
		return nil, fmt.Errorf("no hotels found")
	}

	return &hotels[0], nil
}

func writeHotelToSharedMemory(hotel *HotelSearch) error {
	filePath := "/dev/shm/hotels"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var hotels []HotelSearch
	if len(data) > 0 {
		if err := json.Unmarshal(data, &hotels); err != nil {
			return err
		}
	}

	hotels = append(hotels, *hotel)

	newData, err := json.Marshal(hotels)
	if err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	_, err = file.Write(newData)
	return err

}

func main() {
	p := &HotelSearch{}
	json.NewDecoder(os.Stdin).Decode(p)
	if p.InDate != "" {
		if err := writeHotelToSharedMemory(p); err != nil {
			fmt.Println(err)
			return
		}
	}

	hotel, err := getHotelFromSharedMemory()
	if err != nil {
		fmt.Println(err)
	}
	if hotel != nil {
		fmt.Println(hotel)
	}
}
