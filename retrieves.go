package db

import (
	"errors"
	"reflect"
	"strings"
)

// Find looks for the record with primary key equal to val
func (q *queryable) Find(id interface{}, val interface{}) error {
	return q.EqualTo(q.source.ID.Column(), id).Retrieve(val)
}

func (q *queryable) Retrieve(val interface{}) error {
	query, values := q.source.conn.Dialect.Query(q)
	row := q.source.runQueryRow(query, values)

	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("Must Supply Ptr to Destination")
	}
	value := reflect.ValueOf(val)
	rfltr := reflector{value}
	plan := q.source.mapPlan(rfltr)

	e := row.Scan(plan.Items()...)
	if e == nil {
		e = q.Initialize(val, plan)
		plan.Finalize(val)
	}
	return e
}

func (q *queryable) RetrieveAll(dest interface{}) error {
	query, values := q.source.conn.Dialect.Query(q)
	rows, err := q.source.runQuery(query, values)
	if err != nil {
		return err
	}
	defer rows.Close()

	destVal := reflect.ValueOf(dest)
	destSliceVal := destVal.Elem()
	tempSliceVal := reflect.Zero(destSliceVal.Type())
	element := destSliceVal.Type().Elem()
	vn := reflect.New(element)
	rfltr := reflector{vn}
	plan := q.source.mapPlan(rfltr)
	for rows.Next() {
		err = rows.Scan(plan.Items()...)
		if err != nil {
			return err
		}
		q.Initialize(vn.Interface())
		plan.Finalize(vn.Interface())
		tempSliceVal = reflect.Append(tempSliceVal, vn.Elem())
		rfltr.item = reflect.New(element)
	}
	destSliceVal.Set(tempSliceVal)
	return nil
}

func (q *queryable) Pluck(selection interface{}, val interface{}) error {
	qq := q.Identity().(*queryable)
	switch sv := selection.(type) {
	case string:
		if strings.Index(sv, ".") == -1 {
			sv = q.source.SqlName + "." + sv
		}
		qq.selection = []selector{selector{Formula: sv}}
	case SqlFunc:
		qq.selection = []selector{selector{Formula: sv.String()}}
	}

	query, values := qq.source.conn.Dialect.Query(qq)
	rows, err := qq.source.runQuery(query, values)
	if err != nil {
		return err
	}
	defer rows.Close()

	destVal := reflect.ValueOf(val)
	destSliceVal := destVal.Elem()
	tempSliceVal := reflect.Zero(destSliceVal.Type())
	element := destSliceVal.Type().Elem()
	vn := reflect.New(element)
	for rows.Next() {
		err = rows.Scan(vn.Interface())
		if err != nil {
			return err
		}
		tempSliceVal = reflect.Append(tempSliceVal, vn.Elem())
	}
	destSliceVal.Set(tempSliceVal)

	return err
}
