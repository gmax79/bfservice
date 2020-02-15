#!/bin/bash

if [[ -z $ABF_HOST ]]; then
    echo "error: enviroment variable ABF_HOST not declared"
    exit 1
fi

./tests $ABF_HOST || exit 1

echo "Integration tests via abfcli"
./abfcli use $ABF_HOST || exit 1

./abfcli pass "192.168.100.1"
./abfcli unpass "192.168.100.1"
./abfcli block "192.168.100.1"
./abfcli unblock "192.168.100.1"

echo "Check via abfcli check [11](login,password,host), 10-passed, 11-failed"
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep passed >/dev/null 2>&1 || exit 1
./abfcli check "login" "password" "127.0.0.1" | grep failed >/dev/null 2>&1 || exit 1
echo "Succefully tested"

./abfcli clear "login" "127.0.0.1" || exit 1

