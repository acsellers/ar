package arm

import (
	"github.com/acsellers/ar"
)

type Migration struct {
	Conn      *ar.Connection
	Log       io.Writer
	LogFormat int
}

func (mg *Migration) CreateTableForModel(model *ar.Model) error {

}
