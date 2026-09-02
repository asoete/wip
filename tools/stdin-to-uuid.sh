#!/usr/bin/sh

###
### 1. read STDIN
### 2. Calculate md5sum
### 3. Format md5sum as UUID (eg.: fb2e186c-1544-4456-9712-ff4f294002fc)
###
### USAGE:
###   <command> | ./stdin-to-uuid.sh
###
### EXAMPLE
###   printf "johnd" | ./stdin-to-uuid.sh
###     # 0db52b4e-a61f-fc3f-58b4-d21c237151a1

set -euo pipefail

md5sum | awk '{ print $1}' | sed -E 's/(.{8})(.{4})(.{4})(.{4})/\1-\2-\3-\4-/'
