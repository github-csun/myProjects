#!/bin/bash

# Function to print timestamp in each line of log
function addTimeStamp() {
    while IFS= read -r line; do
      tempStamp=$(date "+%Y-%m-%d %H:%M:%S")
        echo "${tempStamp}  $line"
    done
}


# command | addTimeStamp
# example usage
for count in $(seq 1 10)
do
  echo "line ${count}" | addTimeStamp
  sleep 1
done

