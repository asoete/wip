bats_require_minimum_version 1.5.0
load "../lib.bash"

# ---------------------------------------------------------------------------- #

@test "wip/cli: ensure help has options" {

  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"-ctl.address address"* ]]

  [[ "${lines[@]}" = *"-db.dsn dsn"* ]]

  [[ "${lines[@]}" = *"-dump-config"* ]]

  [[ "${lines[@]}" = *"-pidfile path"* ]]

  [[ "${lines[@]}" = *"-web.address address"* ]]

}
