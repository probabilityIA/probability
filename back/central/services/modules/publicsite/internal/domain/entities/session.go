package entities

type TiendaSession struct {
	Name       string
	Email      string
	Phone      string
	Dni        string
	AvatarURL  string
	CustomerID uint
	Address    *CheckoutAddress
}

type CustomerAddress struct {
	ID         uint
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
	IsPrimary  bool
}

type TiendaOrder struct {
	ID          string
	OrderNumber string
	Total       float64
	Currency    string
	Status      string
	CreatedAt   string
	ItemsCount  int
}
