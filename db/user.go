package db

// User represents an authenticated Facebook user.
type User struct {
	ID              int64
	FacebookID      string `gorm:"not null;unique_index"`
	FacebookName    string `gorm:"not null"`
	FacebookPicture string `gorm:"not null"`
	IsAdmin         bool   `gorm:"not null"`

	// Profile data
	Link     string `gorm:"not null"`
	IsMale   bool   `gorm:"not null"`
	IsFemale bool   `gorm:"not null"`
}

func (u *User) String() string {
	return u.FacebookName
}
