#!/bin/bash

if [[ -z $ABF_HOST ]]; then
    echo "error: enviroment variable ABF_HOST not declared"
    stop=1
fi

./tests || exit 1

echo "Integration tests via abfcli"
./abfcli use $ABF_HOST || exit 1
