package db

import (
	"database/sql"
	"fmt"
	"reflect"
)

type planner struct {
	scanners []*reflectScanner
}

func (s *source) mapPlan(v reflector) *planner {
	p := &planner{[]*reflectScanner{}}

	for _, col := range s.Fields {
		if col.columnInfo != nil && col.structOptions != nil {
			p.scanners = append(
				p.scanners,
				&reflectScanner{parent: v, column: col},
			)
		}
	}

	return p
}

func (s *source) selectColumns() []string {
	output := []string{}
	for _, col := range s.Fields {
		if col.columnInfo != nil && col.structOptions != nil {
			output = append(
				output,
				fmt.Sprintf("%s.%s", s.SqlName, col.columnInfo.SqlColumn),
			)
		}
	}
	return output
}

func (p *planner) Items() []interface{} {
	output := make([]interface{}, len(p.scanners))
	for i, _ := range output {
		output[i] = p.scanners[i].iface()
	}

	return output
}

func (p *planner) Finalize() {
	for _, s := range p.scanners {
		if s.column.Nullable {
			s.finalize()
		}
	}
}

type reflectScanner struct {
	parent reflector
	column *sourceMapping
	b      sql.NullBool
	f      sql.NullFloat64
	i      sql.NullInt64
	s      sql.NullString
	isnull bool
}

type reflector struct {
	item reflect.Value
}

func (rf *reflectScanner) iface() interface{} {
	if rf.column.Nullable {
		switch rf.column.Kind {
		case reflect.String:
			return &rf.s
		case reflect.Bool:
			return &rf.b
		case reflect.Float32, reflect.Float64:
			return &rf.f
		default:
			return &rf.i
		}
	} else {
		return rf.parent.item.Elem().Field(rf.column.Index).Addr().Interface()
	}
}

func (rf *reflectScanner) finalize() {
	switch rf.column.Kind {
	case reflect.String:
		if rf.s.Valid {
			rf.parent.item.Elem().Field(rf.column.Index).SetString(rf.s.String)
		} else {
			rf.isnull = true
		}
	case reflect.Bool:
		if rf.b.Valid {
			rf.parent.item.Elem().Field(rf.column.Index).SetBool(rf.b.Bool)
		} else {
			rf.isnull = true
		}
	case reflect.Float32, reflect.Float64:
		if rf.f.Valid {
			rf.parent.item.Elem().Field(rf.column.Index).SetFloat(rf.f.Float64)
		} else {
			rf.isnull = true
		}
	default:
		if rf.i.Valid {
			rf.parent.item.Elem().Field(rf.column.Index).SetInt(rf.i.Int64)
		} else {
			rf.isnull = true
		}
	}
}
