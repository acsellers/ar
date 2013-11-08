package db

type Config struct {
	StructToTable  func(structName string) string
	FieldToColumn  func(fieldName string) string
	ForeignKeyName func(fieldName string, structName string) string
	IdName         string
	CreatedColumn  string
	UpdatedColumn  string
}
