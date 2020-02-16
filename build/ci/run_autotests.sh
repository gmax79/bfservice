#!/bin/bash

./abf &
abf_pid=$!

export ABF_HOST="localhost:9000"
./tests_entrypoint.sh
status=$?

kill $abf_pid
exit $status
