bats_require_minimum_version 1.5.0
load "../lib.bash"

set -ueo pipefail


# ---------------------------------------------------------------------------- #

@test "ensure TABLE users is seeded" {

	skip

	run -0 libexec/db.migrate.sh t/db/seeds/insert-users.sql

	expected() {
		json_fmt \
			$(jo user_id=1 username=johnd sso_sub=0db52b4e-a61f-fc3f-58b4-d21c237151a1 ) \
			$(jo user_id=2 username=janed sso_sub=d2f9f861-496c-4639-6a1c-666a5a781406 )
	}

	actual() {
		sqlite_json "SELECT * FROM users WHERE user_id IN (1,2)"
	}

	diff expected actual
}

# ---------------------------------------------------------------------------- #

@test "ensure TABLE ntfy_channels is seeded" {

	run -0 libexec/db.migrate.sh t/db/seeds/insert-ntfy_channels.sql

	expected() {
		json_fmt $(jo channel_id=1 url=https://ntfy.sh/WiP-testing-channel_OfTD36Z4uARQeJsQB76pNQMoloPv )
	}

	actual() {
		sqlite_json "SELECT * FROM ntfy_channels WHERE channel_id = 1"
	}

	diff expected actual
}


