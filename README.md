AR
===

##Current Plan
* Build in support for ActiveRecord's pluralization ideas (I already have
a library written, just need to include the library)
* Move to a ARConnection that you call once and then use in a threadsafe 
manner, instead of the get new QBS thing.
* Better documentation in the code (plus examples)
* Scopes?

##Exising Features from qbs

* Define table schema in struct type, create table if not exists.
* Detect table columns in database and alter table add new column automatically.
* Define selection clause in struct type, fields in the struct type become the columns to be selected.
* Define join query in struct type by add pointer fields which point to the parent table's struct.
* Do CRUD query by struct value.
* After a query, all the data you need will be filled into the struct value.
* Compose where clause by condition, which can easily handle complex precedence of "AND/OR" sub conditions.
* If Id value in the struct is provided, it will be added to the where clause.
* "Created" column will be set to current time when insert, "Updated" column will be set to current time when insert and update.
* Struct type can implement Validator interface to do validation before insert or update.
* Support MySQL, PosgreSQL and SQLite3.
* Support connection pool.

## Qbs Performance Claims

`Qbs.Find` is about 60% faster on mysql, 130% faster on postgreSQL than raw `Db.Query`, about 20% slower than raw `Stmt.Query`. (benchmarked on windows).
The reason why it is faster than `Db.Query` is because all prepared Statements are cached in map.

##Install

Don't install yet for production, as I'm making breaking changes, but it 
will likely support only go1.1.

## API Documentation

Will be on godoc when I've gotten some work into it.

##History
AR is a fork of qbs by [Ewan Chou](https://github.com/coocood)

##Contributors on qbs
[Erik Aigner](https://github.com/eaigner)
Qbs was originally a fork from [hood](https://github.com/eaigner/hood) by [Erik Aigner](https://github.com/eaigner), 
but I changed more than 80% of the code, then it ended up become a totally different ORM.

[NuVivo314](https://github.com/NuVivo314),  [Jason McVetta](https://github.com/jmcvetta)
