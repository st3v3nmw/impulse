#!/bin/bash

make build

N=128
declare -A arr
for ((i=0;i<$N;i++)); do
    arr["{$i}"]=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c32 ; echo '')
done

put () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        curl -X PUT -d ${arr[${key}]} http://127.0.0.1:3000/${key}
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    echo "PUT: $(printf %.2f $rps) requests per second"
}

get () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        fetched_value=$(curl -s -X GET http://127.0.0.1:3000/${key})
        if [[ $1 != true ]]; then
            # We know the value should exist in store
            if [ "${arr[${key}]}" != "$fetched_value" ]; then
                echo "${arr[${key}]} != $fetched_value"
            fi
        else
            # We know the value doesn't exist in store
            if [ "Not Found" != "$fetched_value" ]; then
                echo "Not Found != $fetched_value"
            fi
        fi
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    echo "GET: $(printf %.2f $rps) requests per second"
}

delete () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        curl -X DELETE http://127.0.0.1:3000/${key}
    done
    end=`date +%s%N`
    rps=$(bc -l <<< "1 / ((($end - $start) / $N / 1000000000))")
    echo "DELETE: $(printf %.2f $rps) requests per second"
}


declare -a engines=("LEVELDB" "IN_MEMORY_MAP")

for engine in "${engines[@]}"
do
    # setup
    printf "\n$engine\n"
    echo $(printf '%*s' ${#engine} "" | tr ' ' '=')
    ./impulse --engine=$engine --leveldb=level.db &>/dev/null &
    sleep 2

    # test
    put
    get false
    delete
    get true

    # teardown
    kill -9 $(ps -e | grep impulse | awk '{print $1}')
done
