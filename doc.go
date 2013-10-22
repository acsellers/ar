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


*/
package db
