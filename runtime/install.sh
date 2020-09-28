#!/bin/bash

CPARG=$1


res=""
if [ -f /install/${CPARG}.sh ]; then
  res=$(/bin/bash /install/${CPARG}.sh)
fi

if [ -z "$res" ]; then
  res="{}"
fi

echo ${res}

exit 0
