#!/bin/bash

N=100
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
    echo "PUT took $(( ($end - $start) / 1000 / $N )) microseconds"
}

get () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        fetched_value=$(curl -s -X GET http://127.0.0.1:3000/${key})
        if [[ $1 != true ]]; then
            if [ "${arr[${key}]}" != "$fetched_value" ]; then
                echo "${arr[${key}]} != $fetched_value"
            fi
        else
            if [ "Not Found" != "$fetched_value" ]; then
                echo "Not Found != $fetched_value"
            fi
        fi
    done
    end=`date +%s%N`
    echo "GET took $(( ($end - $start) / 1000 / $N )) microseconds"
}

delete () {
    start=`date +%s%N`
    for key in ${!arr[@]}; do
        curl -X DELETE http://127.0.0.1:3000/${key}
    done
    end=`date +%s%N`
    echo "DELETE took $(( ($end - $start) / 1000 / $N )) microseconds"
}


put
get false
get false
delete
get true
