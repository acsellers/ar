package ar

type Mapper struct {
	model *model
}

// Identity will create a queryable for a specific mapper,
// in this way, the different Queryable methods are simply
// implemented in the manner mapper.Identity().Method().
func (m *Mapper) Identity() *Queryable {
	return &Queryable{model: m.model}
}

func (m *Mapper) Where(fragment string, args ...interface{}) *Queryable {
	return m.Identity().Where(fragment, args...)
}

func (m *Mapper) EqualTo(column string, val interface{}) *Queryable {
	return m.Identity().EqualTo(column, val)
}

func (m *Mapper) Between(column string, lower, upper interface{}) *Queryable {
	return m.Identity().Between(column, lower, upper)
}

func (m *Mapper) In(column string, vals []interface{}) *Queryable {
	return m.Identity().In(column, vals)
}

func (m *Mapper) Limit(limit int) *Queryable {
	return m.Identity().Limit(limit)
}

func (m *Mapper) Offset(offset int) *Queryable {
	return m.Identity().Offset(offset)
}

func (m *Mapper) OrderBy(column, direction string) *Queryable {
	return m.Identity().OrderBy(column, direction)
}

func (m *Mapper) Order(ordering string) *Queryable {
	return m.Identity().Order(ordering)
}

func (m *Mapper) Reorder(ordering string) *Queryable {
	return m.Identity().Reorder(ordering)
}

// Find looks for the record with primary key equal to val
func (m *Mapper) Find(val interface{}) *Queryable {
	return m.Identity().Find(val)
}
