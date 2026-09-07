# =============================================================================
# DIFF
# =============================================================================

# Wrapper around difftastic (or any diff utility)
# Compares the output of 2 provided function ($1() <> $2()) and reports
# differences
function diff() {
	#delta --diff-args=-U99 --diff-so-fancy --detect-dark-light always <($1) <($2)
	difft \
		--skip-unchanged \
		--exit-code \
		--color always \
		--override '*:JSON' \
		--context 99 \
		<($1) <($2)
}

# =============================================================================
# HTTP
# =============================================================================

function wip_curl() {

	path="${1:-/}" ; shift
	url="${BASE_URL}${path}"

	#log "CURL %s %s" "${url}" "$*"
	log "CURL %s %s" "${path}" "$*"

	run -0 --separate-stderr \
		curl -i -s "${url}" \
		-w "%{stderr}Status: %{http_code}\n" \
		"$@"
}

# =============================================================================
# SQLITE
# =============================================================================

# Execute SQLITE command and
#	-> format output as JSON
#	-> sort
#	-> if empyt result, replace with []
function sqlite_json() {

	printf "SQLITE3 : %s" "$*" 1>&2

	mapfile -t result < <(sqlite3 -json "${DB_FILE}" "$@")

	if [ ${#result[@]} -le 0 ] ;then
		printf "   (0 rows)\n" 1>&2
		printf "[]\n"
	else 
		printf "   (%d rows)\n" "${#result[@]}" 1>&2
		printf "%s\n" "${result[@]}" | jq . --sort-keys
	fi
}

# =============================================================================
# FMT
# =============================================================================

# Pass inputs to jo (1) and format sorted via jq (1)
function json_fmt() {
	jo -a "$@" | jq --sort-keys
}

# =============================================================================
# ASSERT
# =============================================================================

function assert_equal() {

	test "$1" = "$2" || {
		test -n "$3" && printf "  *** $3 *** \n"
		printf "  ASSERT_EQUAL FAILED: \$1(%s) != \$2(%s)\n" "$1" "$2"
		get_trace "\t |> "
		false
	}
}

function get_trace() {

	prefix="${1:-}"
	nr_fns=${#FUNCNAME[@]}

		for index in $(seq 2 $((${nr_fns}-1)) ) ; do
		get_trace_index $index "$prefix"
	done
}

function get_trace_index() {

	index=$1
	prefix="${2:-}"

	printf "$prefix\e[3m%s\e[23m called at %s:%s\n" \
		"${FUNCNAME[$index]}" \
		"${BASH_SOURCE[$(($index+1))]}" \
		"${BASH_LINENO[$index]}"
}

# =============================================================================
# LOGGING
# =============================================================================

function log() {

	printf "[I] "
	printf "$@"
	printf "\n"
}
