bats_require_minimum_version 1.5.0

@test "verify_environment: ensure BASE_URL is defined" {

	test -n "${BASE_URL}"

	[[ "${BASE_URL}" =~ https?://[a-zA-Z0-9.]+:[0-9]+ ]]
}

# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure WiP service is accepting requests" {

	run -0 curl -vis "${BASE_URL}"
}



# ---------------------------------------------------------------------------- #

@test "verify_environment: ensure DB_FILE is defined" {

	test -n "${DB_FILE}"
}


