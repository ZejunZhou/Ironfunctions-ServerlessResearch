package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getUserFromSharedMemory(username string) (*User, error) {
	filePath := "/dev/shm/hotel_user"
	err := os.Chown("/dev/shm/hotel_user", os.Getuid(), os.Getgid())
	if err != nil {
		return nil, fmt.Errorf("failed to change owner of /dev/shm/hotel_user: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var users []User
	if stat, _ := file.Stat(); stat.Size() != 0 {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &users); err != nil {
			return nil, err
		}
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, nil
}

func writeUserToSharedMemory(user *User) error {
	filePath := "/dev/shm/hotel_user"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var users []User
	if stat, _ := file.Stat(); stat.Size() != 0 {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &users); err != nil {
			return err
		}
	}

	users = append(users, *user)

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func main() {
	p := &User{}
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("Hello %v!\n", p.Username)
}
