#!/usr/bin/sh

###
### This tools ingests all provided migration files and spits out a single
### bash-script.
### This bash script is idempotent an can thus be executed as much as you
### want. For example, at each service start...
###
### USAGE:
###     bundle-migrations.sh <SQL-file> [<SQL-file>...]
###
### EXAMPLE:
###     bundle-migrations.sh sql/schema/*.sql > bundle.sh

set -euo pipefail

CREATE_TABLE_MIGRATIONS_FILE="sql/schema/20260825145000_create-table-migrations.sql"

function main() {

	gen_head "$@"
	gen_usage
	gen_main "$@"
	gen_functions "$@"
	gen_init_new_database
	gen_tail

}

function gen_head() {

	printf '#!/usr/bin/sh\n'

	printf "\n"
	printf '###################################################################### \n'
	printf '###                                                                ###\n'
	printf '###          AUTO GENERATED DO NOT MODIFY (MANUALLY)               ###\n'
	printf '###                                                                ###\n'
	printf '###################################################################### \n'

	printf "\n"
	printf '### \n'
	printf '### Auto-generated via %s \n' "$0"
	printf '### Created: %s\n' "$(date)"
	printf '### Version: %s\n' "$(git describe)"
	printf '### Commit: %s\n' "$(git rev-parse HEAD)"
	printf '### Input-files: %s\n' "$@"
	printf '### \n'

	printf "\n"
	printf 'set -euo pipefail\n'

	printf "\n"
	printf 'DB_FILE="${1:-${DB_FILE:-}}"\n'
}

function gen_main() {

	printf "function main() {\n"

	gen_guard_database_exists

	for filename in "$@" ; do
		function_name="$(file2function "$filename" )"

		printf "\t%s\n" "$function_name"
	done

	printf "}\n\n"
}

function gen_functions() {

	for filename in "$@" ; do

		function_name="$(file2function "$filename" )"

		printf "function %s() {\n" "$function_name"

		printf "\n"
		printf '\tselected_filename="$(%s)"\n' "$(gen_sqlite3_run "SELECT file FROM migrations WHERE file = '%s';" "$filename")"

		printf "\n"
		printf '\tprintf "%%-64s" "%s"\n' "$filename"

		printf "\n"
		printf '\tif [ "$selected_filename" = "%s" ] ; then\n' "$filename"
		printf '\t\tprintf "[  DONE  ] %%s\\n" "$(%s)"\n' "$(gen_sqlite3_run "SELECT date FROM migrations WHERE file = '%s';" "$filename")"
		printf '\t\treturn\n'
		printf '\tfi\n'

		printf "\n"
		printf "\tmigrate_output=\$(cat <<'EOSQL' | sqlite3 \"\${DB_FILE}\"\n"
		cat $filename
		printf 'EOSQL\n'
		printf ')\n'

		printf "\n"
		printf '\tprintf "[MIGRATED] %%s\\n" "$migrate_output"\n'

		printf "}\n"
		printf "\n"
	done
}


function gen_tail() {

	cat <<'EOHEAD'
main "$@"
EOHEAD
}

function gen_guard_database_exists() {

	cat <<'EOGRD'
	if [ -z "${DB_FILE}" ] ; then
		printf " *** error: no database provided\n" 1>&2
		usage 1>&2
		exit 1
	fi

	if [ ! -f "${DB_FILE}" ] ; then
		printf " *** warning: new (uninitialized) database provided: '%s'\n" "${DB_FILE}" 1>&2
		init_new_database
	fi

	if [ ! -r "${DB_FILE}" ] ; then
		printf " *** error: invalid database provided: '%s': unreadable\n" "${DB_FILE}" 1>&2
		exit 1
	fi

EOGRD

}

function gen_init_new_database() {

	printf 'function init_new_database() {\n'

		printf "\n"
		printf '\tprintf " *** creating new database: ' ; printf "'%s' " '${DB_FILE}' ; printf "and populating with '%s'" "$CREATE_TABLE_MIGRATIONS_FILE" ;  printf '\\n"\n'

		printf "\tmigrate_output=\$(cat <<'EOSQL' | sqlite3 \"\${DB_FILE}\"\n"
		cat $CREATE_TABLE_MIGRATIONS_FILE
		printf 'EOSQL\n'
		printf ')\n'

		printf '\tprintf " *** init database ' ; printf "'%s' " '${DB_FILE}' ; printf 'completed\\n"\n'
		printf '\tprintf " *** %%s\\n" "%s" \n' '$(ls -l "${DB_FILE}")'


	printf '}\n'

}

function gen_sqlite3_run() {

	printf 'sqlite3 "${DB_FILE}" "'
	printf "$@"
	printf '"'
}

function file2function() {

	if [ -z "${1:-}" ] ; then
		printf "error: file2function(<file>) requiers a file" 1>&2
		exit 1
	fi

	printf "$@" | tr '/.' '_'
}

function gen_usage() {

	cat <<'EOF'
function usage() {
	echo TODO
}

EOF
}

main "$@"
