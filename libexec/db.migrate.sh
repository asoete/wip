#!/usr/bin/sh

set -euo pipefail

DB_FILE="${DB_FILE:-data/work-in-peace.sqlite}"

VERBOSITY=${VERBOSITY:-1}
LVL_SILENT=-1
LVL_QUIET=0
LVL_INFO=1
LVL_DEBUG=2
LVL_TRACE=3

function main() {

	default_migration_source_dir="sql/schema"
	sql_file_pattern='*.sql'
	declare -a sources

	if [ $# -lt 1 ] ; then

		test $VERBOSITY -gt $LVL_QUIET &&
			printf "\e[33m[W] no SQL migration files specified, falling back to '%s/%s' \e[0m\n" \
				"${default_migration_source_dir}" \
				"${sql_file_pattern}" \
				1>&2

		readarray -t sources < <(find "${default_migration_source_dir}" -iname "${sql_file_pattern}" | sort -V)
	else
		sources=("$@")
	fi

	test $VERBOSITY -ge $LVL_DEBUG && printf "\e[37m sources: %s\e[0m\n" "${sources[@]}"

	for file in "${sources[@]}" ; do

		if [ -d "$file" ] ; then

			printf "\e[33m[W] detected a directory as input ($file). Scanning for SQL files... \e[0m\n" 1>&2

			for subfile in $(find "${file}" -iname "${sql_file_pattern}" | sort -V) ; do
				migrate_single_file "$subfile"
			done

		else

			migrate_single_file "$file"
		fi
	done

}

function migrate_single_file() {

	file="${1:-}"

	test $VERBOSITY -ge $LVL_DEBUG && printf "\e[037m process %s\e[0m\n" "$file"

	test -z "$file" && die "(<file>): no SQL file provided"

	test -f "$file" || die "migrate_single_file(<file>): no such file: %s." "$file"

	test $VERBOSITY -ge $LVL_TRACE && printf "\e[37m Check if file is already in migrations table \e[0m\n"

	sql_res=$(sqlite_run "SELECT file FROM migrations WHERE file = '%s';" "$file")

	if [ "${sql_res}" = "${file}" ] ; then

		printf "\e[32;4m [  NOOP  ]  %-60s \e[4m%s \e[0m\n" \
			"$file" \
			"$(sqlite_run "SELECT datetime(date, 'localtime') FROM migrations WHERE file = '%s';" "$file")"

		return
	fi

	mapfile migration_result < <(sqlite < "$file" | tr -d '\n')

	printf "\e[32;1;4m \e[7m[MIGRATED]\e[27m  %-60s \e[4m%s \e[0m\n" \
		"$file" \
		"${migration_result[@]}"


}

function sqlite_run() {

	test -z "${1:-}" && die "db.migrate.sh: sqlite_run(<query> ,[<params>...]): no query provided"

	printf "$@" | sqlite || die "sqlite_run(<query>, [<params>...]) failed"
}

function sqlite() {

	test -f "${DB_FILE}" || die 'db.migrate.sh: sqlite(<query> ,[<params>...]): DB_FILE not found: `%s`.\n(Run `make db.create` first... ?)' "${DB_FILE}"

	sqlite3 ${sqlite_flags:-} "${DB_FILE}" || die "sqlite3 $DB_FILE (pipe) failed"
}


function die() {

	{
		printf "\e[31;1m[E] db.migrate.sh: "
		printf "$@"
		printf "\n"
		printf "\e[0;31m"
		get_trace "\t |> "
		printf "\e[0m\n"
	} 1>&2

	exit 1
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

main "$@"
