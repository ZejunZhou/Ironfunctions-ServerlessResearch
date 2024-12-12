package main

import (
	"context"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Point struct {
	Pid  string  `bson:"hotelId"`
	Plat float64 `bson:"lat"`
	Plon float64 `bson:"lon"`
}

type Hotel struct {
	Id          string   `bson:"id"`
	Name        string   `bson:"name"`
	PhoneNumber string   `bson:"phoneNumber"`
	Description string   `bson:"description"`
	Address     *Address `bson:"address"`
}

type RecommendHotel struct {
	HId    string  `bson:"hotelId"`
	HLat   float64 `bson:"lat"`
	HLon   float64 `bson:"lon"`
	HRate  float64 `bson:"rate"`
	HPrice float64 `bson:"price"`
}

type Address struct {
	StreetNumber string  `bson:"streetNumber"`
	StreetName   string  `bson:"streetName"`
	City         string  `bson:"city"`
	State        string  `bson:"state"`
	Country      string  `bson:"country"`
	PostalCode   string  `bson:"postalCode"`
	Lat          float32 `bson:"lat"`
	Lon          float32 `bson:"lon"`
}
type RoomType struct {
	BookableRate       float64 `bson:"bookableRate"`
	Code               string  `bson:"code"`
	RoomDescription    string  `bson:"roomDescription"`
	TotalRate          float64 `bson:"totalRate"`
	TotalRateInclusive float64 `bson:"totalRateInclusive"`
}

type RatePlan struct {
	HotelId  string    `bson:"hotelId"`
	Code     string    `bson:"code"`
	InDate   string    `bson:"inDate"`
	OutDate  string    `bson:"outDate"`
	RoomType *RoomType `bson:"roomType"`
}

type Reservation struct {
	HotelId      string `bson:"hotelId"`
	CustomerName string `bson:"customerName"`
	InDate       string `bson:"inDate"`
	OutDate      string `bson:"outDate"`
	Number       int    `bson:"number"`
}

type Number struct {
	HotelId string `bson:"hotelId"`
	Number  int    `bson:"numberOfRoom"`
}

const MongoDBURL = "mongodb://pc21.cloudlab.umass.edu"

func initializeUserDatabase(client *mongo.Client) bool {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
		return false
	}

	collection := client.Database("user-db").Collection("user")
	for i := 0; i <= 10000; i++ {
		suffix := strconv.Itoa(i)
		userName := "Cornell_" + suffix
		password := ""
		for j := 0; j < 10; j++ {
			password += suffix
		}

		filter := bson.M{"username": userName}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatalf("Failed to count documents: %v", err)
			return false
		}

		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), User{Username: userName, Password: password})
			if err != nil {
				log.Fatalf("Failed to insert document: %v", err)
				return false
			}
		} else {
			// log.Printf("User %s already exists", userName)
		}
	}

	// indexModel := mongo.IndexModel{
	// 	Keys: bson.M{"username": 1},
	// }
	// _, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	// if err != nil {
	// 	log.Fatalf("Failed to create index: %v", err)
	// 	return false
	// }

	return true
}

