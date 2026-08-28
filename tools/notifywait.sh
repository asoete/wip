#!/bin/sh

set -u

function main() {

   test $# -gt 0 || {
      printf "usage: $0 '<command>' [<path>...] \n" 1>&2
      exit 1
   }

   cmd="$1" ; shift

   declare -a paths

   test $# -gt 0 && paths=("$@") || paths=(".")

   while [ 1 ]; do

      inotify_output=$(inotifywait -r -q \
         -e create \
         -e modify \
         -e delete \
         --exclude '(.*sw[px]$|~$)' \
         --format '%e `%w%f`' \
         "${paths[@]}")

      printf ' *** %s -> EXEC `%s` \n' "$inotify_output" "$cmd"
      eval "${cmd}" | sed -e 's/^/[ntfyw] /'

      test -n "${ONESHOT:-}" && exit

   done
}

main "$@"
