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
	var vals []interface{}
	for _, condition := range oc.conditions {
		val := condition.Values()
		if len(val) > 0 {
			vals = append(vals, val...)
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
			vals = append(vals, val...)
		}
	}
	return vals
}

type inCondition struct {
	column string
	vals   []interface{}
}

func (ic *inCondition) String() string {
	return withVars(ic.column+" IN (?)", []interface{}{ic.vals})
}
func (ic *inCondition) Fragment() string {
	return ic.column + " IN (?)"
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

type varyCondition struct {
	column string
	cond   int
	val    interface{}
}

func (vc *varyCondition) String() string {
	switch vc.cond {
	case EQUAL:
		if isNil(vc.val) {
			return vc.column + " IS NULL"
		}
		return vc.column + " = " + printString(vc.val)
	case NOT_EQUAL:
		if isNil(vc.val) {
			return vc.column + " IS NOT NULL"
		}
		return vc.column + " <> " + printString(vc.val)
	case LESS_THAN:
		return vc.column + " < " + printString(vc.val)
	case LESS_OR_EQUAL:
		return vc.column + " <= " + printString(vc.val)
	case GREATER_THAN:
		return vc.column + " > " + printString(vc.val)
	case GREATER_OR_EQUAL:
		return vc.column + " >= " + printString(vc.val)
	}

	return ""
}

func (vc *varyCondition) Fragment() string {
	switch vc.cond {
	case EQUAL:
		if isNil(vc.val) {
			return vc.column + " IS NULL"
		}
		return vc.column + " = ?"
	case NOT_EQUAL:
		if isNil(vc.val) {
			return vc.column + " IS NOT NULL"
		}
		return vc.column + " <> ?"
	case LESS_THAN:
		return vc.column + " < ?"
	case LESS_OR_EQUAL:
		return vc.column + " <= ?"
	case GREATER_THAN:
		return vc.column + " > ?"
	case GREATER_OR_EQUAL:
		return vc.column + " >= ?"
	}

	return ""
}

func (vc *varyCondition) Values() []interface{} {
	if isNil(vc.val) {
		return []interface{}{}
	}
	return []interface{}{vc.val}
}
func withVars(sqlFragment string, vals []interface{}) string {
	input := strings.Split(sqlFragment, "?")
	output := []byte{}
	for i, subInput := range input {
		output = append(output, []byte(subInput)...)
		if i < len(vals) {
			if isArray(vals[i]) {
				output = append(output, []byte(printArray(vals[i]))...)
			} else {
				output = append(output, []byte(fmt.Sprint(vals[i]))...)
			}
		}
	}
	return string(output)
}

func isNil(v interface{}) bool {
	return v == nil || (couldBeNil(v) && reflect.ValueOf(v).IsNil())
}

func couldBeNil(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Slice || kind == reflect.Ptr || kind == reflect.Map
}

func isArray(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func printArray(v interface{}) string {
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	if vType.Kind() == reflect.Slice || vType.Kind() == reflect.Array {
		output := make([]string, vValue.Len())
		for i := 0; i < vValue.Len(); i++ {
			output[i] = fmt.Sprint(vValue.Index(i).Interface())
		}
		return strings.Join(output, ",")
	}

	return fmt.Sprint(v)
}

func printString(v interface{}) string {
	if s, ok := v.(string); ok {
		return "'" + s + "'"
	}
	return fmt.Sprint(v)
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
	var current, seperator string

	for scanner.Scan() {
		if isBinder(scanner.Text()) {
			output = append(output, current, scanner.Text())
			current = ""
			seperator = ""
		} else {
			current = current + seperator + scanner.Text()
			seperator = " "
		}
	}
	if current != "" {
		output = append(output, current)
	}

	return output
}

func unbind(sqlFragment string) string {
	output := make([]string, 0)

	for _, fragment := range piecewiseSplit(sqlFragment) {
		if isBinder(fragment) {
			output = append(output, "?")
		} else {
			output = append(output, fragment)
		}
	}
	return strings.Join(output, " ")
}

func outputBindsInOrder(sqlFragment string, bindVals interface{}) []interface{} {
	output := make([]interface{}, 0)
	bind := makeBinder(bindVals)

	for _, fragment := range piecewiseSplit(sqlFragment) {
		if isBinder(fragment) {
			output = append(output, bind.Get(fragment))
		}
	}

	return output
}
