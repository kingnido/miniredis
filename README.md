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

## Approach

The aproach to avoid much contention and keep the structures consistent, was to
split the mutexes in two domains: keys and values.

The keys domain is the global map. It's responsible by mapping keys to
reference to values.

The values domain is responsible accessing and modifiying the inner structure
of values.

The main pattern for operation is:

```
key_map.lock()
value = key_map[key] or key_map[key] = value
key_map.unlock()

value.lock()
value.do_something()
value.unlock()
```

This way it's possible to do long operations (like range in sets), without
blocking operations in other keys.

In this implementation, an operation on a value that is not being referenced
anymore is possible. It happens when the operation might take too long, and
some other operation that change references in the global map occur, like SET
or DEL. Since in the begging of the request the value was available, this
solution considers that the operation can proceed until the end successfully.

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

set value for key, returns OK

##### `set <key:string> <value:string> ex <delta:int>`

set value for key with expiration time in seconds. returns error if delta is
not int or is less than 1, otherwise OK

##### `incr <key:string>`

increments value for key. returns the new value as int. returns error if
value is not convertable to int

### Sorted Set

A sorted set does not keep duplicated elements, and can order the elements
based in a score (and comparing the values lexicografically in case of tie),
enabling range queries, and checking a value's position in the set.

It's implemented using an AVL (a self-balancing binary search tree), with order
statistics, to enable ranking an element, and selecting the k-th smallest
element in the set. Being n the number of elements int the set, this strategy
allow log(n) complexity in search, insert operations and rank operations. The
range operation has log(n)+m complexity, same as the original Redis, being m
the number of elements in the range.

It has a read-write mutex, allowing many readers at same time. On modifying
operations, only one thread can access the structure, to ensure consistency.

#### Commands

##### `zadd <key:string> <score:int> <member:string>`

adds a member with score to set for key create a set for key, if key does not
exist if member is updated, it will be placed in the proper position according
to the new score returns error if value for key is not a set. in success return
OK.

##### `zcard <key:string>`

returns the number of members for key. returns error if key does not exist
or value is not a set

##### `zrank <key:string> <member:string>`

returns position of the member in the set for key. returns error if key does
not exist or member does not exist or value is not a set

##### `zrange <key:string> <int:start> <int:stop>`

returns list of members in positions start to stop in the set for key
invalid indexes in the start - stop range are ignored. returns error if key
does not exist or value is not a set

### Other Commands

Some commands are independent of type.

#### Commands

##### `del <key:string>`

delete key. returns error if key does not exist

##### `dbsize`

returns the number of keys

##### `help`

returns all the commands usage.

## Build and Run

### Build
It's implemented and tested with Go 1.11. Having the Go environment properly set:

```
# on macos
cd <project_directory>
go build .
```

It will build a binary.

Functional tests were implemented for a large range of cases. The RedisCmd and
the server's request handler were tested with many ZADD requests in parallel,
to ensure consistency in the set structure. A reasonable coverage was
implemented for other structures as well. More can be done.

```
# on macos
go test .
```

It wont show much if everything is ok. For a more detailed log, run with -v flag.

### Run

This application can be run as a CLI, or HTTP server, or both.

`<application> [-cli] [-port=<port>]`

If cli is present, it will run a prompt where it's possible to run the commands
already mentioned.

If port is present, it will listen HTTP requests on the specified port.
Commands should be sent in the body of a POST request. Only one command per request is accepted.

To exit the application, press Ctrl-C.

### Docker

It's possible run it with Docker too.

```
cd <project_directory>
docker build -t miniredis .
docker run -p <desired_port>:8000 miniredis
```

### Requests

Requests can be done using the curl command:

```
curl -d "<cmd>" <server_url>
```

Successful requets will return 200. For for those that not return a bulk OK, the
result is JSON marshalled. It can be a integer, a string, or a list of string.

If some HTTP error is return, the body contains the error message.