func initializeGeoDatabase(client *mongo.Client) bool {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("geo-db").Collection("geo")
	count, err := collection.CountDocuments(context.TODO(), bson.M{"hotelId": "1"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"1", 37.7867, -122.4112})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "2"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"2", 37.7854, -122.4005})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "3"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"3", 37.7854, -122.4071})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "4"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"4", 37.7936, -122.3930})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "5"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"5", 37.7831, -122.4181})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "6"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = collection.InsertOne(context.TODO(), Point{"6", 37.7863, -122.4015})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 7; i <= 1000; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": hotel_id})
		if err != nil {
			log.Fatal(err)
		}
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		if count == 0 {
			_, err = collection.InsertOne(context.TODO(), Point{hotel_id, lat, lon})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{"hotelId": 1},
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func initializeCollectionDatabase(client *mongo.Client) bool {
	collection := client.Database("geo-db").Collection("hotels")

	hotels := []Hotel{
		{
			Id:          "1",
			Name:        "Clift Hotel",
			PhoneNumber: "(415) 775-4700",
			Description: "A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
			Address: &Address{
				StreetNumber: "495",
				StreetName:   "Geary St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94102",
				Lat:          37.7867,
				Lon:          -122.4112,
			},
		},
		{
			Id:          "2",
			Name:        "W San Francisco",
			PhoneNumber: "(415) 777-5300",
			Description: "Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
			Address: &Address{
				StreetNumber: "181",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7854,
				Lon:          -122.4005,
			},
		},
		{
			Id:          "3",
			Name:        "Hotel Zetta",
			PhoneNumber: "(415) 543-8555",
			Description: "A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
			Address: &Address{
				StreetNumber: "55",
				StreetName:   "5th St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7834,
				Lon:          -122.4071,
			},
		},
		{
			Id:          "4",
			Name:        "Hotel Vitale",
			PhoneNumber: "(415) 278-3700",
			Description: "This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
			Address: &Address{
				StreetNumber: "8",
				StreetName:   "Mission St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94105",
				Lat:          37.7936,
				Lon:          -122.3930,
			},
		},
		{
			Id:          "5",
			Name:        "Phoenix Hotel",
			PhoneNumber: "(415) 776-1380",
			Description: "Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It’s a 4-minute walk from the historic Great American Music Hall nightclub.",
			Address: &Address{
				StreetNumber: "601",
				StreetName:   "Eddy St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7831,
				Lon:          -122.4181,
			},
		},
		{
			Id:          "6",
			Name:        "St. Regis San Francisco",
			PhoneNumber: "(415) 284-4000",
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: &Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7863,
				Lon:          -122.4015,
			},
		},
	}

	for _, hotel := range hotels {
		filter := bson.M{"id": hotel.Id}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), hotel)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for i := 7; i <= 1000; i++ {
		hotel_id := strconv.Itoa(i)
		filter := bson.M{"id": hotel_id}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		phone_num := "(415) 284-40" + hotel_id
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), Hotel{
				Id:          hotel_id,
				Name:        "St. Regis San Francisco",
				PhoneNumber: phone_num,
				Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
				Address: &Address{
					StreetNumber: "125",
					StreetName:   "3rd St",
					City:         "San Francisco",
					State:        "CA",
					Country:      "United States",
					PostalCode:   "94109",
					Lat:          float32(lat),
					Lon:          float32(lon),
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{"id": 1},
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func initializeInventoryDatabase(client *mongo.Client) bool {
	collection := client.Database("rate-db").Collection("inventory")

	ratePlans := []RatePlan{
		{
			HotelId: "1",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &RoomType{
				BookableRate:       109.00,
				Code:               "KNG",
				RoomDescription:    "King sized bed",
				TotalRate:          109.00,
				TotalRateInclusive: 123.17,
			},
		},
		{
			HotelId: "2",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &RoomType{
				BookableRate:       139.00,
				Code:               "QN",
				RoomDescription:    "Queen sized bed",
				TotalRate:          139.00,
				TotalRateInclusive: 153.09,
			},
		},
		{
			HotelId: "3",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &RoomType{
				BookableRate:       109.00,
				Code:               "KNG",
				RoomDescription:    "King sized bed",
				TotalRate:          109.00,
				TotalRateInclusive: 123.17,
			},
		},
	}

	for _, ratePlan := range ratePlans {
		filter := bson.M{"hotelId": ratePlan.HotelId}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), ratePlan)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for i := 7; i <= 1000; i++ {
		if i%3 == 0 {
			hotel_id := strconv.Itoa(i)
			filter := bson.M{"hotelId": hotel_id}
			count, err := collection.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
			}
			end_date := "2015-04-"
			rate := 109.00
			rate_inc := 123.17
			if i%2 == 0 {
				end_date = end_date + "17"
			} else {
				end_date = end_date + "24"
			}

			if i%5 == 1 {
				rate = 120.00
				rate_inc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rate_inc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rate_inc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rate_inc = 258.00
			}

			if count == 0 {
				_, err := collection.InsertOne(context.TODO(), RatePlan{
					HotelId: hotel_id,
					Code:    "RACK",
					InDate:  "2015-04-09",
					OutDate: end_date,
					RoomType: &RoomType{
						BookableRate:       rate,
						Code:               "KNG",
						RoomDescription:    "King sized bed",
						TotalRate:          rate,
						TotalRateInclusive: rate_inc,
					},
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{"hotelId": 1},
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func initializeRecommendationDatabase(client *mongo.Client) bool {
	collection := client.Database("recommendation-db").Collection("recommendation")

	recommendHotels := []RecommendHotel{
		{"1", 37.7867, -122.4112, 109.00, 150.00},
		{"2", 37.7854, -122.4005, 139.00, 120.00},
		{"3", 37.7834, -122.4071, 109.00, 190.00},
		{"4", 37.7936, -122.3930, 129.00, 160.00},
		{"5", 37.7831, -122.4181, 119.00, 140.00},
		{"6", 37.7863, -122.4015, 149.00, 200.00},
	}

	for _, hotel := range recommendHotels {
		filter := bson.M{"hotelId": hotel.HId}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), hotel)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for i := 7; i <= 1000; i++ {
		hotel_id := strconv.Itoa(i)
		filter := bson.M{"hotelId": hotel_id}
		count, err := collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		rate := 135.00
		rate_inc := 179.00
		if i%3 == 0 {
			switch i % 5 {
			case 0:
				rate = 109.00
				rate_inc = 123.17
			case 1:
				rate = 120.00
				rate_inc = 140.00
			case 2:
				rate = 124.00
				rate_inc = 144.00
			case 3:
				rate = 132.00
				rate_inc = 158.00
			case 4:
				rate = 232.00
				rate_inc = 258.00
			}
		}
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), RecommendHotel{hotel_id, lat, lon, rate, rate_inc})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{"hotelId": 1},
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func initializeReservationDatabase(client *mongo.Client) bool {
	collection := client.Database("reservation-db").Collection("reservation")

	count, err := collection.CountDocuments(context.TODO(), bson.M{"hotelId": "4"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Reservation{"4", "Alice", "2015-04-09", "2015-04-10", 1})
		if err != nil {
			log.Fatal(err)
		}
	}

	collection = client.Database("reservation-db").Collection("number")
	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "1"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"1", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "2"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"2", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "3"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"3", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "4"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"4", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "5"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"5", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": "6"})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err := collection.InsertOne(context.TODO(), Number{"6", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 7; i <= 1000; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = collection.CountDocuments(context.TODO(), bson.M{"hotelId": hotel_id})
		if err != nil {
			log.Fatal(err)
		}
		room_num := 200
		if i%3 == 1 {
			room_num = 300
		} else if i%3 == 2 {
			room_num = 250
		}
		if count == 0 {
			_, err := collection.InsertOne(context.TODO(), Number{hotel_id, room_num})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{"hotelId": 1},
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func main() {
	clientOptions := options.Client().ApplyURI(MongoDBURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)

	}

	defer client.Disconnect(context.Background())

	res := initializeUserDatabase(client)
	if !res {
		log.Fatal("Failed to initialize user database.")
	}
	// initializeGeoDatabase(client)
	// initializeCollectionDatabase(client)
	// initializeInventoryDatabase(client)
	// initializeRecommendationDatabase(client)
	// initializeReservationDatabase(client)

	log.Println("Database initialization completed.")
}
