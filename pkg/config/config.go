package config

// logging holds the information about logrus
type logging struct {
	Level string
}

// server holds the information about gin web server
type server struct {
	Debug bool
	StaticPath string
	ServerPort int
	ServerHost string
	RandomCodeLength int
	TokenTimeToLive int
}

// database holds the information about mongo
type database struct {
	Uri string
	DatabaseName string
}

var (
	Logging = &logging{}

	Server = &server{}

	Database = &database{}
)