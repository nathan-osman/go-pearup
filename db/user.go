package db

// User represents an authenticated Facebook user.
type User struct {
	ID int64

	// Each user authenticates with one or more of these services
	GoogleID    string
	FacebookID  string
	MicrosoftID string

	// This data is filled in with the most recent login method
	Name    string `gorm:"not null"`
	Email   string
	Picture string

	// Management data
	IsAdmin bool `gorm:"not null"`

	// Profile data
	Link     string `gorm:"not null"`
	IsMale   bool   `gorm:"not null"`
	IsFemale bool   `gorm:"not null"`
}

func (u *User) String() string {
	return u.Name
}
