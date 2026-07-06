package config

type Target struct {
	Name   string
	Width  int
	Height int
	Color bool
}

var Targets = map[string]Target{
	"none": {
		Name:   "Original Size",
		Width:  0,
		Height: 0,
		Color: true,
	},
	"kindle-paperwhite-7": {
		Name:   "Kindle Paperwhite Gen 7",
		Width:  1072,
		Height: 1448,
		Color: false,
	},
}