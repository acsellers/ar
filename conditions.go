package ar

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
)

type condition interface {
	String() string
	Fragment() string
	Values() []interface{}
}

type orCondition struct {
	conditions []condition
}

func (oc *orCondition) String() string {
	conds := make([]string, len(oc.conditions))
	for i, condition := range oc.conditions {
		conds[i] = condition.String()
	}
	return "(" + strings.Join(conds, " OR ") + ")"
}
func (oc *orCondition) Fragment() string {
	conds := make([]string, len(oc.conditions))
	for i, condition := range oc.conditions {
		conds[i] = condition.Fragment()
	}
	return "(" + strings.Join(conds, " OR ") + ")"
}
func (oc *orCondition) Values() []interface{} {
	vals := make([]interface{}, 0, len(oc.conditions))
	for _, condition := range oc.conditions {
		val := condition.Values()
		if len(val) > 0 {
			vals = append(vals, val)
		}
	}
	return vals
}

type andCondition struct {
	conditions []condition
}

func (ac *andCondition) String() string {
	conds := make([]string, len(ac.conditions))
	for i, condition := range ac.conditions {
		conds[i] = condition.String()
	}
	return "(" + strings.Join(conds, " AND ") + ")"
}
func (ac *andCondition) Fragment() string {
	conds := make([]string, len(ac.conditions))
	for i, condition := range ac.conditions {
		conds[i] = condition.Fragment()
	}
	return "(" + strings.Join(conds, " AND ") + ")"
}
func (ac *andCondition) Values() []interface{} {
	vals := make([]interface{}, 0, len(ac.conditions))
	for _, condition := range ac.conditions {
		val := condition.Values()
		if len(val) > 0 {
			vals = append(vals, val)
		}
	}
	return vals
}

type inCondition struct {
	column string
	vals   []interface{}
}

func (ic *inCondition) String() string {
	return withVars(ic.column+" IN ?", []interface{}{ic.vals})
}
func (ic *inCondition) Fragment() string {
	return ic.column + " IN [?]"
}
func (ic *inCondition) Values() []interface{} {
	return []interface{}{ic.vals}
}

type betweenCondition struct {
	column       string
	lower, upper interface{}
}

func (bc *betweenCondition) String() string {
	return withVars(bc.Fragment(), bc.Values())
}
func (bc *betweenCondition) Fragment() string {
	return bc.column + " BETWEEN ? AND ?"
}
func (bc *betweenCondition) Values() []interface{} {
	return []interface{}{bc.lower, bc.upper}
}

type equalCondition struct {
	column string
	val    interface{}
}

func (ec *equalCondition) String() string {
	return withVars(ec.Fragment(), ec.Values())
}
func (ec *equalCondition) Fragment() string {
	if isNil(ec.val) {
		return ec.column + "IS NULL"
	}
	return ec.column + " = ?"
}
func (ec *equalCondition) Values() []interface{} {
	return []interface{}{ec.val}
}

type whereCondition struct {
	fragment string
	args     []interface{}
}

func (wc *whereCondition) String() string {
	switch {
	case len(wc.args) == 0:
		return wc.fragment
	case len(wc.args) == 1 && isBindVars(wc.args[0]):
		return bindedWith(wc.fragment, wc.args[0])
	default:
		return withVars(wc.fragment, wc.args)
	}
}
func (wc *whereCondition) Fragment() string {
	if len(wc.args) == 1 && isBindVars(wc.args[0]) {
		return unbind(wc.fragment)
	}

	return wc.fragment
}
func (wc *whereCondition) Values() []interface{} {
	switch {
	case len(wc.args) == 0:
		return []interface{}{}
	case len(wc.args) == 1 && isBindVars(wc.args[0]):
		return outputBindsInOrder(wc.fragment, wc.args[0])
	}

	return wc.args
}

func withVars(sqlFragment string, vals []interface{}) string {
	return ""
}

func isNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

func isBindVars(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func bindedWith(sqlFragment string, bindVals interface{}) string {
	output := make([]string, 0)
	bind := makeBinder(bindVals)

	for _, fragment := range piecewiseSplit(sqlFragment) {
		if isBinder(fragment) {
			output = append(output, fmt.Sprint(bind.Get(fragment)))
		} else {
			output = append(output, fragment)
		}
	}
	return strings.Join(output, " ")
}

func isBinder(fragment string) bool {
	return fragment[0] == ':' && fragment[len(fragment)-1] == ':'
}

func makeBinder(v interface{}) *binder {
	rv := reflect.ValueOf(v)
	b := new(binder)
	if rv.Type().Kind() == reflect.Map {
		b.mapValue = rv
		b.useful = true
	}
	return b
}

type binder struct {
	mapValue reflect.Value
	useful   bool
}

func (b *binder) Get(item string) interface{} {
	if b.useful {
		vv := b.mapValue.MapIndex(reflect.ValueOf(strings.Trim(item, ":")))
		if vv.IsValid() {
			return vv.Interface()
		}
	}
	return item
}

func piecewiseSplit(sqlFragment string) []string {
	scanner := bufio.NewScanner(strings.NewReader(sqlFragment))
	scanner.Split(bufio.ScanWords)

	var output []string
	var current string

	for scanner.Scan() {
		if isBinder(scanner.Text()) {
			output = append(output, current, scanner.Text())
			current = ""
		} else {
			current = current + " " + scanner.Text()
		}
	}
	if current != "" {
		output = append(output, current)
	}

	return output
}

func unbind(sqlFragment string) string {

}

func outputBindsInOrder(sqlFragment string, bindVals interface{}) []interface{} {
}
