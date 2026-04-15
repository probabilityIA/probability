package entities

import "time"

type CustomerAddress struct {
	ID         uint
	CustomerID uint
	BusinessID uint
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
	Latitude   *float64
	Longitude  *float64
	TimesUsed  int
	LastUsedAt time.Time
}
