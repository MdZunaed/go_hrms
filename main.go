package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := connectDb(); err != nil {
		log.Fatal("connect DB error", err)
	}
	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection(EmployeeCollection).Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		//options.Find(c.context())
		var employees []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees)
	})
	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection(EmployeeCollection)

		employee := new(Employee)

		err := c.BodyParser(&employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		employee.Id = ""
		result, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		id := bson.D{{Key: "_id", Value: result.InsertedID}}
		createdData := collection.FindOne(c.Context(), id)

		createdEmployee := &employee
		createdData.Decode(createdEmployee)

		return c.Status(201).JSON(createdEmployee)

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		employeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(404).SendString("id not found")
		}
		collection := mg.Db.Collection(EmployeeCollection)

		employee := new(Employee)

		err = c.BodyParser(&employee)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		query := bson.D{{Key: "_id", Value: employeId}}

		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}
		err = collection.FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString("Employee not found")
			}
			return c.Status(500).SendString(err.Error())
		}
		employee.Id = id
		return c.Status(201).JSON(employee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.SendStatus(404)
		}
		query := bson.D{
			{
				Key:   "_id",
				Value: id,
			},
		}
		result, err := mg.Db.Collection(EmployeeCollection).DeleteOne(c.Context(), query)
		if err != nil {
			return c.SendStatus(500)
		}
		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("Record deleted")
	})

	app.Listen(":3000")
}

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "hrms"
const MongoURI = "mongodb://localhost:27017/" + dbName
const EmployeeCollection = "employees"

func connectDb() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoURI))
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

type Employee struct {
	Id     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}
