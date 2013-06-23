package ar

type Config struct {
	StructToTable func(string) string
	TableToStruct func(string) string
	FieldToColumn func(string) string
	ColumnToField func(string) string
	IdName        string
	CreatedColumn string
	UpdatedColumn string
}
