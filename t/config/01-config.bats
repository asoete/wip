bats_require_minimum_version 1.5.0
load "../lib.bash"

# ---------------------------------------------------------------------------- #

@test "wip/config: --web.address address" {


  flag="-web.address"

  # Test --help
  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"${flag} address"* ]]

  # Test default value
  run -0 ./bin/wip --dump-config
  [[ "${lines[@]}" = *"-${flag} = 127.0.0.1:8080"* ]]

  run -0 ./bin/wip --dump-config -${flag} 0.0.0.0:8888
  [[ "${lines[@]}" = *"-${flag} = 0.0.0.0:8888"* ]]

  run -0 ./bin/wip --dump-config -${flag} 192.168.10.32:33997
  [[ "${lines[@]}" = *"-${flag} = 192.168.10.32:33997"* ]]

  run -0 ./bin/wip --dump-config -${flag} :8000
  [[ "${lines[@]}" = *"-${flag} = :8000"* ]]

  run -2 ./bin/wip --dump-config -${flag}

}


# ---------------------------------------------------------------------------- #

@test "wip/config: --ctl.address address" {

  flag="-ctl.address"

  # Test --help
  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"${flag} address"* ]]

  # Test default value
  run -0 ./bin/wip --dump-config
  [[ "${lines[@]}" = *"-${flag} = 127.0.0.1:8100"* ]]

  run -0 ./bin/wip --dump-config -${flag} 0.0.0.0:8412
  [[ "${lines[@]}" = *"-${flag} = 0.0.0.0:8412"* ]]

  run -0 ./bin/wip --dump-config -${flag} 192.168.10.32:64123
  [[ "${lines[@]}" = *"-${flag} = 192.168.10.32:64123"* ]]

  run -0 ./bin/wip --dump-config -${flag} :9100
  [[ "${lines[@]}" = *"-${flag} = :9100"* ]]

  run -2 ./bin/wip --dump-config -${flag}
}

# ---------------------------------------------------------------------------- #

@test "wip/config: --db.dsn dsn" {

  flag="-db.dsn"

  # Test --help
  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"-db.dsn dsn"* ]]

  # Test default value
  run -0 ./bin/wip --dump-config
  [[ "${lines[@]}" = *"-${flag} = sqlite::memory:"* ]]

  run -0 ./bin/wip --dump-config -${flag} "sqlite:/path/to/db.sqlite"
  [[ "${lines[@]}" = *"-${flag} = sqlite:/path/to/db.sqlite"* ]]

  run -2 ./bin/wip --dump-config -${flag}
}

# ---------------------------------------------------------------------------- #

@test "wip/config: --pidfile path" {

  flag="-pidfile"

  # Test --help
  run -0 ./bin/wip --help

  [[ "${lines[@]}" = *"-pidfile path"* ]]

  # Test default value (-> empty)
  run bats_pipe -0 ./bin/wip --dump-config \| grep -qE -- '--pidfile = \s*$ '

  run -0 ./bin/wip --dump-config -${flag} "/dev/shh/wip/wip.sqlite"
  [[ "${lines[@]}" = *"-${flag} = "/dev/shh/wip/wip.sqlite""* ]]

  run -2 ./bin/wip --dump-config -${flag}
}


# ---------------------------------------------------------------------------- #

@test "wip/config: X-WiP-Ctl-Token" {

  # Test default value
  run -0 ./bin/wip --dump-config

  [[ "${lines[@]}" = *"X-WiP-Ctl-Token = wip-ctl-token-please-change!"* ]]


  # Test value set by ENV var
  WIP_CTL_TOKEN="testing-the-ctl-token" \
     run -0 ./bin/wip --dump-config

  [[ "${lines[@]}" = *"X-WiP-Ctl-Token = testing-the-ctl-token"* ]]

}

