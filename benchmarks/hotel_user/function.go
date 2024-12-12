package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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
		data, err := io.ReadAll(file)
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
		data, err := io.ReadAll(file)
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

const MongoDBURL = "mongodb://pc21.cloudlab.umass.edu"

func getUserFromMongoDB(username string) (*User, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MongoDBURL))
	if err != nil {
		return nil, err
	}

	collection := client.Database("user-db").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err = collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
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
			user := &User{}
			if err := json.Unmarshal(p, user); err != nil {
				res.StatusCode = 400
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			}

			existingUser, err := getUserFromMongoDB(user.Username)
			if err != nil {
				res.StatusCode = 500
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			} else if existingUser != nil {
				if existingUser.Password != user.Password {
					res.StatusCode = 400
					res.Status = http.StatusText(res.StatusCode)
					fmt.Fprintln(&buf, "Password is incorrect")
				} else if existingUser.Password == user.Password {
					fmt.Fprintf(&buf, "Hello, %s, welcome back!\n", user.Username)
				}
			} else {
				res.StatusCode = 404
				fmt.Fprintf(&buf, "Non-existing user %s, please sign up first\n", user.Username)
			}
			// for k, vs := range req.Header {
			// 	fmt.Fprintf(&buf, "ENV: %s %#v\n", k, vs)
			// }
		}

		res.Body = io.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
	}
}
