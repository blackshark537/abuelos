package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Variable struct {
	Id        primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt time.Time          `json:"CreatedAt" xml:"CreatedAt" form:"CreatedAt"`
	UpdatedAt time.Time          `json:"UpdatedAt" xml:"UpdatedAt" form:"UpdatedAt"`
	Name      string
	Value     float64
	Month     time.Month
	Year      int
}

func (v *Variable) Save() (interface{}, error) {
	return nil, errors.New("Not Implemented")
}

func (v *Variable) Update() error {
	return errors.New("Not Implemented")
}

func (v *Variable) Delete(filters string) error {
	return errors.New("Not Implemented")
}

func (v *Variable) GetAll(filters string) ([]Variable, error) {
	return nil, errors.New("Not Implemented")
}

func (v *Variable) FindOne(filters string) error {
	return errors.New("Not Implemented")
}
