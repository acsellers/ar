/*
Package db is a way to interface with sql databases as more than
just simple data stores. At the moment, it is an ORM-like in an
alpha state.

Guiding Ideas (Some unimplemented)

A Connection is both a sql.DB connection to an active database and
a mechanism for translating that database connection into structs.

Structs can be configured by the Connection.Config mechanism or
struct tags for special cases.

Mappers turn rows from sql tables into structs. Mappers can be used to
save structs to databases or retrieve subsets of the table using
scopes. Scopes are implemented on Mappers, which return a Queryable
that you can then chain off of to further filter the results that
would be returned with RetrieveAll.

Scopes are useful for more than just filtering results, you can also
act on them. You will be able to pull out specific attributes into
arrays of simple types. For example you should be able to run the
following code to get a list of users who meet certain conditions
and pull back their email addresses. Either way would work.

Option 1
  User.Where("welcomed_at", nil).
    Where("created_at >=", time.Now().Add(time.Duration(time.Minute * -15)).
    Pluck("email").([]string)

Option 2
  var emails []string
  User.Unwelcomed().AddedSince(time.Minute * -15).Pluck(email, &emails)

  type UserMapper *MapperPlus
  func (um UserMapper) Unwelcomed() *UserMapper {
    return &UserMapper{um.Where("welcomed_at", nil)}
  }
  func (um UserMapper) AddedSince(d time.Duration) *UserMapper {
    return &UserMapper{um.Where("created_at >= ", time.Now().Add(d))}
  }

But wait there's more. You can also run delete or update statements
from scopes. You can update via a column and value, a map of columns
to values, or straight sql. You can also run a sql delete from a scope.

Scopes are used for more than just that. Has Many relations can have
both an array of related structs, and a scope that are filled in by
the retrieval process. Calling user.Posts.Save(post) would set the
user relation on the post before saving it. user.Posts.Where(...)
would start with a queryable scoped to the user's posts before applying
the rest of the scope.

Connections

Start building an environment by calling NewConnection with you dialect
name ("mysql", "sqlite", "postgres") the name of the database to
use, and the connector string for the database adapter (document this).

There is code to manage Transactions using Connections, this isn't well
supported and should be removed if it is not well supported or
semi-supported until EphemeralConnections come around.

Mappers

Mappers are created by running CreateMapper on a Connection. You will receive
back a Mapper object that you can then use to retrieve database rows, save
structs into the database and create scopes of the struct.

  Posts := Connection.MustCreateMapper("Post", &Post{})

  // Retrieve posts created from ~14 days ago to present
  timeCutoff := time.Now().Add(time.Duration(-14*24*time.Hour))

  var recentPosts []Post
  Posts.Cond("created_at", db.GT, timeCutoff).RetrieveAll(&recentPosts)

  var myRecentPosts []Post
  Posts.EqualTo("author_id", Current.User.Id).Order("hits_counter").RetrieveAll(&myRecentPosts)

Mapper+

Mapper+'s are like mappers, but they are designed to be part of a user defined
struct, that the user then defines their own Scopes, which would be composed
of 1 or more regular scopes. The main difference between a Mapper and a Mapper+
is that Mappers return a Scope when you call a Queryable method on it, while
Mapper+ will return a new Mapper+ for the first scope, and then all further
scopes will be applied to the Mapper+ object. The reason for that is so that
user structs do not have to be continually casted to or created, the Scopes simply
add to the current Scope.

  type UserMapper struct {
    MapperPlus
  }
  func (um *UserMapper) Activated() *UserMapper {
    return &UserMapper{um.Cond("activated_at", db.NE, nil)}
  }
  func (um *UserMapper) RecentlyActive() *UserMapper {
    timeCutoff := time.Now().Add(time.Duration(-30*24*time.Hour))
    return &UserMapper{um.Cond("last_login", db.GTE, timeCutoff)}
  }

  var Users UserMapper
  Connection.CreateMapperPlus("User", &Users)

  goodUsers := Users.RecentlyActive()
  moreGoodUsers := goodUsers.Activated()
  // at this moment the results of moreGoodUsers == goodUsers
  // but they are different instances

Scopes

Scopes are the cornerstone of db. Scopes are the common case for creating SQL
queries. Scopes come in different forms, depending on what you are wanting to
accomplish.

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
    "user_ids IN (:students) AND cancelled_on IS NULL AND begin_time BETWEEN :begin AND :end",
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
*/
package db
