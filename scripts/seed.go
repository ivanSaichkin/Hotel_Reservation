package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GoDev/Hotel-reservatrion/api"
	"github.com/GoDev/Hotel-reservatrion/db"
	"github.com/GoDev/Hotel-reservatrion/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookngStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, pw string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  pw,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		SeaSide: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking: ", resp.ID)
}

func main() {
	james := seedUser(false, "james", "foo", "james@foo.com", "supersecurepassword")
	seedUser(true, "admin", "admin", "admin@admin.com", "adminpassword123")

	seedHotel("Ballucia", "France", 3)
	seedHotel("The cozy", "The Nederlands", 4)
	hotel := seedHotel("Dont die in your sleep", "London", 1)

	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", true, 189.99, hotel.ID)
	room := seedRoom("large", false, 289.99, hotel.ID)

	seedBooking(james.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NemMongoHotelStore(client)
	roomStore = db.NemMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
