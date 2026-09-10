bats_require_minimum_version 1.5.0
load "../lib.bash"

set -ueo pipefail

# ---------------------------------------------------------------------------- #

@test "validate current userinfo: johnd" {

	run -0 --separate-stderr \
		ctl_curl /userinfo +ecode \
		-H "REMOTE_USER: johnd" \
		-H "SSO_SUB: 30541499-fd52-24c1-17eb-e8377fdfc86c" \

	# Ensure HTTP returned OK
	[[ ${stderr_lines[-1]} == "Status: 200" ]]

	# Validate response
	expected() {
		jo \
			username=johnd \
			sso_sub=30541499-fd52-24c1-17eb-e8377fdfc86c \
		| jq . --sort-keys
	}

	actual() {
		printf "%s\n" "${lines[*]}" | jq --sort-keys
	}

	diff expected actual
}

# ---------------------------------------------------------------------------- #

@test "validate current userinfo: janed" {

	run -0 --separate-stderr \
		ctl_curl /userinfo +ecode \
		-H "REMOTE_USER: janed" \
		-H "SSO_SUB: e36dfb45-2485-950b-5a13-b88d5d1776cf" \

	# Ensure HTTP returned OK
	[[ ${stderr_lines[-1]} == "Status: 200" ]]

	# Validate response
	expected() {
		jo \
			username=janed \
			sso_sub=e36dfb45-2485-950b-5a13-b88d5d1776cf \
		| jq . --sort-keys
	}

	actual() {
		printf "%s\n" "${lines[*]}" | jq --sort-keys
	}

	diff expected actual
}


