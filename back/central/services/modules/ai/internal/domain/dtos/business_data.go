package dtos

import "time"

type OrderQuery struct {
	Status       string
	From         *time.Time
	To           *time.Time
	CodOnly      bool
	WithoutGuide bool
	WithNovelty  bool
	Limit        int
}
