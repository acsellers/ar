package ar

type MapperPlus struct {
	model *model
	query *Queryable
}

func (mp *MapperPlus) Identity() *MapperPlus {
	if mp.query != nil {
		return &MapperPlus{model: mp.model, query: mp.query}
	}

	return &MapperPlus{model: mp.model, query: &Queryable{model: mp.model}}
}

func (mp *MapperPlus) Where(fragment string, args ...interface{}) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Where(fragment, args...)}
}

func (mp *MapperPlus) EqualTo(column string, val interface{}) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.EqualTo(column, val)}
}

func (mp *MapperPlus) Between(column string, lower, upper interface{}) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Between(column, lower, upper)}
}

func (mp *MapperPlus) In(column string, vals []interface{}) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.In(column, vals)}
}

func (mp *MapperPlus) Limit(limit int) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Limit(limit)}
}

func (mp *MapperPlus) Offset(offset int) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Offset(offset)}
}

func (mp *MapperPlus) OrderBy(column, direction string) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.OrderBy(column, direction)}
}

func (mp *MapperPlus) Order(ordering string) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Order(ordering)}
}

func (mp *MapperPlus) Reorder(ordering string) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Reorder(ordering)}
}

// Find looks for the record with primary key equal to val
func (mp *MapperPlus) Find(val interface{}) *MapperPlus {
	return &MapperPlus{model: mp.model, query: mp.query.Find(val)}
}
