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
	source     *source
	order      []string
	groupBy    string
	offset     int
	limit      int
	selection  []selector
	joins      []join
	conditions []condition
}

func (q *Queryable) selectorSql() string {
	if len(q.selection) == 0 {
		return "*"
	}
	output := make([]string, len(q.selection))
	for i, selection := range q.selection {
		output[i] = selection.String()
	}
	return strings.Join(output, ", ")
}

func (q *Queryable) conditionSql() (string, []interface{}) {
	ac := &andCondition{q.conditions}
	return ac.Fragment(), ac.Values()
}

func (q *Queryable) joinSql() string {
	if len(q.selection) == 0 {
		return ""
	}
	output := make([]string, len(q.joins))
	for i, join := range q.joins {
		output[i] = join.String()
	}
	return strings.Join(output, " ")
}

func (q *Queryable) endingSql() string {
	var output string
	if queryable.groupBy != "" {
		output += " GROUP BY " + queryable.groupBy
	}
	if len(queryable.order) > 0 {
		output += strings.Join(queryable.order, ", ")
	}
	if queryable.limit != 0 {
		output += " LIMIT " + fmt.Sprint(queryable.limit)
	}
	if queryable.offset != 0 {
		output += " OFFSET " + fmt.Sprint(queryable.offset)
	}

	return output
}

// Identity is the way to clone a Queryable, it is used everywhere
func (q *Queryable) Identity() *Queryable {
	return &Queryable{
		source:     q.source,
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
		&equalCondition{fmt.Sprint(q.source.ID.Column()), val})
	return nq
}
