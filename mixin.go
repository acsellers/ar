package db

import (
	"fmt"
	"reflect"
)

/*

The Mixin struct can be embedded into your structs to allow for ActiveRecored (the pattern)
like operations. It does need to be Initialized, which you can do manually when you create
the instances with User{} or new(User), or it will be automatically done when you map it using Find,
Retrieve, or RetrieveAll. If you are using constructor functions for your structs, you should
add this to the function.

  func NewUser(p UserParams) *User {
    u := &User{}
    u.Init(u)
    u.SetFromParams(p)

    return u
  }

*/
type Mixin struct {
	instance interface{}
	nulls    []string
	model    *source
}

// Manually initialize a Mixin. Pass a pointer to the current instance in to initialize.
//
//  user.Init(&user)
func (m *Mixin) Init(instance interface{}) error {
	tn := fullNameFor(getType(instance))
	s := mappedStructs[tn]
	if s == nil {
		return fmt.Errorf("Could not locate a mapper for this struct")
	}
	if s.multiMapped {
		return fmt.Errorf("Must use InitWithConn for initializing this struct")
	}

	s.Initialize(instance)
	return nil
}

// Manually initialize a struct that is mapped on multiple connections.
//
//  user.InitWithConn(pgConn, &user)
func (m *Mixin) InitWithConn(conn *Connection, instance interface{}) error {
	tn := fullNameFor(getType(instance))
	s := conn.mappedStructs[tn]
	if s == nil {
		return fmt.Errorf("Could not locate a mapper for this struct")
	}

	s.Initialize(instance)

	return nil
}

// Save the instance to the database. If the primary key field isn't set, this will create
// the object, otherwise it will use the primary key field to update the fields. Mixin should
// log the fields so we don't have to update all fields, but that hasn't been written yet.
//
//  user.FirstName = "Bob"
//  user.LastName = "Zealot"
//  e := user.Save()
func (m *Mixin) Save() error {
	return m.model.SaveAll(m.instance)
}

func (m *Mixin) selfScope() Scope {
	id := m.model.extractID(reflect.ValueOf(m.instance).Elem())
	return m.model.EqualTo(m.model.ID.SqlColumn, id)
}

// Delete the database record associated with this instance
//
//  user.Delete()
func (m *Mixin) Delete() error {
	return m.selfScope().Delete()
}

// Update the database record record with the column name attr with the value passed
// Note: the instance that you are calling this on will not get the updated values
//
//  user.UpdateAttribute("visit_time", time.Now())
func (m *Mixin) UpdateAttribute(attr string, value interface{}) error {
	return m.selfScope().UpdateAttribute(attr, value)
}

// Update multiple columns of the database record with the passed values
// Note: the instance you are calling this on will not get the updated values
//
//  apptAttendance.UpdateAttributes(db.Attributes{
//    "attended": true,
//    "arrival": time.Now(),
//  })
func (m *Mixin) UpdateAttributes(values Attributes) error {
	return m.selfScope().UpdateAttributes(values)
}

// Return whether a column had the value NULL when retrieved from the database
// In this manner, you don't need to use sql.NullString, or *string values in your
// structs for fields that may be nullable in the database.
//
//  if userAttendance.IsNull("arrival") {
//    fmt.Println(userAttendance.User.Name, "did not attend")
//  }
func (m *Mixin) IsNull(column string) bool {
	for _, n := range m.nulls {
		if column == n {
			return true
		}
	}
	return false
}

// Sets whether a column should have a null value in the database.
// TODO: this only works for IsNull atm, should also be used for saving
func (m *Mixin) SetNull(column string) {
	m.nulls = append(m.nulls, column)
}
