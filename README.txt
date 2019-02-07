# Mini Redis

## Introduction

This simple project implements a small set of commands provided by Redis, with
minor changes in the responses.

## Commands

###	get <key:string>

returns value for key.
returns error if key does not exist or value is not string

###	set <key:string> <value:string>

set value for key

### set <key:string> <value:string> ex <delta:int>

set value for key with expiration time in seconds
returns error if delta is not int or is less than 1

# del <key:string>

delete key
returns error if key does not exist

### dbsize

returns the number of keys

### incr <key:string>

increments value for key
returns the new value as int
returns error if value is not convertable to int

### zadd <key:string> <score:int> <member:string>

adds a member with score to set for key
create a set for key, if key does not exist
if member is updated, it will be placed in the proper position according to the new score
returns error if value for key is not a set

### zcard <key:string>

returns the number of members for key
returns error if key does not exist or value is not a set

### zrank <key:string> <member:string>

returns position of the member in the set for key
returns error if key does not exist or member does not exist or value is not a set

### zrange <key:string> <int:start> <int:stop>

returns list of members in positions start to stop in the set for key
invalid indexes in the start - stop range are ignored
returns error if key does not exist or value is not a set`, nil
