package db

import (
	"errors"
	"reflect"
	"strings"
)

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
	e := rows.Scan(plan.Items()...)
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

func (q *queryable) Pluck(column string, val interface{}) error {
	if strings.Index(column, ".") == -1 {
		column = q.source.SqlName + "." + column
	}
	qq := q.Identity().(*queryable)
	qq.selection = []selector{selector{Formula: column}}

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
