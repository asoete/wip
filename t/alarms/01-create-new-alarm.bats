bats_require_minimum_version 1.5.0
load "../lib.bash"

set -ueo pipefail

# ---------------------------------------------------------------------------- #

@test "POST /alarm deadline=17:00" {

# Directly submit to WiP (and not the proxy) so we can set the User data
	description="bats[$$]: POST /alarm deadline=15:00"
	deadline=$(date -d '17:00' '+%Y-%m-%d %H:%M:%S')
	run -0 --separate-stderr wip_curl /alarm +ecode \
		-H "REMOTE_USER: janed" \
		-H "SSO_SUB: d2f9f861-496c-4639-6a1c-666a5a781406" \
		-d deadline="$deadline" \
		-d description="$description"

	assert "Status: 200" = "${stderr_lines[-1]}" \
		"POST /alarm returned invalid HTTP response code"

	# ---

	expected() {
		json_fmt "$(jo deadline="$deadline" description="$description" )"
	}

	actual() {
		sqlite_json "SELECT DATETIME(deadline, 'localtime') AS deadline,description FROM alarms WHERE description = '$description'"
	}

	diff expected actual

	# ---

	query() {
		sqlite "$(printf "SELECT %s FROM alarms WHERE description = '$description' LIMIT 1" "$@")"
	}

	assert "$(query alarm_id)" != ""
	assert "$(query alarm_id)" != "0"
	echo $(query alarm_id) | grep -P '^\d+$'

	assert "$(query user)" != ""
	assert "$(query user)" = "janed"
	
	assert "\
		$(query "strftime('%s', created_at)")" '-gt' "$(date +%s -d 'now - 2 seconds')" \
		"'created_at' timestamp is outside expected range (now - 3 secs)"
	assert \
		"$(query "strftime('%s', created_at)")" '-le' "$(date +%s -d 'now')" \
		"'created_at' timestamp is outside expected range (timestamp is in the future??)"

	assert "$(query "DATETIME(deadline, 'localtime')")" = "$deadline"
	assert "$(query cancelled_at)" = ""
	assert "$(query description)" = "$description"
}
