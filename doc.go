/*
Package db is a way to interface with sql databases as more than
just simple data stores. At the moment, it is an ORM-like in an
alpha state.

Guiding Ideas

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

Mappers can also save structs to the database using the SaveAll function. You
can pass either a single instance or multiple instances to it, and it will use
the primary key value of the passed instance(s) to determine whether it needs
to update existing records or create new records.

  expiration := time.Now().Add(time.Duration(7*24*time.Hour))
  for _, post := range newPosts {
    if post.AboveAverage(newPosts) {
      post.FeaturedUntil = expiration
    }
  }

  Posts.SaveAll(newPosts)


MapperPlus

Mapper+'s are like mappers, but they are designed to be part of a user defined
struct, that the user then defines their own Scopes on, where each custom scope would be composed
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
  Connection.InitMapperPlus("User", &Users)

  goodUsers := Users.RecentlyActive()
  moreGoodUsers := goodUsers.Activated()
  // at this moment the results of moreGoodUsers == goodUsers
  // but they are different instances

When you need to duplicate a Mapper+ scope, you could use Identity, but that
will return you a Scope, not a MapperPlus. To assist in this situation, the
MapperPlus interface has a Dupe method that will return a MapperPlus for you
to use for this situation.

  type UserMapper struct {
    db.MapperPlus
  }
  var Users *UserMapper
  db.InitMapperPlus("User", &Users)

  // will error out
  &UserMapper{Users.Where(...).Identity()}

  // will compile correctly
  &UserMapper{Users.Where(...).Dupe()}

Scopes

Scopes are the cornerstone of db. Scopes are the common case for creating SQL
queries. Scopes come in different forms, depending on what you are wanting to
accomplish.

For instance, lets say I needed all students who have attended at least 1 class
and have a score of 90 on the 3 tests they've taken.

  myClass.Students.
    Join(myClass.Meetings).
    EqualTo("meeting_attendance.attended", true).
    Join(myClass.Quizzes).
    GroupBy(myClass.Students)
    Having("AVG(quiz_score.overall) >= ?", 90)

Or perhaps you would rather see the top 5 most popular posts on your blog
from the articles you released in the past month.

  Posts.Cond("created_at", GTE, time.Now().AddDate(0, -1, 0).Order("hits").Limit(5)

For detailed descriptions of each Scope function, please see the Queryable
interface documentation.

Database Mapping

There are multiple ways to retrieve data and map the data into struct instancess.

The Find function takes two parameters, the primary key of the record you want and
a pointer to the struct instance you want to map to. The Find function may start from
Mappers, Scopes and Mapper+'s. Note that the Find will still respect any conditions on
the Scope or Mapper+ if you are calling it from one of them.

  // find user by id
  var CurrentUser User
  err := Users.Find(session["user_id"].Int(), &User)

  // find user if the user is an admin
  var AdminUser User
  err := Users.Admins().Find(session["user_id"].Int(), &AdminUser)

The Retrieve function takes 1 parameter, the struct instance to map the first records
data into. If there are more than 1 records that would be returned from the current
Scope, Mapper, or Mapper+, then the first record will be the mapped record.

  // retrieve head writer for section
  var SectionHead User
  Users.Joins(SectionAssignments).EqualTo("section_assignments.section_id", section.Id).Retrieve(&SectionHead)

  // retrieve first created user
  var FirstUser User
  Users.OrderBy("created_at", "ASC").Retrieve(&User)

The RetrieveAll function takes 1 parameter, which is a pointer to an array of the struct
you want to map into.

  // get all the Users
  var Many []User
  Users.RetrieveAll(&Many)

  // get all the commentors for an article
  var Commentors []User
  Users.Joins(Comments).EqualTo("comments.article_id", CurrentArticle.Id).RetrieveAll(&Commentors)


Saving and Updating Values

You can save slices of new structs into the database using a Mapper using the SaveAll
call. You can also save single instances of structs as well using SaveAll, but you
will need to pass a pointer to the struct instance, so the mapper can update the
instance with the primary key assigned to that struct.

    // newPost is an unsaved post and newPost.Id is the zero value
    Posts.SaveAll(&post)
    // now post.Id will equal the primary key of the db record associated with it

    // otherPosts is an array of posts some of which are new, some of which need to be updated
    Posts.SaveAll(otherPosts)
    // posts that need to be saved, will be saved and their slice instance should be updated,
    // no matter whether the slice is of []Post or []*Post

You can also update columns in the database off of a Scope or a Mapper. There are three
functions, UpdateAttribute, UpdateAttributes, and UpdateSql that will to this for you.
UpdateAttribute takes a column name and a value, and will then update that column to
the value for all the database rows that would match the scope. UpdateAttributes takes a
map of column names to values so you may update more than 1 column at once. UpdateSql takes
a sql fragment and will allow you to write sql that uses sql functions instead of using dumb
values. UpdateSql will be less used when db.Formula objects are implemented. UpdateSql is
not yet implemented as well.

  Posts.EqualTo("late", true).UpdateAttribute("delete_on", time.Now().Add(10 * time.Minute))

The Count method allows you to retrieve a count of the rows that would be retrieved from
a Scope or Mapper.

  // Returns the number of Posts saved in the database
  Posts.Count()

  // Returns the number of Posts written by a specific user
  user.Posts.Count()

The Pluck method allows you to retrieve a selected column from a Scope, Mapper, etc. It
is then mapped into a simple array value that was passed as the second value.

  // get emails for users who haven't paid for last month
  var emails []string
  Users.Joins(Payments).Where("payments.month = ? AND payments.paid_on IS NULL", month).Pluck("email", &emails)

Possible Future Retrieval Methods

The CountOn method is a user controlled version of Count. If you would like to specify
a specific column, perhaps to do a DISTINCT count on, this is what you want.

  // get total number of distinct authors for a category of posts
  Posts.CountOn("DISTINCT category_id")

The PluckSeveral is similar to Pluck, but allows you to specify multiple parameters and arrays to
map results into. It uses a string array for the first parameters, then a variable amount
of pointers to the arrays for the data.

  // get emails and names for users who have paid for last month
  var emails, names []string
  Users.
    Joins(Payments).
    Where("payments.month = ? AND payments.paid_on IS NOT NULL", month).
    Pluck([]string{"name", "email"}, &names, &emails)

The Select function allows you to map specially selected columns and/or formulas into
purpose-written or anonymous structs. If a table has many columns, or you are returning
quite a bit of data, this can be a performance boost to use special structs instead of the
default mapper.

  // get weekly newsletter readers
  type weeklyReaders struct {
    Name, Email string
    Sections string
  }
  var readers []weeklyReaders

  columns := "users.name, users.email, GROUP_CONCAT(subscription_sections.name SEPARATOR '|')"
  Users.Joins(Subscriptions).Joins(SubscriptionSections).GroupBy("users.id").Select(columns, &readers)


There are also TableInformation and ScopeInformation interfaces. I would caution use
of the two interfaces at the moment, as they are intended to be improved heavily before
a stable release of db. A stable version of db will provide a comprehensive informational
interface for both Scopes and Mappers, but there are more pressing features than it at the
moment.

Mixin Functionality

If your use case involves significant use of the database, instead of using the
database as a simple persistence mechanism, you will enjoy the Mixin functionality
offered by db. When you add the db.Mixin struct as an embedded field to your sturcts,
you will have the ability to Save, Delete, and UpdateAttribute(s) from struct instances
directly instead of having to use the mapper objects.

Mixins need to be initialized
explicitly, this can be done by sending the instances individually, as a slice, or any
number of individual instances to the mapper for that sturct type's Initialize function.
You can also initialize individual instances by calling that instances Init function with
a pointer to the instance. This is only required if you are constructing your instances
manually and not using the Find/Retrieve/RetrieveAll Scope/Mapper functions. Find, Retrieve,
and RetrieveAll will all initialize the instances they retrieve if the instances have
Mixin instances. Instances do not need to be resident in the database for Initialization
to succeed. Instances also don't need to be initialized to be saved using the
Mapper.SaveAll function.

  // initalize an instance
  post := new(Post)
  post.Init(&post)
  post.Name = "Hello World"
  // Save the post to a new record in the corresponding database table
  post.Save()

  // initialize an instance that is mapped on multiple connections
  // this is only necessary when a struct is mapped on different connections
  // Init will return an error in situations when you must use this function
  post := new(Post)
  post.InitWithConn(pgConn, &post)
  post = new(Post)
  post.InitWithConn(myConn, &post)

  // initalize three instances at once
  post1, post2, post3 := new(Post), new(Post), new(Post)
  Posts.Initialize(&post1, &post2, &post3)

  var newPosts []Post
  ... // code that add instances to newPosts
  // initialize all instances in newPosts
  Posts.Initialize(newPosts)

Joining and Sub-struct Operations

While joining in db can be divided multiple ways, the simplest division may be the
division between automatic joins and manual joins. Manual Joins may be specified by
the user in the joins query and may add specifiers to the join call, or may be joining
on non-intuitive columns. Automatic joins are discovered during mapping by db and
can the be retrieved using the mapper or mixin, or using the Include Scope method.

Manual scopes are intended for use either in cases when you need a filtering that is
created from the existence of the join, or you need to select columns/formulas/etc.
from the query using the Select method of retrieval.

  // find posts that have not been commented on by a specific user
  Posts.Join("LEFT JOIN comments ON comments.user_id = 123456 AND comments.post_id = posts.id").
    EqualTo("comments.id", nil)

  // find the average score of good and bad comments for posts that have both good and bad comments
  var scoredPosts = []struct {
    Post,
    BadAverage, GoodAverage float64,
  }
  Posts.
    Join("INNER JOIN comments AS bad_comments ON bad_comments.post_id = posts.id AND bad_comments.score < 0").
    Join("INNER JOIN comments AS good_comments ON good_comments.post_id = posts.id AND good_comments.score > 0").
    Group("posts.id").
    Select(db.Selections{
      db.TableSelection("posts.*"),
      db.FormulaSelection("AVG(bad_comments.score)", "BadAverage"),
      db.FormulaSelection("AVG(good_comments.score)", "GoodAverage"),
      &scoredPosts,
    )


Automatic joins are declared as part of the struct, and then can be used in Join calls
by simply passing in a string or mapper corresponding to the joined struct. See an example
below. By default, joins are implemented as outer joins, but you can default specific joins
to be inner joins in the sql statment by setting the struct tag of db_join to be inner.
You can also use the alternative Join function, InnerJoin to have the join run as an inner
join. Finally, if you have set a db_join to default to inner, but want it to be a outer join
instead, you can use the OuterJoin function.

  type Post struct {
    Title, Body string

    AuthorId int
    Author *User `db_join:"inner"`
    Comments []*Comment
  }

  // find all posts with authors that are guests
  Posts.Join("Author").EqualTo("author.role", Author_GUEST)

  // retrieve the posts from a specific author, that have featured comments
  Posts.EqualTo("Author", theAuthor).Join(Comments).EqualTo("comments.featured", true)

Possible Future Join Enhancements

The FullJoin function allows you to retrieve records that don't have a match to the primary
mapped struct. You pass the normal Join paramerters, but add a pointer to an array of the
struct you are asking to be joined, which will be filled with the non-matching records when
the first Retrieve/RetrieveAll call is made

  // find all author's posts, also get one's that are missing an author
  var authors []User
  var orphaned []Post
  User.FullJoin(&orphaned, Posts).Cond("posts.created_at", GTE, oneYearAgo).RetrieveAll(&authors)

SQL Sundries

If you need to use functions to be evaluated by the sql server as part of conditions,
you can pass a formula created with the Func function. Func's can have their own
parameters, which you should specify using ?'s to denote where the values should appear.
Where scopings do not respect Func's at the moment, but they will in the future.

  // Find all appointments with a length of less than 5 minutes
  var tooShort []Appointment
  shortFunc := db.Func("TIMESTAMPADD(MINUTE, 5, appointments.begin_date_time)")
  Appointments.Cond("end_date_time", db.LT, shortFunc).RetrieveAll(&tooShort)

  // Where calls do not need Func instances to use functions
  Appointments.Where("end_date_time < TIMESTAMPADD(MINUTE, 5, appointments.begin_date_time)")

The Col function allows you to specify a column to be used as a parameter
in the same manner as a value or Func.

  // Find all appointments that have been updated
  Appointment.Cond("begin_date_time", db.LT, db.Col("end_date_time"))

Dialects

A Dialect creates a way for db to talk to a specific RDBMS. The current internal ones are
mysql and sqlite3, with postgres planned for the near future. You can replace existing
dialects or add your own dialects by writing a struct that corresponds to the Dialect
interface and then calling RegisterDialect with the name you want the dialect to be
accessible under and an instance of your dialect struct.

Logging

Before a public announcement of a db version, I need to implement the Logging facilities.
It won't be difficult, but it takes time. Time that I haven't invested yet.
*/
package db
