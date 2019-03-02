#!/bin/sh

FILE1=$1
FILE2=$2
echo 1
>&2 echo scorer error

[ -f "$FILE1" ] && [ -f "$FILE2" ] || exit 1
