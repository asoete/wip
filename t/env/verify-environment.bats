bats_require_minimum_version 1.5.0
load "../lib.bash"

@test "verify_environment: ensure WIP_URL is defined" {

	test -n "${WIP_URL}"

	[[ "${WIP_URL}" =~ https?://[a-zA-Z0-9.]+:[0-9]+ ]]
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure WiP service is accepting requests" {

	run -0 curl -vis "${WIP_URL}"
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure wip_curl() is functioning" {

	run --separate-stderr -0 wip_curl /some-fake-page +ecode

	[[ ${stderr_lines[-1]} = "Status: 404" ]]
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure CTL_URL is defined" {

	test -n "${CTL_URL}"
	test -n "${CTL_AUTH_HEADER}"

	[[ "${CTL_URL}" =~ https?://[a-zA-Z0-9.]+:[0-9]+ ]]
	printf "%s\n" "${CTL_AUTH_HEADER}" | tee /dev/stderr | grep -E '^X-' -q
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure CTL service is accepting requests" {

	run -0 curl -vis "${CTL_URL}"
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure ctl_curl is functioning" {

	run --separate-stderr -0 ctl_curl /ping +ecode

	[[ ${lines[-1]} = "pong" ]]
	[[ ${stderr_lines[-1]} = "Status: 200" ]]
}


# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure DB_FILE is defined" {

	test -n "${DB_FILE}"
}


