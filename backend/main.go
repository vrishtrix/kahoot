package main

import (
	"context"
	"github.com/gofiber/contrib/websocket"
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
	var (
		cursor *mongo.Cursor
		err    error
		ctx    = context.Background()
	)

	if cursor, err = quizCollection.Find(ctx, bson.M{}); err != nil {
		log.Fatal("Error getting cursor while fetching quizzes: ", err)
	}

	var quizzes []fiber.Map

	if err = cursor.All(ctx, &quizzes); err != nil {
		log.Fatal("Error getting all quizzes: ", err)
	}

	return c.JSON(quizzes)
}

func main() {
	setupDb()

	var app = fiber.New()

	// Middlewares
	app.Use(cors.New()) // Remove in production!

	// HTTP Request Handlers (For serving frontend)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// API Route Handlers
	app.Get("/api/quizzes", getQuizzes)

	// WebSocket Connection Handler
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("Error reading from websocket: ", err)
				break
			}

			log.Println("Message from client: ", string(msg))

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("Error writing to client: ", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
