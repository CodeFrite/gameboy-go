#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

go test -v ./gameboy | while read line; do
  if echo "$line" | grep -q "FAIL"; then
    echo -e "${RED}${line}${NC}"
  elif echo "$line" | grep -q "PASS"; then
    echo -e "${GREEN}${line}${NC}"
  else
    echo "$line"
  fi
done
