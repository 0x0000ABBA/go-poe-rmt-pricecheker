package main

import (
	"context"
	"log"
	"poe-rmt-pricechecker/pricecheckers"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const URI = "mongodb://root:qweasdzxc1337@192.168.1.158:27017"

var funpayCollection *mongo.Collection
var g2gCollection *mongo.Collection
var eldoradoCollection *mongo.Collection

type Record struct {
	ID        primitive.ObjectID `bson:"_id"`
	Timestamp int64             `bson:"created_at"`
	Prices    []float64          `bson:"prices"`
}

var ctx = context.TODO()

func init() {

	clientOptions := options.Client().ApplyURI(URI)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	funpayCollection = client.Database("poe-rmt-prices").Collection("funpay")
	g2gCollection = client.Database("poe-rmt-prices").Collection("g2g")
	eldoradoCollection = client.Database("poe-rmt-prices").Collection("eldorado")

}

func main() {

	go func(){
		for {
			c := colly.NewCollector(
				colly.AllowedDomains("funpay.com", "www.g2g.com"),
			)
			timestamp := time.Now().Unix()
			funpayPrices, err := pricecheckers.GetFunpayPrices(c)
			if err != nil {
				log.Print(err)
			}
			funpayRecord := &Record{
				ID:        primitive.NewObjectID(),
				Timestamp: timestamp,
				Prices:    funpayPrices,
			}
			err = createRecord(funpayRecord, funpayCollection)
			if err != nil {
				log.Print(err)
			}
			g2gPrices, err := pricecheckers.GetG2GPrices(c)
			g2gRecord := &Record{
				ID:        primitive.NewObjectID(),
				Timestamp: timestamp,
				Prices:    g2gPrices,
			}
			if err != nil {
				log.Print(err)
			}
			err = createRecord(g2gRecord, g2gCollection)
			if err != nil {
				log.Print(err)
			}
			eldoradoPrices, err := pricecheckers.GetEldoradoPrices()
			if err != nil {
				log.Print(err)
			}
			eldoradoRecord := &Record{
				ID:        primitive.NewObjectID(),
				Timestamp: timestamp,
				Prices:    eldoradoPrices,
			}
			err = createRecord(eldoradoRecord, eldoradoCollection)
			if err != nil {
				log.Print(err)
			}
			time.Sleep(time.Minute * 10)
		}
	}()

}

func createRecord(record *Record, collection *mongo.Collection) error {
	_, err := collection.InsertOne(ctx, record)
	return err
}
