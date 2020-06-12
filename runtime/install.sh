#!/bin/bash

CPARG=$1


res=""
if [ -f /install/${CPARG}.sh ]; then
  res=$(/install/${CPARG}.sh)
fi

if [ -z "$res" ]; then
  res="{}"
fi

echo ${res}

exit 0
