#!/bin/bash

make build

N=${N:-128}
declare -A arr
for ((i=0;i<$N;i++)); do
    arr["$i"]=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c32 ; echo '')
done

put () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        curl -X PUT -d ${arr[${key}]} http://127.0.0.1:3000/$key
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    echo "PUT: $(printf %.2f $rps) requests per second"
}

get () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        fetched_value=$(curl -s -X GET http://127.0.0.1:3000/$key)
        if [[ $1 == true ]]; then
            # We know the value should exist in store
            if [ "${arr[${key}]}" != "$fetched_value" ]; then
                echo "${arr[${key}]} != $fetched_value"
            fi
        else
            # We know the value doesn't exist in store
            if [ "" != "$fetched_value" ]; then
                echo "$fetched_value shouldn't exist in the store"
            fi
        fi
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    fetch_type=$([ $1 == true ] && echo "present" || echo "missing")
    echo "GET ($fetch_type): $(printf %.2f $rps) requests per second"
}

delete () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        curl -X DELETE http://127.0.0.1:3000/$key
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    echo "DELETE: $(printf %.2f $rps) requests per second"
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
    get true
    delete
    get false

    # Teardown
    kill -9 $(ps -e | grep impulse | awk '{print $1}')
done
