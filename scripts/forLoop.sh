#!/bin/bash

function counter() {
  MAX=8
  for COUNTER in $(seq 1 ${MAX})
  do
    echo "#${COUNTER} of ${MAX}"
  done
}

function fileInDir() {
  FOLDER=~/*
  for FILENAME in ${FOLDER}; do
    echo "${FILENAME} is found"
  done
}