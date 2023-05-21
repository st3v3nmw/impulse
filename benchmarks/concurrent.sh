#!/bin/bash

make build

T=${T:-10}

put () {
    printf "\nPUT\n"
    go run github.com/tsliwowicz/go-wrk@latest -c 100 -d $T -body "bar" -M PUT http://127.0.0.1:3000/foo
}

get () {
    printf "\nGET\n"
    go run github.com/tsliwowicz/go-wrk@latest -c 100 -d $T http://127.0.0.1:3000/foo
}

delete () {
    printf "\nDELETE\n"
    go run github.com/tsliwowicz/go-wrk@latest -c 100 -d $T -M DELETE http://127.0.0.1:3000/foo
}

declare -a engines=("HASH_MAP" "LEVELDB")

for engine in "${engines[@]}"
do
    # Setup
    printf "\n$engine\n"
    echo $(printf '%*s' ${#engine} "" | tr ' ' '=')
    ./impulse --engine=$engine --leveldb=level.db &>/dev/null &
    SERVER_PID=
    sleep 2

    # Run Tests
    put
    get
    delete

    # Teardown
    kill -9 $(ps -e | grep impulse | awk '{print $1}')
done
