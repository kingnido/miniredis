# MiniRedis

## Introduction

This simple project implements a small set of commands provided by Redis, with
minor changes in the responses.

Redis is an open source (BSD licensed), in-memory data structure store, used as
a database, cache and message broker. It supports data structures such as
strings, hashes, lists, sets, sorted sets with range queries, bitmaps,
hyperloglogs, geospatial indexes with radius queries and streams. 

In this project suport for string and sorted sets with range queries were
implemented.

## Data Structures

### Global Map

MiniRedis is a key-value structure, thread-safe, that can hold strings and
sorted sets. It's basically a map with a read-write mutex. This mutex allow
many readers access the data without blocking each other. If the write lock
is requested, it will wait for the readers to finish, block any new reader
request, run it's task, and unlock the write mutex, giving access to
readers once again.

### Strings

It's a simple string. The only operation that can change it's value is
INCR, if the string is parseable to an interger. This operation does not
change the value reference. Due to INCR operation, it has its own mutex.

New values can be set using SET. It replaces the old reference for a new
one in the global map.

Strings can be set with expiration time. If a new string is set for that
key, the expiration time is discarded. While a string does not expire, it
can be INCRed, since the reference string is the same.

#### Commands

##### `get <key:string>`

returns value for key. returns error if key does not exist or value is not
string

##### `set <key:string> <value:string>`

set value for key

##### `set <key:string> <value:string> ex <delta:int>`

set value for key with expiration time in seconds. returns error if delta is
not int or is less than 1

##### `incr <key:string>`

increments value for key. returns the new value as int returns error if
value is not convertable to int







##### `del <key:string>`

    delete key returns error if key does not exist

### `dbsize`

    returns the number of keys

### `zadd <key:string> <score:int> <member:string>`

    adds a member with score to set for key create a set for key, if key does
    not exist if member is updated, it will be placed in the proper position
    according to the new score returns error if value for key is not a set

### `zcard <key:string>`

    returns the number of members for key returns error if key does not exist
    or value is not a set

### `zrank <key:string> <member:string>`

    returns position of the member in the set for key returns error if key does
    not exist or member does not exist or value is not a set

### `zrange <key:string> <int:start> <int:stop>`

    returns list of members in positions start to stop in the set for key
    invalid indexes in the start - stop range are ignored returns error if key
    does not exist or value is not a set`, nil
