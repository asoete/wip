# WiP - Database

WiP uses SQLite as database.

The default (development) location is `data/work-in-peace.sqlite`.

This or alternate locations can be specified at startup via: `--db.dsn
<string:path/to/db.sqlite>`

## MANAGE SQLITE DATABASE

As with most things related to WiP, a `make` "interface" is provided to manage
database interactions.

`make` variables / parameters

| NAME           | DEFAULT VALUE               |
|----------------|-----------------------------|
| DB.SQLITE_FILE | `data/work-in-peace.sqlite` |

### Init database

```bash
make db [DB.SQLITE_FILE=/path/to/db.sqlite]
```

`make db` is a shorthand for

```bash
make db.create
	# sqlite3 data/work-in-peace.sqlite < sql/schema/20260825145000_create-table-migrations.sql
	# 2026-08-27 16:18:16

make db.migrate
	# libexec/db.migrate.sh
	# [W] no SQL migration files specified, falling back to 'sql/schema/*.sql'
	#  [  NOOP  ]  sql/schema/20260825145000_create-table-migrations.sql        2026-08-27 16:18:16
	#  [MIGRATED]  sql/schema/20260825145030_create-table-users.sql             2026-08-27 16:18:50
	#  [MIGRATED]  sql/schema/20260825145205_create-table-alarms.sql            2026-08-27 16:18:50
	#  [MIGRATED]  sql/schema/20260825145411_create-table-ntfy_channels.sql     2026-08-27 16:18:50
	#  [MIGRATED]  sql/schema/20260825145522_create-table-user_channels.sql     2026-08-27 16:18:50
```

`make db.create` will create the SQLite file and create a single table: `migrations`.
(See: `sql/schema/20260825145000_create-table-migrations.sql`)

