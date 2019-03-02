#!/bin/sh

FILE1=$1
FILE2=$2
echo model output
echo 1 > $FILE2
>&2 echo model error

[ -f "$FILE1" ] || exit 1
