package ar

type Queryable struct {
	source *Mapper
}

// Identity is a no-op, but it is kept to maintain parity
// between mapper and queryable
func (q *Queryable) Identity() *Queryable {
	return q
}

func (q *Queryable) Where(fragment string, args ...interface{}) *Queryable {
	return q
}

func (q *Queryable) EqualTo(column string, val interface{}) *Queryable {
	return q
}

func (q *Queryable) Between(column string, lower, upper interface{}) *Queryable {
	return q
}

func (q *Queryable) In(column string, vals []interface{}) *Queryable {
	return q
}

func (q *Queryable) Limit(limit int) *Queryable {
	return q
}

func (q *Queryable) Offset(offset int) *Queryable {
	return q
}

func (q *Queryable) OrderBy(column, direction string) *Queryable {
	return q
}

func (q *Queryable) Order(ordering string) *Queryable {
	return q
}

func (q *Queryable) Reorder(ordering string) *Queryable {
	return q
}

// Find looks for the record with primary key equal to val
func (q *Queryable) Find(val interface{}) *Queryable {
	return q
}
