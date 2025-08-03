package model

type User struct {
	Id       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"-"` // Prevents marshalling to JSON
}

type Employee struct {
	Id     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}
