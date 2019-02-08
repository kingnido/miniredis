#/bin/bash

PORT=9000
n=99999

echo "Zadd $n + 1 members"
(
    curl -s -d "del asd" localhost:$PORT && echo

    for x in `seq 0 $n` ; do
        echo -n "$x: " && curl -s -d " zadd asd $x $x" localhost:$PORT && echo
    done

    curl -s -d "zcard asd" localhost:$PORT && echo
    curl -s -d "zrank asd 1000" localhost:$PORT && echo
    curl -s -d "zrank asd 83849" localhost:$PORT && echo

    curl -s -d "zcard asd" localhost:$PORT && echo
)

echo "Zadd $n + 1 members and range"
(
    curl -s -d "del asd" localhost:$PORT && echo

    for x in `seq 0 $n` ; do
        echo -n "$x: " && curl -s -d " zadd asd $x $x" localhost:$PORT && echo
    done

    curl -s -d "zcard asd" localhost:$PORT && echo
    curl -s -d "zrank asd 1000" localhost:$PORT && echo
    curl -s -d "zrank asd 83849" localhost:$PORT && echo
    curl -s -d "zrange asd 500 1000" localhost:$PORT && echo
)

echo "incr $n + 1 time"
(
    curl -s -d "del asd" localhost:$PORT && echo

    for x in `seq 0 $n` ; do
        curl -s -d "incr asd" localhost:$PORT && echo
    done

    curl -s -d "get asd" localhost:$PORT && echo
)
