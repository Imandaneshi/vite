package config

type logging struct {
	Level string
}

type server struct {
	Debug bool
	StaticPath string
	ServerPort int
	ServerHost string
}

type database struct {
	Uri string
	DatabaseName string
}

var (
	// Logging holds the information about logging and logrus
	Logging = &logging{}

	// Server holds the information about vite web server
	Server = &server{}

	// Database holds the information about mongo db
	Database = &database{}
)