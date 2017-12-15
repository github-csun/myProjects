#!/bin/bash

# script name
PROG_NAME=$(echo $0 | awk -F "/" '{print $NF}')
# echo "script name is ${PROG_NAME}"

# absolute path of current
CURRENT_DIR=$(pwd)
# echo "current absolute path is:     ${CURRENT_DIR}"

# relative path of command
RELATIVE_PATH=$(echo $0 | awk 'BEGIN{FS=OFS="/"} {$NF=""; NF--; print}')
# echo "script relative path is:      ${RELATIVE_PATH}"

# concat 
FULL_PATH="${CURRENT_DIR}/${RELATIVE_PATH}"
# echo "script absolute path is:      ${FULL_PATH}"