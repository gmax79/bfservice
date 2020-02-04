#!/bin/bash

stop=0
if [[ -z $ABF_HOST ]]; then
    echo "error: enviroment variable ABF_HOST not declared"
    stop=1
fi
if [[ -z $ABF_LOGIN_RATE ]]; then
    echo "error: enviroment variable ABF_LOGIN_RATE not declared"
    stop=1
fi
if [[ -z $ABF_PASSWORD_RATE ]]; then
    echo "error: enviroment variable ABF_PASSWORD_RATE not declared"
    stop=1
fi
if [[ -z $ABF_IP_RATE ]]; then
    echo "error: enviroment variable ABF_IP_RATE not declared"
    stop=1
fi
if [[ -z $ABF_REDIS_HOST ]]; then
    echo "error: enviroment variable ABF_REDIS_HOST not declared"
    stop=1
fi

if [[ $stop == 1 ]]; then
    echo "abf stopped due to errors"
    exit 1
fi

if [[ -z $ABF_LOG_LEVEL ]]; then
    export ABF_LOG_LEVEL="info"
fi
if [[ -z $ABF_LOG_FILE ]]; then
    export ABF_LOG_FILE="/var/log/abf"
fi

logdir=$(dirname "$ABF_LOG_FILE")
mkdir -p "$logdir"

if [[ -z $ABF_LOG_ENCODING ]]; then
    export ABF_LOG_ENCODING="console"
fi

if [[ -z $ABF_REDIS_DB ]]; then
    export ABF_REDIS_DB="0"
fi

envsubst < abf_config.json.template > config.json
cat config.json

./abf || exit 1
