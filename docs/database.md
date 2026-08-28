# WiP - Database

WiP uses SQLite as database.

The default (development) location is `data/work-in-peace.sqlite`.

This or alternate locations can be specified at startup via: `--db.dsn
<string:path/to/db.sqlite>`

## Management

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

## Delete / Destroy database

** !! WARNING: this will permanently delete data !! **

```bash
make db.delete
	# >_ rm -v data/work-in-peace.sqlite
	# are you sure (y|yes|NO): y
	# removed 'data/work-in-peace.sqlite'
```
