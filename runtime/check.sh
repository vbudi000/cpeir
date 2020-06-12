#!/bin/bash

CPARG=$1

res=""
if [ -f /check/${CPARG}.sh ]; then
  res=$(/check/${CPARG}.sh)
fi

if [ -z "$res" ]; then
  res="{}"
fi

echo ${res}

exit 0
