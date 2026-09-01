bats_require_minimum_version 1.5.0
load "lib/diff.bash"

# ---------------------------------------------------------------------------- #

@test "wip/cli: ensure help has options" {

  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"-ctl.address address"* ]]

  [[ "${lines[@]}" = *"-db.dsn dsn"* ]]

  [[ "${lines[@]}" = *"-dump-config"* ]]

  [[ "${lines[@]}" = *"-pidfile path"* ]]

  [[ "${lines[@]}" = *"-web.address address"* ]]

}

# ---------------------------------------------------------------------------- #

@test "wip/cli" {
  skip
  readarray -t input <<<'
    apache:!!:17344::::::
    arnes:$6$GCtvSSz7$dhd4GZ:18747:0:99999:7:::
    boris:$6$GCtvSSz7$dhd4GZ:18512:0:99999:7::21123:
    camillel:$6$GCtvSSz7$dhd4GZ:20348:0:99999:7::20635:
    benjaminr:$6$GCtvSSz7$dhd4GZ:18852:0:99999:7::20454:
  '
  readarray passwd_file_content <<< '
    apache:x:48:48:Apache:/usr/share/httpd:/sbin/nologin
    arnes:x:1504:100:Arne.Soete,,,,:/home/arnes:/bin/bash
    benjaminr:x:2643:100:Benjamin.Rombaut,,,,:/home/benjaminr:/bin/bash
  '

  readarray -t expected_output <<< $(jq --sort-keys <<<'{
     "apache": "48",
     "arnes": "1504",
     "benjaminr": "2643"
   }')

  run -0 --separate-stderr bats_pipe printf "%s\n" "${input[@]}" \
    \| ./unix-account-expiry-parser.pl --debug-print-passwd-lookup --min-uid 500 --passwd-file <(printf "%s" "${passwd_file_content[@]}")  \
    \| jq --sort-keys -r

  # We provided a --passwd-file: no warning may be present
  [[ "${lines[@]}" != "[W] no --passwd-file specified, falling back to /etc/passwd. (required by --min-uid)"* ]]

  @diff "expected_output" "lines"
}

