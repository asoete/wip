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
		--display side-by-side-show-both \
		--color always \
		--override '*:JSON' \
		--context 99 \
		--ignore-comments \
		<(echo '// *** EXPECTED *** //' ; $1) <(echo '// *** RECEIVED *** //' ; $2)
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

	printf "[curl debug]: /usr/bin/curl" 1>&2
	for arg in "${curl_args[@]}" ; do
		case "$arg" in
			-*)
				printf "\n[curl debug]:    %s" "$arg" 1>&2
				;;
			http*)
				printf "\n[curl debug]:    %s" "$arg" 1>&2
				;;
			*)
				printf " '%s'" "$arg" 1>&2
				;;
		esac

	done
	printf "\n" 1>&2

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

	printf "[SQLITE3] : %s" "$*" 1>&2

	mapfile -t result < <(sqlite3 -json "${DB_FILE}" "$@")

	if [ ${#result[@]} -le 0 ] ;then
		printf "   (0 rows)\n" 1>&2
		printf "[]\n"
	else 
		printf "   (%d rows)\n" "${#result[@]}" 1>&2
		printf "%s\n" "${result[@]}" | jq . --sort-keys || {
			printf "[SQLITE3] unable to parse JSON response:\n" 1>&2
			printf "[SQLITE3 response]:  %s\n" "${result[@]}" 1>&2
			false
		}
		printf "%s\n" "${result[@]}" | sed -e 's/^/[SQLITE3 response]: /' 1>&2
	fi
}

function sqlite {

	printf "[SQLITE3] : %s\n" "$*" 1>&2

	sqlite3 "${DB_FILE}" "$@"

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

function assert() {

	if [ $# -lt 3 ] ; then
		printf "assert() requires <\$left> <\$cmp> <\$right>" 1>&2
		get_trace "\t |> " 1>&2
		false
	fi

	printf "[ASSERT]: '%s' %s '%s'" "$1" "$2" "$3" 1>&2

	left="$1" ; shift
	cmp="$1" ; shift
	right="$1" ; shift

	test "$left" "$cmp" "$right" && {
		printf " => <OK> \n"
	} || {

		printf " => <FAILED> \n"
		if [ $# -gt 0 ] ; then
			printf "ASSERTION ERROR: " 1>&2
			printf "$@" 1>&2
			printf "\n" 1>&2
		fi

		printf "  ASSERT FAILED: \$left(%s) \$cmp(%s) \$right(%s)\n" "$left" "$cmp" "$right" 1>&2
		get_trace "\t |> " 1>&2

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

# =============================================================================
# MISC
# =============================================================================

function rid() {

	n="${1:-3}"

	base32 /dev/urandom | head -c $n
}
