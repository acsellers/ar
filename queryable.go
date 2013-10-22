package db

import "errors"
import "fmt"
import "reflect"
import "strings"

// Conditions for Cond calls
const (
	EQUAL = iota
	NOT_EQUAL
	LESS_THAN
	LESS_OR_EQUAL
	GREATER_THAN
	GREATER_OR_EQUAL
)

// Shorthand conditions for Cond calls
const (
	EQ = iota
	NE
	LT
	LTE
	GT
	GTE
)

type queryable struct {
	*source
	order      []string
	groupBy    string
	offset     int
	limit      int
	selection  []selector
	joins      []join
	conditions []condition
}

type Formula string

func (q *queryable) SelectorSql() string {
	if len(q.selection) == 0 {
		return strings.Join(q.source.selectColumns(), ", ")
	}
	output := make([]string, len(q.selection))
	for i, selection := range q.selection {
		output[i] = selection.String()
	}
	return strings.Join(output, ", ")
}

func (q *queryable) ConditionSql() (string, []interface{}) {
	if len(q.conditions) > 0 {
		ac := &andCondition{q.conditions}
		return ac.Fragment(), ac.Values()
	}
	return "", []interface{}{}
}

func (q *queryable) JoinSql() string {
	if len(q.selection) == 0 {
		return ""
	}
	output := make([]string, len(q.joins))
	for i, join := range q.joins {
		output[i] = join.String()
	}
	return strings.Join(output, " ")
}

func (queryable *queryable) EndingSql() string {
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

// Identity is the way to clone a queryable, it is used everywhere
func (q *queryable) Identity() Scope {
	return &queryable{
		source:     q.source,
		order:      q.order,
		offset:     q.offset,
		limit:      q.limit,
		conditions: q.conditions,
	}
}

func (q *queryable) Where(fragment string, args ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, &whereCondition{fragment, args})
	return nq
}

func (q *queryable) EqualTo(column string, val interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, &equalCondition{column, val})
	return nq
}

func (q *queryable) Between(column string, lower, upper interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, &betweenCondition{column, lower, upper})
	return nq
}

func (q *queryable) In(column string, vals []interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, &inCondition{column, vals})
	return nq
}

func (q *queryable) Cond(column string, condition int, val ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, &varyCondition{column, condition, val})

	return nq
}

func (q *queryable) Limit(limit int) Scope {
	nq := q.Identity().(*queryable)
	nq.limit = limit
	return nq
}

func (q *queryable) Offset(offset int) Scope {
	nq := q.Identity().(*queryable)
	nq.offset = offset
	return nq
}

func (q *queryable) OrderBy(column, direction string) Scope {
	nq := q.Identity().(*queryable)
	nq.order = append(nq.order, column+" "+direction)
	return nq
}

func (q *queryable) Order(ordering string) Scope {
	nq := q.Identity().(*queryable)
	if !(strings.HasSuffix(ordering, "DESC") || strings.HasSuffix("ASC")) {
		ordering = ordering + " ASC"
	}
	nq.order = append(nq.order, ordering)
	return nq
}

func (q *queryable) Reorder(ordering string) Scope {
	nq := q.Identity().(*queryable)
	nq.order = []string{}
	return nq.Order(ordering)
}

// Find looks for the record with primary key equal to val
func (q *queryable) Find(id interface{}, val interface{}) error {
	return q.EqualTo(q.source.ID.Column(), id).Retrieve(val)
	//	nq := q.Identity().(*queryable)
	//	nq.conditions = append(q.conditions,
	//		&equalCondition{fmt.Sprint(q.source.ID.Column()), val})
	//	return nq.Retrieve(val)
}

func (q *queryable) Retrieve(val interface{}) error {
	query, values := q.source.conn.Dialect.Query(q)
	rows, err := q.source.runQuery(query, values)
	if err != nil {
		return err
	}
	defer rows.Close()

	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("Must Supply Ptr to Destination")
	}
	value := reflect.ValueOf(val)
	rfltr := reflector{value}
	plan := q.source.mapPlan(rfltr)

	rows.Next()
	return rows.Scan(plan.Items()...)
}

func (q *queryable) RetrieveAll(dest interface{}) error {
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

func (q *queryable) UpdateAttribute(column string, val interface{}) error {
	query, vals := q.source.conn.Dialect.Update(q, map[string]interface{}{column: val})
	_, err := q.source.runExec(query, vals)

	return err
}
func (q *queryable) UpdateAttributes(values map[string]interface{}) error {
	query, vals := q.source.conn.Dialect.Update(q, values)
	_, err := q.source.runExec(query, vals)
	return err
}
func (q *queryable) UpdateSql(sql string, vals ...interface{}) error {
	panic("UNIMPLEMENTED")
}
func (q *queryable) Delete() error {
	query, vals := q.source.conn.Dialect.Delete(q)
	_, err := q.source.runExec(query, vals)
	return err
}
