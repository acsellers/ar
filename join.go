package db

import (
	"strings"
)

type join struct {
	Type    string
	Table   string
	Alias   string
	Matches []string
}

func (j *join) String() string {
	output := j.Type + " JOIN " + j.Table
	if j.Alias != "" {
		output += " AS " + j.Alias
	}
	return output + " ON " + strings.Join(j.Matches, " AND ")
}
