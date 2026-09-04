bats_require_minimum_version 1.5.0
load "../lib.bash"

set -ueo pipefail

# ---------------------------------------------------------------------------- #

@test "ensure BASE_URL is defined" {

	test -n "${BASE_URL}"

	[[ "${BASE_URL}" =~ https?://[a-zA-Z0-9.]+:[0-9]+ ]]

}

# ---------------------------------------------------------------------------- #

@test "ensure DB_FILE is defined" {

	test -n "${DB_FILE}"

	test -f "${DB_FILE}"

}


# ---------------------------------------------------------------------------- #

@test "POST /timers deadline=15:00" {

	wip_curl /timers -d deadline="15:00"

	assert_equal "Status: 200" "${stderr_lines[0]}" \
		"POST /timers returned invalid HTTP response code"

	log "TEST HTTP Header Location="
	test "${stderr_lines[1]}" = "Location:"

	# -- Validate info is found in database
	expected() {
		json_fmt \
			$(jo user_id=1 username=johnd sso_sub=0db52b4e-a61f-fc3f-58b4-d21c237151a1 ) \
			$(jo user_id=2 username=janed sso_sub=0db52b4e-a61f-fc3f-58b4-d21c237151a1 )
	}

	actual() {
		sqlite_json "SELECT * FROM users"
	}

	diff expected actual




	# run -0 --separate-stderr \
	# 	curl -vis --trace /dev/stdout "${BASE_URL}/timers" -d deadline="" \
	# 	-w "%{stderr}Status: %{http_code}\nLocation: %{redirect_url}"

	# 	# test "${stderr_lines[0]}" = "Status: 200"
	# 	test "${stderr_lines[1]}" = "Location:"


	# 	assert_sqlite_json \
	# 		"SELECT * FROM alarms ORDER BY alarm_id DESC LIMIT 1" \
	# 		"==" \
	# 		"$(jo alarm_id=1 owner_id=1 created_at=now deadline=15:00 cancelled_at="" )"
}


