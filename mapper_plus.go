package db

type MapperPlus struct {
	source *source
	query  *Queryable
}

func (mp *MapperPlus) Identity() *MapperPlus {
	if mp.query != nil {
		return &MapperPlus{source: mp.source, query: mp.query}
	}

	return &MapperPlus{source: mp.source, query: &Queryable{source: mp.source}}
}

func (mp *MapperPlus) Where(fragment string, args ...interface{}) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Where(fragment, args...)}
}

func (mp *MapperPlus) EqualTo(column string, val interface{}) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.EqualTo(column, val)}
}

func (mp *MapperPlus) Between(column string, lower, upper interface{}) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Between(column, lower, upper)}
}

func (mp *MapperPlus) In(column string, vals []interface{}) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.In(column, vals)}
}

func (mp *MapperPlus) Limit(limit int) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Limit(limit)}
}

func (mp *MapperPlus) Offset(offset int) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Offset(offset)}
}

func (mp *MapperPlus) OrderBy(column, direction string) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.OrderBy(column, direction)}
}

func (mp *MapperPlus) Order(ordering string) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Order(ordering)}
}

func (mp *MapperPlus) Reorder(ordering string) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Reorder(ordering)}
}

// Find looks for the record with primary key equal to val
func (mp *MapperPlus) Find(val interface{}) *MapperPlus {
	return &MapperPlus{source: mp.source, query: mp.query.Find(val)}
}
