package tui

type FieldType int

const (
	FieldText FieldType = iota
	FieldBool
)

type Field struct {
	Name        string
	Placeholder string
	Type        FieldType
}