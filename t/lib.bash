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

function curl_wrapped() {

	declare -a curl_args

	curl_args+=("--silent")

	for arg in "$@" ; do

		case "$arg" in
			+ecode)
				curl_args+=("--write-out" "%{stderr}Status: %{http_code}\n")
				;;
			+elocation)
				curl_args+=("--write-out" "%{stderr}Location: %{redirect_url}\n")
				;;
			+headers)
				curl_args+=("--show-headers")
				;;
			+*)
				printf "curl_wrapped(): invalid expansion code: %s\n" "$arg" 1>&2
				get_trace "\t |> " 1>&2
				exit 1
				;;
			*)
				curl_args+=("$arg")
				;;
		esac

	done

	curl "${curl_args[@]}"
}

function wip_curl() {

	path="${1:-/}" ; shift
	url="${WIP_URL}${path}"

	curl_wrapped "${url}" "$@"
}

# ----------------------------------------------------------------------------

function ctl_curl() {

	path="${1:-/}" ; shift

	# If path is not prefixed with /ctl -> add the prefix
	if [[ "$path" != "/ctl/" ]] ; then
		path="/ctl${path}"
	fi

	url="${CTL_URL}${path}"

	curl_wrapped -H "${CTL_AUTH_HEADER}" "${url}" "$@"
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

	printf "[I] " 1>&2
	printf "$@" 1>&2
	printf "\n" 1>&2
}
