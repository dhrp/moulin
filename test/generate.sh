#!/bin/bash

# Generate test data
# Usage: ./generate.sh  <number of lines per file>

# echo an an integer and a sleep of one second into a file


for i in $(seq 1 $1)
do
    echo "sh -c \"echo task $i && sleep 1\""
done