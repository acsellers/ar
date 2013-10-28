DB - Sql Database Adapter
=========================

[![Build Status](https://travis-ci.org/acsellers/db.png)](https://travis-ci.org/acsellers/db)

There are plenty of mappers for Go, take your pick. DB has no desire to be just
a mapper, DB wants to be a very good mapper, as well as a something more.


DB Features Current + Planned
-----------------------------

* Map data into structs from sql tables
  * Not require a bunch of struct tags for each mapped column
  * Allow users to have string's instead of sql.NullStrings, etc.
  * For complex mapping, allow developers to customize mapping behavior 
* Map custom queries into adhoc structs or existing structs with a subset of attributes
* Track related structs using a mixin object to do recursive saving
* Save either via a Mapper.Save or instance.Save (via activated mixin)
* Initialize structs using a Mapper.Init function
* Retrieve related struct in the original .Retrive[All?] using Include
* Retrieve related structs later using a call from the mixin
* Multiple database Mysql, sqlite3, ...

