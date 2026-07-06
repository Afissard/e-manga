package config

type Config struct {
	LibraryPath string `json:"library_path"`
}

var AppConfig Config = Config{
	LibraryPath: "library",
}