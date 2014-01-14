package db

/*
The Identity Scope will always return a copy of the current scope, whether it
is on a Mapper, Scope or Mapper+. Internally it is the canonical method to
duplicate a Scope.

  recent := Posts.Where(...).Joins(...).Order(...)
  recent2 := recent.Identity()
  // recent != recent2
  // recent.RetrieveAll(...) == recent2.RetrieveAll()

The Cond Scope is the matching shorthand scope. It allows you to run common SQL
matching conditions, but not have to write out the sql code like you would for
the Where condition. Cond supports the matchers: Equal, Not Equal, Less Than,
Less Than or Equal To, Greater Than, and Greater Than or Equal To.

  // find all posts published since recentTime
  recentScope := Posts.Cond("publish_date", db.GT, recentTime)

  // find all users created before longTimeAgo
  ancientScope := Users.Cond("created_at", db.LESS_OR_EQUAL, longTimeAgo)

  // find all appointments happening today
  todaysAppointments := Appointments.Cond("end_time", db.GTE, beginningOfDay).Cond("begin_time", db.LTE, endingOfDay)

The EqualTo scope is a short hand way of saying Cond(column, db.EQUAL_TO, value).

  // find all admins
  admins := Users.EqualTo("is_admin", true)

  // find nicer users
  niceUsers := Users.Joins(Bans).EqualTo("bans.id", nil)

  // find non-permabanned users
  meanUsers := Users.EqualTo("permabanned_on", nil)

The Between Scope is a way to specify a SQL BETWEEN without writing a Where fragment.

  // find last weeks posts
  lastWeek := Posts.Between(twoWeeksAgo, oneWeekAgo)

  // find the days appointments
  appts := Appointment.Between(beginningOfDay, endingOfDay)

The In Scope allows you to specify that a value should match 1 or more values in an array

  // find appointment conflicts
  conflicts := Appointments.Joins(Attendees).In("attendees.user_id", userIds)

  // find powerful users
  daPower := Users.In("user_type", []string{ADMIN, AUTHOR, BDFL, CTHULHU})

The Where Scope is a generalized way to write SQL Expressions. It can do simple matching
like EqualTo, Cond, Between or In written in SQL. It will also handle binding variables
within a SQL statement.

  // find awesome users
  Users.Where("kudos_count >= ?", coolnessQuotient)

  // find users with names like "mike"
  Users.Where("first_name LIKE ?", "mike")

  // find appointments that will be missed
  Attendees.Where(
    "user_ids IN (:students:) AND cancelled_on IS NULL AND begin_time BETWEEN :begin: AND :end:",
    map[string]interface{}{
      "students": userIds,
      "begin":    tripBeginning,
      "end":      tripEnding,
    },
  )

The Limit Scope allows you to specify the maximum number of records returned in a
RetrieveAll call.

  // find the 10 newest users
  Users.Limit(10).Order("created_at DESC")

The Offset specifies the number of records that will be skipped when returning records

  // find the users for the grid
  Users.Limit(25).Offset((pageNum - 1) * 25)

The ordering scopes are 3 different Scopes. The simplest is Order, which could be a sql fragment
or just a column name. If you do not pass a specific direction, then the direction will be
set to ascending. OrderBy allows you to specify the ordering as a string. Reorder will wipe
any previous orders and replace them with the passed ordering.

  // order by age
  Users.Order("birth_date")

  // order by creation date
  Users.OrderBy("created_at", "DESC")

  // order only by beginning time
  Appointments.Reorder("begin_time")

Joining and Including Related Tables

There's two different ways of dealing with joining in db, joining for the
use in Scopes and joining for the use of returning related data. When you
just want to use the data for scoping conditions or for Pluck/custom
selections, you can use the *Join functions (InnerJoin, LeftJoin, FullJoin,
RightJoin).

If you need to retrieve the joined records and have them
present in the data structures you pull back, you can use the *Include
functions. LeftInclude and InnerInclude you can specify normally, but
RightInclude and FullInclude are different. in that you are expected to
send along an array to put the records that can't be joined. Read up on how
Right and Full Outer Joins work if you don't understand why you have to send
in the array.

Specifying Joins or Includes without writing SQL is done in several ways. The
simplest is to just pass in the Mapper or MapperPlus for the table you wish
to join or include. That works fine for then the join is a simple unaliased
join from one mapper to another. Passing in the mapper will work whether the
relationship is 1-to-1 or 1-to-many.

But not every join is simple, perhaps you
have multiple joins of a mapper in a single struct, some or all being aliased,
in that case you can pass in a string that describes the join. For instance,
if an Appointments mapper had an Coordinator User as well as Users through
an Attendee Join Mapper, you could use the string "Coordinator" to indicate
that Users join to work as opposed to the Attendee Users that would happen
if you just passed in the Users Mapper.

Simple Join Examples

  // Join Comments to Users for a recentPost, then select some things
  Comments.EqualTo("post_id", recentPost.Id).InnerJoin(User).Select(...)

  // Join From Threads To Entries as well as the Original Post (First Entry)
  Threads.LeftJoin(Entries).InnerJoin("OP")...

Aliased Join Examples

  // Join all supervised departments, then if there are any recent employees add them
  Departments.InnerJoin("Supervisor").LeftJoin(Employees)

  // Join students and the creator of an enrollment (an instance of a college class in a specific semester)
  Enrollment.InnerJoin("Students").InnerJoin("Creator")

SQL Join Examples

  // Custom polymorphic join
  Meeting.JoinSql(`INNER JOIN calendared ON
    calendared.parent_id = meeting.id AND calendared.parent_type = 'Meeting'
  `)

*/
type Queryable interface {
	// Identity is the canonical way to duplicate a Scope, it doesn't do anything else
	Identity() Scope

	// Cond is a quick interface to the simple compare operations
	Cond(column string, condition COND, val interface{}) Scope
	// The EqualTo scope is a short hand way of saying Cond(column, db.EQUAL_TO, value).
	EqualTo(column string, val interface{}) Scope
	// The Between Scope is a way to specify a SQL BETWEEN without writing a Where fragment.
	Between(column string, lower, upper interface{}) Scope
	// The In Scope is a way to specify an SQL IN for variable sized arrays of items
	In(column string, items interface{}) Scope
	// The Where Scope is a generalized way to write SQL Expressions. It can do simple matching
	// like EqualTo, Cond, Between or In written in SQL. It will also handle binding variables
	// within a SQL statement.
	Where(fragment string, args ...interface{}) Scope

	// The Having SQL clause allows you to filter on aggregated
	// values from a GROUP BY. Since Having always is using SQL
	// functions, it should be simpler to just write the SQL
	// fragment directly instead of using a SqlFunc constructor.
	Having(fragment string, values ...interface{}) Scope

	// GroupBy allows you to group by a table or column, this is
	// necessary for aggregation functions.
	GroupBy(groupItem string) Scope

	// Limit sets the number of results to return from the full query
	Limit(limit int) Scope
	// Offset sets the number of results to skip over before returning results
	Offset(offset int) Scope

	// Order sets an ordering column for the query, direction is ASC unless specified
	Order(ordering string) Scope
	// Specify both an ordering and direction as separate parameters
	OrderBy(column, direction string) Scope
	// Drop all previous order declarations and only order by the parameter passed
	Reorder(ordering string) Scope

	// Search for a record with the primary key of id, then place the result in the val pointer
	Find(id, val interface{}) error
	// Return the first result from the scope and place it into the val pointer
	Retrieve(val interface{}) error
	// Return all results from the Scope and put the in the array pointed at by the dest parameter
	RetrieveAll(dest interface{}) error
	// Return the count results, uses the primary key of the originating mapper to count on
	// Not distinct, need to add a CountSql function
	Count() (int64, error)
	// Retrieve a single column using joins, limits, conditions from the Scope and place
	// the results into the array pointed at by values
	Pluck(column, values interface{}) error

	// Run a DELETE FROM query using the conditions from the Scope
	Delete() error

	/*
	   WARNING: None of the Update* functions update timestamps
	*/
	// Update a single column using a UPDATE query
	UpdateAttribute(column string, val interface{}) error
	// Update multiple columns at once, using an
	UpdateAttributes(values Attributes) error
	// Update however you want off of a scope, go wild with SQL, you can pass in some values to be
	// substituted if you need them
	UpdateSql(sql string, vals ...interface{}) error

	// LeftJoin is a Left Outer Join to the Mapper or Scoped Mapper String
	LeftJoin(joins ...interface{}) Scope
	// InnerJoin will require both sides of the relationship to exist
	InnerJoin(joins ...interface{}) Scope
	// Full Join will make a Full Outer Join query using the passed arguments
	FullJoin(joins ...interface{}) Scope
	// Right Join will do a right join, if you need it
	RightJoin(joins ...interface{}) Scope
	// JoinSql will allow you to write straight SQL for the JOIN
	JoinSql(sql string, args ...interface{}) Scope

	// LeftInclude
	LeftInclude(include ...interface{}) Scope
	// InnerInclude
	InnerInclude(include ...interface{}) Scope
	// FullInclude
	FullInclude(include interface{}, nullRecords interface{}) Scope
	// RightInclude
	RightInclude(include interface{}, nullRecords interface{}) Scope
	// IncludeSql, magic come true
	IncludeSql(il IncludeList, query string, args ...interface{}) Scope
}

