package domain

type MeliItemDetail struct {
	ID                string
	Status            string
	SubStatus         []string
	AvailableQuantity int
	UserProductID     string
	LogisticType      string
	Tags              []string
	CatalogListing    bool
	SellerSKU         string
	Variations        []MeliItemVariation
}

type MeliItemVariation struct {
	ID                int64
	UserProductID     string
	AvailableQuantity int
	SellerSKU         string
}

type UserProductStock struct {
	Version   string
	Locations []StockLocation
}

type StockLocation struct {
	Type          string
	StoreID       string
	NetworkNodeID string
	Quantity      int
}

type MeliPack struct {
	ID       int64
	OrderIDs []int64
}

type MeliClaim struct {
	ID           int64
	ResourceType string
	ResourceID   int64
	Reason       string
	Status       string
	Messages     []string
}
