package ar

import "fmt"

const (
	EQUAL = iota
	NOT_EQUAL
	LESS_THAN
	LESS_OR_EQUAL
	GREATER_THAN
	GREATER_OR_EQUAL
)

type Queryable struct {
	model      *model
	order      []string
	offset     int
	limit      int
	conditions []condition
}

// Identity is the way to clone a Queryable, it is used everywhere
func (q *Queryable) Identity() *Queryable {
	return &Queryable{
		model:      q.model,
		order:      q.order,
		offset:     q.offset,
		limit:      q.limit,
		conditions: q.conditions,
	}
}

func (q *Queryable) Where(fragment string, args ...interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(nq.conditions, &whereCondition{fragment, args})
	return nq
}

func (q *Queryable) EqualTo(column string, val interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(nq.conditions, &equalCondition{column, val})
	return nq
}

func (q *Queryable) Between(column string, lower, upper interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(nq.conditions, &betweenCondition{column, lower, upper})
	return nq
}

func (q *Queryable) In(column string, vals []interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(nq.conditions, &inCondition{column, vals})
	return nq
}

func (q *Queryable) Cond(column string, condition int, val ...interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(nq.conditions, &varyCondition{column, condition, val})

	return nq
}

func (q *Queryable) Limit(limit int) *Queryable {
	nq := q.Identity()
	nq.limit = limit
	return nq
}

func (q *Queryable) Offset(offset int) *Queryable {
	nq := q.Identity()
	nq.offset = offset
	return nq
}

func (q *Queryable) OrderBy(column, direction string) *Queryable {
	nq := q.Identity()
	nq.order = append(nq.order, column+" "+direction)
	return nq
}

func (q *Queryable) Order(ordering string) *Queryable {
	nq := q.Identity()
	nq.order = append(nq.order, ordering)
	return nq
}

func (q *Queryable) Reorder(ordering string) *Queryable {
	nq := q.Identity()
	nq.order = []string{ordering}
	return nq
}

// Find looks for the record with primary key equal to val
func (q *Queryable) Find(val interface{}) *Queryable {
	nq := q.Identity()
	nq.conditions = append(q.conditions,
		&equalCondition{fmt.Sprint(q.model.pk), val})
	return nq
}
