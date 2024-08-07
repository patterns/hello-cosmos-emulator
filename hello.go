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

var mg MongoInstance

func main() {
	mg = connect()
	defer mg.C.Disconnect(mg.Ctx)

	app := fiber.New()
	app.Get("/api/v1/users", users)
	app.Post("/api/v1/users", join)
	app.Get("/api/v1/courses", courses)
	app.Post("/api/v1/courses", compose)

	app.Get("/", numSessions)
	app.Listen(":3000")
}
//TODO middleware to validate CF Access JWT

// api handler, users list
func users(c *fiber.Ctx) error {
	return vanilla[User](c, "users")
}

// api handler, courses list
func courses(c *fiber.Ctx) error {
	return vanilla[Course](c, "courses")
}

// null handler, session count
func numSessions(c *fiber.Ctx) error {
	kv := map[string]int{"count": mg.C.NumberSessionsInProgress()}
	return c.JSON(kv)
}

// vanilla handler, list
func vanilla[T Resource](c *fiber.Ctx, collection string) error {
	// get all records as a cursor
	query := bson.D{{}}
	cursor, err := mg.D.Collection(collection).Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var rows []T = make([]T, 0)

	// iterate the cursor and decode each item into desired items
	if err := cursor.All(c.Context(), &rows); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	// return list in JSON format
	return c.JSON(rows)
}

// api handler, compose (course)
func compose(c *fiber.Ctx) error {
	input := new(Course)
	c.BodyParser(input)
	fresh := &Course{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Description: input.Description,
		URL:         input.URL,
		Published:   time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	oid, err := create(c, fresh, "courses")
	if err != nil {
		return c.Status(500).JSON(err)
	}

	saved, err := find(c, oid, "courses")
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.Status(200).JSON(saved)
}

// api handler, join (user)
func join(c *fiber.Ctx) error {
	parsed := new(User)
	c.BodyParser(parsed)
	fresh := &User{
		ID:          primitive.NewObjectID(),
		Email:       parsed.Email,
		Name:        parsed.Name,
		Role:        parsed.Role,
		Created:     time.Now().UTC(),
		Deactivated: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	oid, err := create(c, fresh, "users")
	if err != nil {
		return c.Status(500).JSON(err)
	}

	saved, err := find(c, oid, "users")
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.Status(200).JSON(saved)
}

// ///////////////////////////////////////
// config/database: data access

// open connection to db
func connect() MongoInstance {
	databaseName := os.Getenv(mongoEnv)
	if databaseName == "" {
		log.Fatalf("missing environment variable, %s", mongoEnv)
	}

	// wait 10 seconds to connect, otherwise give up
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	// skip TLS verify while testing on local emulator
	var cfg = &tls.Config{InsecureSkipVerify: true}
	opts := options.Client().ApplyURI(localuri).SetDirect(true)
	opts.SetRetryWrites(false)
	opts.SetTLSConfig(cfg)
	c, err := mongo.NewClient(opts)

	if err = c.Connect(ctx); err != nil {
		log.Fatalf("Connect failed, %v", err)
	}
	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping failed, %v", err)
	}

	return MongoInstance{
		C:   c,
		D:   c.Database(databaseName),
		Ctx: ctx,
	}
}

// data access, insert
func create(c *fiber.Ctx, fresh interface{}, collection string) (primitive.ObjectID, error) {
	coll := mg.D.Collection(collection)
	r, err := coll.InsertOne(c.Context(), fresh)
	if err != nil {
		log.Printf("Add item failed ")
		return primitive.NilObjectID, err
	}
	// assert result ID
	oid := r.InsertedID.(primitive.ObjectID)
	log.Printf("Item added, %s", oid.Hex())
	return oid, nil
}

// data access, get by id
func find(c *fiber.Ctx, oid primitive.ObjectID, collection string) (bson.M, error) {
	coll := mg.D.Collection(collection)

	var found bson.M
	query := bson.D{{"_id", oid}}

	err := coll.FindOne(c.Context(), query).Decode(&found)
	if err != nil {
		return nil, err
	}

	return found, nil
}

// ///////////////////////////////////////
// type constraint for union of courses and users
type Resource interface {
	User | Course
}
type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Created     time.Time          `json:"created" bson:"created"`
	Deactivated time.Time          `json:"deactivated" bson:"deactivated"`
	Email       string             `json:"email" bson:"email"`
	Name        string             `json:"name" bson:"name"`
	Role        string             `json:"role" bson:"role"`
}
type Course struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Published   time.Time          `json:"published" bson:"published"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	URL         string             `json:"url" bson:"url"`
}

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	// todo Is this struct just for namespace?
	C   *mongo.Client
	D   *mongo.Database
	Ctx context.Context
}

// globals , TODO move into env vars?
// //var emuuri string = "mongodb://localhost:C2y6yDjf5%2FR%2Bob0N8A7Cgv30VRDJIWEHLM%2B4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw%2FJw%3D%3D@localhost:10255/admin?ssl=true"
const (
	localuri = "mongodb://localhost:C2y6yDjf5%2FR%2Bob0N8A7Cgv30VRDJIWEHLM%2B4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw%2FJw%3D%3D@localhost:10255/admin?tls=false"
	mongoEnv = "MONGODB_DATABASE"
)
