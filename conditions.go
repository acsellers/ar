package ar

type condition interface {
	String() string
	Values() []interface{}
}

type inCondition struct {
	column string
	vals   []interface{}
}

type betweenCondition struct {
	column       string
	lower, upper interface{}
}

type equalCondition struct {
	column string
	val    interface{}
}

type whereCondition struct {
	fragment string
	args     []interface{}
}
