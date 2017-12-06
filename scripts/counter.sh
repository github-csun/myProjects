#!/bin/bash

ROUND=8

  TEST_RUN=1
  for TEST_RUN in $(seq 1 ${ROUND})
  do
    echo "current run is #${TEST_RUN} of ${ROUND}"
  done  