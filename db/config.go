package db

// Config stores the connection information for the database.
type Config struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}
