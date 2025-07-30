package employee

import (
	"github.com/MdZunaed/go_hrms/internal/db"
	model "github.com/MdZunaed/go_hrms/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEmployees() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := db.GetEmployeeCollection().Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		//options.Find(c.context())
		var employees []model.Employee = make([]model.Employee, 0)

		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees)
	}
}

func CreateEmployee() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := db.GetEmployeeCollection()

		employee := new(model.Employee)

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

	}
}

func UpdateEmployee() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		employeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(404).SendString("id not found")
		}
		collection := db.GetEmployeeCollection()

		employee := new(model.Employee)

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
	}
}

func DeleteEmployee() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		result, err := db.GetEmployeeCollection().DeleteOne(c.Context(), query)
		if err != nil {
			return c.SendStatus(500)
		}
		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("Record deleted")
	}
}
