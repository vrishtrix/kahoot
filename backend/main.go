package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var quizCollection *mongo.Collection

func setupDb() {
	var client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	quizCollection = client.Database("quiz").Collection("quizzes")
}

func getQuizzes(c *fiber.Ctx) error {
	var ctx = context.Background()

	var cursor, err = quizCollection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal("Error getting cursor while fetching quizzes: ", err)
	}

	var quizzes []fiber.Map
	err = cursor.All(ctx, &quizzes)

	if err != nil {
		log.Fatal("Error fetching quizzes: ", err)
	}

	return c.JSON(quizzes)
}

func main() {
	var app = fiber.New()

	setupDb()
	app.Use(cors.New()) // Remove in production!

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// API Route Handlers
	app.Get("/api/quizzes", getQuizzes)

	log.Fatal(app.Listen(":3000"))
}
