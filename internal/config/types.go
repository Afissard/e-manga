package config

const (
	APPCONFIG_FILE = "config.json"

	SOURCE_DIR_NAME    = "source"
	CACHE_DIR_NAME     = "cache"
	OUTPUT_DIR_NAME    = "output"
	METADATA_FILE_NAME = "metadata.json"
)

var (
	Configuration = Config{}
	LogSrv        *Logger
)

type Config struct {
	LibraryPath string   `json:"library_path"`
	TuiMode     bool     `json:"tui_mode"`
	Targets     []Target `json:"targets"`
}

type Target struct {
	Name       string `json:"name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Color      bool   `json:"color"`
	AutoRotate bool   `json:"auto_rotate"`
}
