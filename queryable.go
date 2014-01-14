package db

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type COND int

// Conditions for Cond calls
const (
	EQUAL COND = iota
	NOT_EQUAL
	LESS_THAN
	LESS_OR_EQUAL
	GREATER_THAN
	GREATER_OR_EQUAL
)

// Shorthand conditions for Cond calls
const (
	EQ COND = iota
	NE
	LT
	LTE
	GT
	GTE
)

/*
This is used to pass multiple columns and values to UpdateAttributes
functions

  // Update multiple attributes
  somePostsScope.UpdateAttributes(db.Attributes{
    "created_at": time.Now(),
    "building": true,
    "state": "build",
  })
*/
type Attributes map[string]interface{}

type queryable struct {
	*source
	order      []string
	groupBy    string
	having     []whereCondition
	offset     int
	limit      int
	selection  []selector
	joins      []*join
	conditions []condition
}

func (q *queryable) SelectorSql() string {
	if len(q.selection) == 0 {
		return strings.Join(q.source.selectColumns(), ", ")
	} else {
		selections := make([]string, len(q.selection))
		for i, selection := range q.selection {
			selections[i] = selection.String()
		}
		return strings.Join(selections, ", ")
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
		return ac.Fragment(), q.cleanValues(ac.Values())
	}
	return "", []interface{}{}
}

// cleanValues takes in an array of values, and exchanges any structs
// for their ids, it prefers mixins
func (q *queryable) cleanValues(vals []interface{}) []interface{} {
	output := make([]interface{}, len(vals))
	for i, val := range vals {
		// Deal with the lone sql struct, time.Time
		if _, ok := val.(*time.Time); ok {
			output[i] = val
			continue
		}

		switch reflect.TypeOf(val).Kind() {
		case reflect.Struct:
			fn := fullNameFor(reflect.TypeOf(val))
			if s, ok := q.conn.mappedStructs[fn]; ok {
				output[i] = s.extractID(reflect.ValueOf(val))
			} else {
				output[i] = val
			}
		default:
			output[i] = val
		}
	}
	return output
}

func (q *queryable) JoinsSql() string {
	if len(q.selection) == 0 {
		return ""
	}
	output := make([]string, len(q.joins))
	for i, join := range q.joins {
		output[i] = join.String()
	}
	return strings.Join(output, " ")
}

func (queryable *queryable) EndingSql() (string, []interface{}) {
	var output string
	var vals []interface{}
	if queryable.groupBy != "" {
		output += " GROUP BY " + queryable.groupBy
	}
	if len(queryable.having) > 0 {
		clauses := []string{}
		for _, h := range queryable.having {
			vals = append(vals, h.Values()...)
			clauses = append(clauses, h.Fragment())
		}
		output += strings.Join(clauses, " AND ")
	}
	if len(queryable.order) > 0 {
		output += " ORDER BY " + strings.Join(queryable.order, ", ")
	}
	if queryable.limit != 0 {
		output += " LIMIT " + fmt.Sprint(queryable.limit)
	}
	if queryable.offset != 0 {
		output += " OFFSET " + fmt.Sprint(queryable.offset)
	}

	return output, vals
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

func (q *queryable) Having(fragment string, args ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.having = append(nq.having, whereCondition{fragment, args})
	return nq
}

func (q *queryable) GroupBy(groupItem string) Scope {
	nq := q.Identity().(*queryable)
	nq.groupBy = groupItem
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

func (q *queryable) In(column string, items interface{}) Scope {
	nq := q.Identity().(*queryable)
	nq.conditions = append(nq.conditions, newInCondition(column, items))
	return nq
}

func (q *queryable) Cond(column string, condition COND, val interface{}) Scope {
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
	if !(strings.HasSuffix(ordering, "DESC") || strings.HasSuffix(ordering, "ASC")) {
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

func (q *queryable) Count() (int64, error) {
	ct := "COUNT(" + q.source.SqlName + "." + q.source.ID.SqlColumn + ")"
	qq := q.Identity().(*queryable)
	qq.selection = []selector{selector{Formula: ct}}

	var count int64
	query, values := qq.source.conn.Dialect.Query(qq)
	row := qq.source.runQueryRow(query, values)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(query)
	}

	return count, err
}

func (q *queryable) UpdateAttribute(column string, val interface{}) error {
	query, vals := q.source.conn.Dialect.Update(q, map[string]interface{}{column: val})
	_, err := q.source.runExec(query, vals)

	return err
}
func (q *queryable) UpdateAttributes(values Attributes) error {
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
func (q *queryable) LeftJoin(joins ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	for _, j := range joins {
		jn, cd := newJoin("LEFT", j, q)
		nq.joins = append(nq.joins, jn...)
		if len(cd) > 0 {
			nq.conditions = append(nq.conditions, cd...)
		}
	}
	return nq
}
func (q *queryable) InnerJoin(joins ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	for _, j := range joins {
		jn, cd := newJoin("INNER", j, q)
		nq.joins = append(nq.joins, jn...)
		if len(cd) > 0 {
			nq.conditions = append(nq.conditions, cd...)
		}
	}
	return nq
}
func (q *queryable) FullJoin(joins ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	for _, j := range joins {
		jn, cd := newJoin("FULL OUTER", j, q)
		nq.joins = append(nq.joins, jn...)
		if len(cd) > 0 {
			nq.conditions = append(nq.conditions, cd...)
		}
	}
	return nq
}
func (q *queryable) RightJoin(joins ...interface{}) Scope {
	nq := q.Identity().(*queryable)
	for _, j := range joins {
		jn, cd := newJoin("RIGHT OUTER", j, q)
		nq.joins = append(nq.joins, jn...)
		if len(cd) > 0 {
			nq.conditions = append(nq.conditions, cd...)
		}
	}
	return nq
}
func (q *queryable) JoinSql(sql string, args ...interface{}) Scope {
	return q.Identity().JoinSql(sql, args...)
}

func (q *queryable) LeftInclude(include ...interface{}) Scope {
	return q.Identity().LeftInclude(include...)
}
func (q *queryable) InnerInclude(include ...interface{}) Scope {
	return q.Identity().InnerInclude(include...)
}
func (q *queryable) FullInclude(include interface{}, nullRecords interface{}) Scope {
	return q.Identity().FullInclude(include, nullRecords)
}
func (q *queryable) RightInclude(include interface{}, nullRecords interface{}) Scope {
	return q.Identity().RightInclude(include, nullRecords)
}
func (q *queryable) IncludeSql(il IncludeList, query string, args ...interface{}) Scope {
	return q.Identity().IncludeSql(il, query, args...)
}
