package generator

import (
	"github.com/Pallinder/go-randomdata"
	"math/rand"
	"time"
)

type Field interface {
	Type() string
	GenerateValue() interface{}
}

type StringField struct {
	length int
}

func (field StringField) Type() string {
	return "string"
}

func (field StringField) GenerateValue() interface{} {
	return randomdata.RandStringRunes(field.length)
}

type IntegerField struct {
	min int
	max int
}

func (field IntegerField) Type() string {
	return "integer"
}

func (field IntegerField) GenerateValue() interface{} {
	return randomdata.Number(field.min, field.max)
}

type FloatField struct {
	min float64
	max float64
}

func (field FloatField) Type() string {
	return "float"
}

func (field FloatField) GenerateValue() interface{} {
	return float64(rand.Intn(int(field.max-field.min))) + field.min + rand.Float64()
}

type DateField struct {
	min time.Time
	max time.Time
}

func (field DateField) Type() string {
	return "date"
}

func (field DateField) ValidBounds() bool {
	return field.min.Before(field.max)
}

func (field DateField) GenerateValue() interface{} {
	format := "2006-01-02"
	return randomdata.FullDateInRange(field.min.Format(format), field.max.Format(format))
}

type DictField struct {
	category string
}

func (field DictField) Type() string {
	return "dict"
}

func (field DictField) GenerateValue() interface{} {
	switch field.category {
	case "last_name":
		return randomdata.LastName()
	case "first_name":
		return randomdata.FirstName(randomdata.RandomGender)
	case "city":
		return randomdata.City()
	case "country":
		return randomdata.Country(randomdata.FullCountry)
	case "state":
		return randomdata.State(randomdata.Small)
	case "street":
		return randomdata.Street()
	case "address":
		return randomdata.Address()
	case "email":
		return randomdata.Email()
	case "zip_code":
		return randomdata.PostalCode("US")
	case "full_name":
		return randomdata.FullName(randomdata.RandomGender)
	case "random_string":
		return randomdata.SillyName()
	default:
		return nil
	}
}
