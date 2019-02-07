#/bin/bash

n=99999

echo "Zadd $n + 1 members"
(
    for x in `seq 0 $n` ; do
        echo zadd asd $x $x
    done

    echo zcard asd
    echo zrank asd 1000
    echo zrank asd 83849
) | time go run . | tail
echo

echo "Zadd $n + 1 members and range"
(
    for x in `seq 0 $n` ; do
        echo zadd asd $x $x
    done

    echo zcard asd
    echo zrank asd 1000
    echo zrank asd 83849
    echo zrange asd 500 1000
) | time go run . | tail
echo

echo "incr $n + 1 time"
(
    for x in `seq 0 $n` ; do
        echo incr asd
    done

    echo get asd
) | time go run . | tail
echo
