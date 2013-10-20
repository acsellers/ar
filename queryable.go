package db

import "errors"
import "fmt"
import "reflect"
import "strings"

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

type Formula string

func (q *Queryable) selectorSql() string {
	if len(q.selection) == 0 {
		return strings.Join(q.source.selectColumns(), ", ")
	}
	output := make([]string, len(q.selection))
	for i, selection := range q.selection {
		output[i] = selection.String()
	}
	return strings.Join(output, ", ")
}

func (q *Queryable) conditionSql() (string, []interface{}) {
	if len(q.conditions) > 0 {
		ac := &andCondition{q.conditions}
		return ac.Fragment(), ac.Values()
	}
	return "", []interface{}{}
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

func (queryable *Queryable) endingSql() string {
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

func (q *Queryable) Retrieve(val interface{}) error {
	query, values := q.source.conn.Dialect.Query(q)
	rows, err := q.source.runQuery(query, values)
	defer rows.Close()
	if err != nil {
		return err
	}

	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("Must Supply Ptr to Destination")
	}
	value := reflect.ValueOf(val)
	rfltr := reflector{value}
	plan := q.source.mapPlan(rfltr)

	rows.Next()
	return rows.Scan(plan.Items()...)
}

func (q *Queryable) RetrieveAll(dest interface{}) error {
	query, values := q.source.conn.Dialect.Query(q)
	rows, err := q.source.runQuery(query, values)
	defer rows.Close()
	if err != nil {
		return err
	}
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("Must Supply Ptr to Destination")
	}
	destVal := reflect.ValueOf(dest)
	destSliceVal := destVal.Elem()
	tempSliceVal := destSliceVal
	element := destSliceVal.Type().Elem()
	vn := reflect.New(element)
	rfltr := reflector{vn}
	plan := q.source.mapPlan(rfltr)
	for rows.Next() {
		err = rows.Scan(plan.Items()...)
		if err != nil {
			return err
		}
		tempSliceVal = reflect.Append(tempSliceVal, vn.Elem())
		rfltr.item = reflect.New(element)
	}
	destSliceVal.Set(tempSliceVal)

	return nil
}

func (q *Queryable) UpdateAttribute(column string, val interface{}) error {
	query, vals := q.source.conn.Dialect.Update(q, map[string]interface{}{column: val})
	_, err := q.source.runExec(query, vals)

	return err
}
func (q *Queryable) UpdateAttributes(values map[string]interface{}) error {
	query, vals := q.source.conn.Dialect.Update(q, values)
	_, err := q.source.runExec(query, vals)
	return err
}
