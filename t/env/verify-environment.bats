bats_require_minimum_version 1.5.0

@test "verify_environment: ensure WIP_URL is defined" {

	test -n "${WIP_URL}"

	[[ "${WIP_URL}" =~ https?://[a-zA-Z0-9.]+:[0-9]+ ]]
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure WiP service is accepting requests" {

	run -0 curl -vis "${WIP_URL}"
}



# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure DB_FILE is defined" {

	test -n "${DB_FILE}"
}


