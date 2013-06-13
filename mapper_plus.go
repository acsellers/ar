package ar

import "reflect"

type MapperPlus struct {
	model *model
}

func (mp *MapperPlus) Identity() *MapperPlus {
	return
}
func (mp *MapperPlus) Add() *MapperPlus {
	return &MapperPlus{val: mp.val, ops: append(mp.ops, "add")}
}
func (mp *MapperPlus) Sub() *MapperPlus {
	return &MapperPlus{val: mp.val, ops: append(mp.ops, "sub")}
}
func (mp *MapperPlus) Exec() string {
	output := fmt.Sprint(mp.val)
	for _, op := range mp.ops {
		output = fmt.Sprintf("%s(%s)", op, output)
	}

	return output
}
func (mp *MapperPlus) Init(nmp *MapperPlus) {
	mp.val = nmp.val
}
func InitMapperPlus(v interface{}) {
	rv := reflect.ValueOf(v).Elem()
	fv := rv.Field(0)
	mp := new(MapperPlus)
	mp.val = 246
	vmp := reflect.ValueOf(mp)

	if fv.Type().Kind() == reflect.Ptr {
		fv.Set(vmp)
	}
}

// end ar library

//start model code
type UserMapper struct {
	*MapperPlus
}

func (m *UserMapper) ASA() *UserMapper {
	return &UserMapper{m.Add().Sub().Add()}
}
func (m *UserMapper) SS() *UserMapper {
	return &UserMapper{m.Sub().Sub()}
}

var User = new(UserMapper)

func init() {
	InitMapperPlus(User)
}

// end model code

func main() {
	fmt.Println(User.ASA().SS().Exec())
}
