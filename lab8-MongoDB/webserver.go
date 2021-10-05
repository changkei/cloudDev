package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

type Post struct {
	ID    primitive.ObjectID `bson:"_id"`
	Item  string             `bson:"item"`
	Price float32            `bson:"price"`
}

func main() {

	// create a mongo client
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	list := func(w http.ResponseWriter, req *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10000*time.Second)
		err = client.Connect(ctx)
		// Disconnect
		defer client.Disconnect(ctx)
		// select collection from database
		col := client.Database("webserver").Collection("items")

		cur, err := col.Find(ctx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			result := Post{}

			err := cur.Decode(&result)
			if err != nil {
				log.Fatal()
			}
			// do something with result
			fmt.Fprintln(w, "Item: ", result.Item)
			fmt.Fprintln(w, "Price: ", result.Price)
		}
	}

	price := func(w http.ResponseWriter, req *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		// Disconnect
		defer client.Disconnect(ctx)
		// select collection from database
		col := client.Database("shop").Collection("items")

		item := req.URL.Query().Get("item")
		cur, err := col.Find(ctx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			result := Post{}

			err := cur.Decode(&result)
			if err != nil {
				log.Fatal()
			}
			// do something with result
			if result.Item == item {
				fmt.Fprintln(w, "Price: ", result.Price)
			} else {
				fmt.Fprintf(w, "Please enter a valid item: ")
			}
		}
		// Disconnect
		defer client.Disconnect(ctx)
	}

	create := func(w http.ResponseWriter, req *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		// select collection from database
		col := client.Database("shop").Collection("items")

		item := req.URL.Query().Get("item")
		price := req.URL.Query().Get("price")
		if priceDol, err := strconv.ParseFloat(price, 2); err == nil {
			res, err := col.InsertOne(ctx, &Post{
				ID:    primitive.NewObjectID(),
				Item:  item,
				Price: float32(priceDol),
			})
			fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
			if err != nil {
				log.Fatal()
			}
			fmt.Fprintln(w, "Created Item: ", item)
		} else {
			fmt.Fprintf(w, "Please enter a valid number: %q\n", price)
		}
		// Disconnect
		defer client.Disconnect(ctx)
	}

	update := func(w http.ResponseWriter, req *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		// select collection from database
		col := client.Database("shop").Collection("items")

		item := req.URL.Query().Get("item")
		price := req.URL.Query().Get("price")
		// id, _ := primitive.ObjectIDFromHex("")
		if priceDol, err := strconv.ParseFloat(price, 2); err == nil {
			_, err := col.UpdateOne(
				ctx,
				bson.M{"item": item},
				bson.D{
					{"$set", bson.D{{"price", priceDol}}},
				},
			)
			if err != nil {
				log.Println(err)
				log.Fatal()
			}
			fmt.Fprintln(w, "Updated Item: ", item)
		} else {
			fmt.Fprintf(w, "Please enter a valid number: %q\n", price)
		}
		// Disconnect
		defer client.Disconnect(ctx)
	}

	delete := func(w http.ResponseWriter, req *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		// select collection from database
		col := client.Database("shop").Collection("items")

		item := req.URL.Query().Get("item")

		db := Post{}
		errf := col.FindOne(ctx, bson.M{"item": item}).Decode(&db)
		if errf != nil {
			fmt.Fprintf(w, "No such item\n")
		} else {
			res, err := col.DeleteOne(ctx, bson.M{"item": item})
			if err != nil {
				log.Fatal()
			}
			fmt.Printf("DeleteOne removed %v document(s)\n", res.DeletedCount)
			fmt.Fprintf(w, "Deleted %s\n", item)
		}
		// Disconnect
		defer client.Disconnect(ctx)
	}

	http.HandleFunc("/list", list)
	http.HandleFunc("/price", price)
	http.HandleFunc("/create", create)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
