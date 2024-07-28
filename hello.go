package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// globals , TODO move into env vars?
// //var emuuri string = "mongodb://localhost:C2y6yDjf5%2FR%2Bob0N8A7Cgv30VRDJIWEHLM%2B4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw%2FJw%3D%3D@localhost:10255/admin?ssl=true"
const (
	localuri             = "mongodb://localhost:C2y6yDjf5%2FR%2Bob0N8A7Cgv30VRDJIWEHLM%2B4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw%2FJw%3D%3D@localhost:10255/admin?tls=false"
	MongoDBDatabaseEnv   = "MONGODB_DATABASE"
	MongoDBCollectionEnv = "MONGODB_COLLECTION"
)

var (
	database   string
	collection string
)

func main() {
	////client := connect()
	////log.Printf("mongo sessions, %d", client.NumberSessionsInProgress())

	app := fiber.New()
	app.Post("/api/v1/courses", add)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Listen(":3000")
}

// api handler, add
func add(c *fiber.Ctx) error {
	course := new(Course)
	c.BodyParser(course)
	oid, err := createCourse(c, course)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	inserted, err := findCourse(c, oid)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.Status(200).JSON(inserted)
}

// data access, open connection to db
func connect() *mongo.Client {
	database = os.Getenv(MongoDBDatabaseEnv)
	if database == "" {
		log.Fatalf("missing environment variable, %s", MongoDBDatabaseEnv)
	}
	// TODO dynamic collection from route
	collection = os.Getenv(MongoDBCollectionEnv)
	if collection == "" {
		log.Fatalf("missing environment variable, %s", MongoDBCollectionEnv)
	}

	// wait 10 seconds to connect, otherwise give up
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	// skip TLS verify while testing on local emulator
	var cfg = &tls.Config{InsecureSkipVerify: true}
	opts := options.Client().ApplyURI(localuri).SetDirect(true)
	opts.SetTLSConfig(cfg)
	c, err := mongo.NewClient(opts)

	if err = c.Connect(ctx); err != nil {
		log.Fatalf("Connect failed, %v", err)
	}
	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping failed, %v", err)
	}

	return c
}

// data access, insert
func createCourse(c *fiber.Ctx, course *Course) (primitive.ObjectID, error) {
	client := connect()
	defer client.Disconnect(c.Context())

	courseCollection := client.Database(database).Collection(collection)
	r, err := courseCollection.InsertOne(c.Context(), course)
	if err != nil {
		log.Printf("Add course failed ")
		return primitive.NilObjectID, err
	}

	oid := r.InsertedID.(primitive.ObjectID)
	log.Printf("Course added, %s", oid.Hex())
	return oid, nil
}

func findCourse(c *fiber.Ctx, oid primitive.ObjectID) (bson.M, error) {
	client := connect()
	defer client.Disconnect(c.Context())
	courseCollection := client.Database(database).Collection(collection)

	var inserted bson.M
	query := bson.D{{"_id", oid}}

	err := courseCollection.FindOne(c.Context(), query).Decode(&inserted)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

type Course struct {
	ID          primitive.ObjectID `bson:"_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	PublishedAt time.Time          `bson:"Published_at"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	URL         string             `bson:"url"`
}