type IncludeList []interface{}

// Information used in query generation and Mapper interrogation. Much more information will
// be exposed at a later time
type TableInformation interface {
	TableName() string
	PrimaryKeyColumn() string
}

/*
These functions are used by Dialects to construct the queries. As Dialects
could be written without explicit support from db, these function (though we'll
get around to adding more) need to be publicly accessible.
*/
type ScopeInformation interface {
	SelectorSql() string
	ConditionSql() (string, []interface{})
	JoinsSql() string
	EndingSql() (string, []interface{})
}

/*
Mappers are the basic mechanism to turn database rows into struct instances.
You can Scope queries off of them, you can get basic SQL information about
their database table, and you can act on multiple instances of the mapped struct.

Initialize allows you to Initialize Mixins on struct instances ranging from just one
instance, to several instances, to many instances in an array. SaveAll is a either
a shortcut to calling Save on structs that have Mixin's or the main way to save
records when the mapped struct doesn't have a Mixin.
*/
type Mapper interface {
	Queryable
	TableInformation
	Initialize(val ...interface{}) error
	SaveAll(val interface{}) error
}

// A MapperPlus is both a Scope-like interface, but also the Mapper for a struct.
// It has to be that in order to be able to allow you to specify custom Scope-like
// functions. I'll add examples soon.
type MapperPlus interface {
	Mapper
	ScopeInformation
	Dupe() MapperPlus
}

// A Scope is an auto-duplicating Query construction mechanism
type Scope interface {
	Queryable
	TableInformation
	ScopeInformation
}

// SqlBit is a common interface shared by many internal SQL fragments
// like joins, SqlFunc's, etc.
type SqlBit interface {
	Fragment() string
	Values() []interface{}
	String() string
}

type mixedin interface {
	Init(interface{}) error
	InitWithConn(interface{}) error
	Save() error
	Delete() error
	UpdateAttribute(string, interface{}) error
	UpdateAttributes(Attributes) error
	IsNull(string) bool
	SetNull(string)
}
