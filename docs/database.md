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

### [Migrations](Migrations)

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


## Delete / Destroy database

** !! WARNING: this will permanently delete data !! **

```bash
make db.delete
	# >_ rm -v data/work-in-peace.sqlite
	# are you sure (y|yes|NO): y
	# removed 'data/work-in-peace.sqlite'
```
