bats_require_minimum_version 1.5.0
load "../lib.bash"

set -ueo pipefail

# ---------------------------------------------------------------------------- #

@test "Cancel alarm" {

	deadline=$(date -d 'now + 30 seconds' '+%Y-%m-%d %H:%M:%S')
	deadline_epoch=$(date -d 'now + 30 seconds' '+%s')
	description="bats[$$]: cancel alarm before $deadline"

	# ----

	# Directly submit to WiP (and not the proxy) so we can set the User data
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

	assert "$(query "DATETIME(deadline, 'localtime')")" = "$deadline"
	assert "$(query cancelled_at)" = ""
	assert "$(query description)" = "$description"

	# ----

	alarm_id="$(query alarm_id)" 

	assert "$(ctl_curl /alarm/${alarm_id}/dispatch-has-alarm)" = "true" \
		"assert the alarm was stored in the dispatcher"

	assert "$(ctl_curl /alarm/${alarm_id}/dispatch-has-timer)" = "true" \
		"assert the alarm has a timer configured in the dispatcher"

	run -0 --separate-stderr wip_curl /alarm/${alarm_id}/cancel +ecode \
		-X POST \
		-H "REMOTE_USER: janed" \
		-H "SSO_SUB: d2f9f861-496c-4639-6a1c-666a5a781406"

	assert "\
		$(query "strftime('%s', cancelled_at)")" '-gt' "$(date +%s -d 'now - 1 seconds')" \
		"'cancelled_at' timestamp is outside expected range (now - 1 secs)"
	assert \
		"$(query "strftime('%s', cancelled_at)")" '-le' "$(date +%s -d 'now + 1 seconds')" \
		"'cancelled_at' timestamp is outside expected range (timestamp is in the future??)"

	assert "$(ctl_curl /alarm/${alarm_id}/dispatch-has-timer)" = "false" \
		"assert the alarm has NO timer configured in the dispatcher"

	assert "$(ctl_curl /alarm/${alarm_id}/dispatch-has-alarm)" = "false" \
		"assert the alarm was removed from the dispatcher"

}
