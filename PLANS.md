Overall Plan

The idea is to get features in early, then fill out implementations based
on feedback and experience from writing examples and revellers. 

0.2 will hopefully be a fast follower release, making it out in 2 weeks or 
less after 0.1 drops. 0.1 is going to be announced at 0.1.pre, so I'll take
comments for a bit and fix the worst 0.1 mistakes then release 0.1.0. In that
time, I've got to get unitemplate to 0.1 as well as helpers. Hopefully I can
work on the migrator design while doing the rote work in helpers. So my schedule
will resemble the following chronologically:

* release 0.1.pre [10/11]
* take comments, fix issues, add tests (especially for MapperPlus) [11/11]
* restart work on unitemplate [11/11]
* release 0.1.0 [15/11]
* finish django part of unitemplate [13/11]
* restart work on helpers [14/11]
* start simpler parts of db 0.2.0 [16/11]
* ready unitemplate 0.1.0 [19/11]
* start revellers with working versions of things [19/11]
* release revellers 0.1, db 0.2.0, helpers 0.1.0 [??/11]

Version Features

0.1 - cheap hotel

* Mapping structs to tables
* Parameterized everywhere and bind variables for Where
* Scopes for most of sql syntax
* sqlite, mysql, postgres testing
* query cache
* documentation for major parts with at least 1 small example each
* Simple Mixin Operations
* Scoping from struct to struct
* Retrieving related structs using Mixin or Include* functions

0.2 - the sprawl

* Mapping sub-structs that are mapped to the same table
* Saving embedded/included structs
* Migration framework
* EphemeralConnections/Transaction
* Scope to Subquery
* Subquery helper
* Scope combining
* Revellers repository :)
* Default Struct Orderings
* Marshalling either by encoding or Marshal/Unmarshal interface

0.3 - zion

* Scope to Materialized View
* Slow Query Report
* JSON/XML Substructs + hstore
* Virtual Tabling
* MariaDB CI
* New Builtin Dialects?

0.4 - freeside

* SuperScoping and Aqua Vitae (I'll fill this out more later)
* More Things

0.5 - straylight

* Prepare for 1.0
* Features
