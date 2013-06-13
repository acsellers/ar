package ar

import (
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
	for _, condition := range oc.conditions {
		conds[i] = condition.String()
	}
	return "(" + strings.Join(conds, " OR ") + ")"
}
func (oc *orCondition) Fragment() string {
	conds := make([]string, len(oc.conditions))
	for _, condition := range oc.conditions {
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
	for _, condition := range ac.conditions {
		conds[i] = condition.String()
	}
	return "(" + strings.Join(conds, " AND ") + ")"
}
func (ac *andCondition) Fragment() string {
	conds := make([]string, len(ac.conditions))
	for _, condition := range ac.conditions {
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
	if isNil(val) {
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
	if len(args) == 1 && isBindVars(args[0]) {
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

}

func bindedWith(sqlFragment string, bindVals interface{}) string {

}

func unbind(sqlFragment string) string {

}

func outputBindsInOrder(sqlFragment string, bindVals interface{}) string {
}
