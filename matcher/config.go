package matcher

import (
	"github.com/nathan-osman/go-pearup/db"
)

// Config stores information needed by the matcher to run the match routine.
type Config struct {
	Conn *db.Conn
}
