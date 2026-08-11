package entities

import "time"

type Prospect struct {
	ID                     uint
	Domain                 string
	MeliSellerID           string
	BusinessName           string
	Platform               string
	City                   string
	Phone                  string
	WhatsAppNumber         string
	InstagramHandle        string
	InstagramFollowers     int
	ProductCount           int
	HasCOD                 bool
	NumCouriersMentioned   int
	UsesProbabilityChannel bool
	RobotsOK               bool
	DiscoveryQuery         string
	Score                  float64
	ScoreBreakdown         map[string]interface{}
	Status                 string
	Seen                   bool
	SourceDiscoveredAt     *time.Time
	SourceLastScannedAt    *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type CountByKey struct {
	Key   string
	Count int64
}

type ScoreBucket struct {
	RangeLabel string
	MinScore   float64
	MaxScore   float64
	Count      int64
}

type ProspectStats struct {
	Total          int64
	ByPlatform     []CountByKey
	ByStatus       []CountByKey
	ByCity         []CountByKey
	ScoreBuckets   []ScoreBucket
	HasCODCount    int64
	SeenCount      int64
	MultiChannel   int64
	AverageScore   float64
	TopScoreDomain string
	TopScore       float64
}