See [Migrations](#Migrations) for more info about `make db.migrate`.

### Migrations

#### Development

`make db.create` is a simple wrapper / _make interface_ to `libexec/db.migrate.sh`.

`libexec/db.migrate.sh` will scan `sql/schema/` for SQL files and use the
`migrations` SQLite table to determine if a migration-file needs to be executed
or if the migration was already applied.

See `libexec/db.migrate.sh -h` for more info about this program and its options.

You can ask `libexec/db.migrate.sh` to run migrations as many times as you
like. Only when the migration is not yet registered in the `migrations` table,
will the file be executed...

```bash
make db
	# sqlite3 data/work-in-peace.sqlite < sql/schema/20260825145000_create-table-migrations.sql
	# 2026-08-27 16:19:26
	# libexec/db.migrate.sh
	# [W] no SQL migration files specified, falling back to 'sql/schema/*.sql'
	#  [  NOOP  ]  sql/schema/20260825145000_create-table-migrations.sql        2026-08-27 16:19:26
	#  [MIGRATED]  sql/schema/20260825145030_create-table-users.sql             2026-08-27 16:19:26
	#  [MIGRATED]  sql/schema/20260825145205_create-table-alarms.sql            2026-08-27 16:19:26
	#  [MIGRATED]  sql/schema/20260825145411_create-table-ntfy_channels.sql     2026-08-27 16:19:26
	#  [MIGRATED]  sql/schema/20260825145522_create-table-user_channels.sql     2026-08-27 16:19:26

make db.migrate
	# libexec/db.migrate.sh
	# [W] no SQL migration files specified, falling back to 'sql/schema/*.sql'
	#  [  NOOP  ]  sql/schema/20260825145000_create-table-migrations.sql        2026-08-27 16:19:26
	#  [  NOOP  ]  sql/schema/20260825145030_create-table-users.sql             2026-08-27 16:19:26
	#  [  NOOP  ]  sql/schema/20260825145205_create-table-alarms.sql            2026-08-27 16:19:26
	#  [  NOOP  ]  sql/schema/20260825145411_create-table-ntfy_channels.sql     2026-08-27 16:19:26
	#  [  NOOP  ]  sql/schema/20260825145522_create-table-user_channels.sql     2026-08-27 16:19:26
```

#### Production

A problem arises when trying to bring these migration files into production. Do
we ship a makefile and a bunch of loose SQL files? Where do we store this data? ~~`/usr/share/wip`~~?

`tools/bundle-migrations.sh` to the rescue!

From the tool itself:
> This tools ingests all provided migration files and spits out a single
> bash-script.
> This bash script is idempotent an can thus be executed as much as you
> want. For example, at each service start...

The created script will:
1. Check if the provided database is initialized
1a. If not: initialize (by creating the `migrations` table 
2. Run all bundled SQL files (in order)

As usual, a make wapper is provided: `make rpm/migrate.sh`

```bash
# 1. Create migration bundle from sql files
make rpm/migrate.sh
	# tools/bundle-migrations.sh \
	#       sql/schema/20260825145000_create-table-migrations.sql \
	#       sql/schema/20260825145030_create-table-users.sql \
	#       sql/schema/20260825145205_create-table-alarms.sql \
	#       sql/schema/20260825145411_create-table-ntfy_channels.sql \
	#       sql/schema/20260825145522_create-table-user_channels.sql \
	#       > rpm/migrate.sh

# 2. Run migration script for the first time
sh rpm/migrate.sh /tmp/db.sqlite
	#  *** warning: new (uninitialized) database provided: '/tmp/db.sqlite'
	#  *** creating new database: '/tmp/db.sqlite' and populating with 'sql/schema/20260825145000_create-table-migrations.sql'
	#  *** init database '/tmp/db.sqlite' completed
	#  *** -rw-r--r--. 1 user user 16384 28 aug 17:22 /tmp/db.sqlite
	# sql/schema/20260825145000_create-table-migrations.sql           [  DONE  ] 2026-08-28 15:22:21
	# sql/schema/20260825145030_create-table-users.sql                [MIGRATED] 2026-08-28 17:22:21
	# sql/schema/20260825145205_create-table-alarms.sql               [MIGRATED] 2026-08-28 17:22:21
	# sql/schema/20260825145411_create-table-ntfy_channels.sql        [MIGRATED] 2026-08-28 17:22:21
	# sql/schema/20260825145522_create-table-user_channels.sql        [MIGRATED] 2026-08-28 17:22:21
	
# 2. Run migration script (many) times again
sh rpm/migrate.sh /tmp/db.sqlite
	# sql/schema/20260825145000_create-table-migrations.sql           [  DONE  ] 2026-08-28 15:22:21
	# sql/schema/20260825145030_create-table-users.sql                [  DONE  ] 2026-08-28 15:22:21
	# sql/schema/20260825145205_create-table-alarms.sql               [  DONE  ] 2026-08-28 15:22:21
	# sql/schema/20260825145411_create-table-ntfy_channels.sql        [  DONE  ] 2026-08-28 15:22:21
	# sql/schema/20260825145522_create-table-user_channels.sql        [  DONE  ] 2026-08-28 15:22:21
```

### Delete / Destroy database

** !! WARNING: this will permanently delete data !! **

```bash
make db.delete
	# >_ rm -v data/work-in-peace.sqlite
	# are you sure (y|yes|NO): y
	# removed 'data/work-in-peace.sqlite'
```

## USE DATABASE

In order to interface with the database, we need some dependencies

1. [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) ( (pure go) sqlite driver for stdlib [database/sql](https://pkg.go.dev/database/sql))
2. [`sqlc`](https://sqlc.dev/) (for generating typesafe database abstraction / interaction code)


### Installation

**`sqlc`** is only required during development as it is performing codegen based on SQL files...

I used `nix` to add this to my development environment

```bash
nix-shell -p sqlc

sqlc version
  # v1.31.1
```

**`modernc.org/sqlite`** should be added as a _normal_ go dependency:

```bash
cat go.mod
  # module vsc.irc.ugent.be/itsupport/work-in-peace
  # 
  # go 1.26.5

go get modernc.org/sqlite
  # go: downloading modernc.org/sqlite v1.58.0
  # go: downloading golang.org/x/sys v0.47.0
  # go: downloading modernc.org/libc v1.75.6
  # go: downloading github.com/ncruces/go-strftime v1.0.0
  # go: downloading modernc.org/mathutil v1.7.1
  # go: downloading github.com/google/uuid v1.6.0
  # go: downloading github.com/mattn/go-isatty v0.0.24
  # go: downloading modernc.org/memory v1.12.1
  # go: downloading github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec
  # go: added github.com/dustin/go-humanize v1.0.1
  # go: added github.com/google/uuid v1.6.0
  # go: added github.com/mattn/go-isatty v0.0.24
  # go: added github.com/ncruces/go-strftime v1.0.0
  # go: added github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec
  # go: added golang.org/x/sys v0.47.0
  # go: added modernc.org/libc v1.75.6
  # go: added modernc.org/mathutil v1.7.1
  # go: added modernc.org/memory v1.12.1
  # go: added modernc.org/sqlite v1.58.0

cat go.mod
  # module vsc.irc.ugent.be/itsupport/work-in-peace
  # 
  # go 1.26.5
  # 
  # require (
  #         github.com/dustin/go-humanize v1.0.1 // indirect
  #         github.com/google/uuid v1.6.0 // indirect
  #         github.com/mattn/go-isatty v0.0.24 // indirect
  #         github.com/ncruces/go-strftime v1.0.0 // indirect
  #         github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
  #         golang.org/x/sys v0.47.0 // indirect
  #         modernc.org/libc v1.75.6 // indirect
  #         modernc.org/mathutil v1.7.1 // indirect
  #         modernc.org/memory v1.12.1 // indirect
  #         modernc.org/sqlite v1.58.0 // indirect
  # )
```
