package db

// A Config is the way you can define defaults for the
// database as far as table, column, foreignkey, and primary
// key naming.
type Config struct {

	// A function that would give the default table name for a struct
	// in the database, for instance Rails would take in "Post" and
	// return "posts"
	StructToTable func(structName string) string

	// A function to give a default name for a databse column based
	// on a struct field
	FieldToColumn func(structName, fieldName string) string

	// A function that gives a guess as to what the default name
	// would be for a foreign key field based on the field name from
	// the struct and the struct that was embedded in the larger
	// struct
	ForeignKeyName func(fieldName string, structName string) string

	// The default name for a primary key field, this will turn
	// the name of the struct into the primary key name. Note that
	// this will return the field name for the struct first, and then
	// the database column next
	IdName func(structName string) (string, string)

	// When timestamping is added, CreatedColumn is the default column
	// to set to the current time when creating records in the database
	CreatedColumn string
	// When timestamping is added, UpdatedColumn is the default column
	// to set to the current time when Saving records in the database.
	// Update* calls would not update this.
	UpdatedColumn string
}
