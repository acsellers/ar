package db

import (
	"strings"
)

type join struct {
	Type     string
	Table    string
	Joined   *source
	Alias    string
	Matches  []string
	Args     []interface{}
	Compiled string
}

func (j *join) Fragment() string {
	if j.Compiled != "" {
		return j.Compiled
	} else {
		output := j.Type + " JOIN " + j.Table
		if j.Alias != "" {
			output += " AS " + j.Alias
		}
		return output + " ON " + strings.Join(j.Matches, " AND ")
	}
}
func (j *join) Values() []interface{} {
	return j.Args
}

func (j *join) String() string {
	return withVars(j.Fragment(), j.Values())
}

func newJoin(Type string, desc interface{}, on *queryable) ([]*join, []condition) {
	//j := &join{Type: Type}
	switch dv := desc.(type) {
	case *source:
		path := sourceVisitor(on.source, dv)
		if len(path) > 0 {
			return []*join{}, []condition{}
		} else {
			for _, qj := range on.joins {
				path = sourceVisitor(qj.Joined, dv)
				if len(path) > 0 {
					return []*join{}, []condition{}
				}
			}
		}
	case *mapperPlus:
		//path := sourceVisitor(on.source, dv.source)
	case *queryable:
	case string:
	}
	return []*join{}, []condition{}
}

func locateRelation(from *queryable, to *source) []*source {
	panic("NOT HERE")
}

type relationRoute struct {
	head *sourceMapping
	body []*sourceMapping
}

// it's breadth first search for relations, amazing
func sourceVisitor(f, t *source) []*sourceMapping {
	queue := []relationRoute{}
	for _, r := range f.relations {
		queue = append(queue, relationRoute{r, []*sourceMapping{r}})
	}
	visited := make(map[*source]bool)
	visited[f] = true

	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		if c.head.Relation == t && !c.head.Aliased() {
			return c.body
		}
		for _, r := range c.head.Relation.relations {
			if !visited[r.Relation] {
				visited[r.Relation] = true
				queue = append(queue, relationRoute{r, append(c.body, r)})
			}
		}
	}
	return []*sourceMapping{}
}

func aliasVisitor(f *source, t string) []*sourceMapping {
	queue := []relationRoute{}
	for _, r := range f.relations {
		queue = append(queue, relationRoute{r, []*sourceMapping{r}})
	}
	visited := make(map[*source]bool)
	visited[f] = true

	for len(queue) > 0 {
		c := queue[0]
		visited[c.head.Relation] = true
		queue = queue[1:]
		if c.head.structOptions.Name == t && c.head.Aliased() {
			return c.body
		}
		for _, r := range c.head.Relation.relations {
			if !visited[r.Relation] {
				queue = append(queue, relationRoute{r, append(c.body, r)})
			}
		}
	}
	return []*sourceMapping{}
}
