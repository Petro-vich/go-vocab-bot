package word

import (
	"time"
)

type Word struct {
	Russian   string
	English   string
	CreatedAt time.Time
	Status    string
	Lvl       int
}
