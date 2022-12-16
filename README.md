# mission

**mission** is a task runner designed for a (daily) system and cloud backup. It loads a set of tasks from a configuration file and executes them in order.

## Features

- Built in (written in Go) tasks or external scripts or binaries.
- Scheduled tasks to run as part of a system cronjob, as well as manually.
- Post-tasks to be run after commit and chown operations are complete. (Useful for when filesystem backups need to run after cloud backups.)
- Multiple types of storage directory (sync, local, ...).
- Handle chowning and committing files in storage directories.
- Helper scripts shared by all tasks.
- Ping (via HTTP) a remote endpoint when a run completes without errors.
- Send an email (via sendmail) with the log when a run completes.
- Show a desktop notification (via notify-send) when a run completes.
- Verbose timestamped logs that can be sent to stdout, saved to disk, and/or emailed.
- Filter out potentially sensitive lines from logs sent over email.
- Fetch credentials from [pass][pass] (or alternatives with the same API).

## Installation

```sh
go install github.com/crdx/mission@latest
```

## CLI

```
Usage:
    mission [options] run [--task SLUG...] [--headless]
    mission [options] list [--verbose]
    mission [options] check
    mission [options] dump

Commands:
    run     Run all tasks or specific tasks
    list    List all available tasks
    check   Validate configuration
    dump    Dump parsed configuration as JSON

Options:
    --headless           Run headlessly
    -c, --config PATH    Configuration file
    -t, --task SLUG      One or more tasks to run
    -q, --quiet          Be quiet
    -v, --verbose        Be verbose
    -C, --no-color       Disable colours
    -h, --help           Show help
```

## Commands

### run

Start a run.

Pass `-t/--task` one or more times to specify which tasks to run, or omit it to run all scheduled tasks.

The `--headless` flag indicates that this is a run triggered from a cron job, changing the behaviour in the following ways.

- Disable colour output.
- Send output to a logfile instead of stdout.
- Send an email on completion.
- Ping endpoint on successful completion.

### list

Show a list of available task slugs.

Pass `-v/--verbose` to see more detailed output.

### check

Validate the configuration file.

To avoid any surprises the configuration has to pass strict validation checks before any runs can start.

### dump

Parse the configuration file and dump it as JSON after transformations and validation checks have taken place.

## Configuration

The format of the configuration file is a more relaxed form of JSON ("jsonc") that supports comments using `//`.

If the path is not passed via `-c/--config` then it will be assumed to be `mission.json` in the working directory.

See `config.sample.json` for a sample configuration file, or carry on reading for a more detailed description of each section.

### tasks

The most important section: the list of tasks. Each task should have the following fields.

| Field | Description |
| ----- | ----------- |
| slug | A unique lowercase label identifying this task.<br><br>If this is an external task then this corresponds to the name of the directory within the `tasks` directory.<br><br>If this is a built in task then this corresponds to the value used by the `GetBuiltIn` method in `task/task.go`. |
| name | A user-friendly name for this task. |
| scheduled | Whether this task should run if called without any specified tasks with `-t/--task`.<br><br>Tasks that are not scheduled can still be run manually. |
| external | Whether this task is external or built in. |
| post | Whether this task should run after all other tasks have run and after commit and chown operations are complete. |
| entrypoint | (Optional) Path to script or binary entrypoint. Overrides the lookup done via `storage.tasks.path`, and only applies to external tasks. |

### user

Should contain one field, `name`, with the username of the system's unprivileged user.

This user is used as the file owner, and as the recipient of notifications and emails.

### passBin

Path to the `pass` binary.

This does not _have_ to be [pass][pass], but it should be a binary with an API compatible with your tasks. For built in tasks there is a `GetPassValue` helper available that will run it with a single argument, and for external tasks it's passed as an environment variable so it depends entirely on how your tasks invoke it.

A `~` prefix in the path will be transformed into the home directory of the user set in the `user.name` field.

### storage

A set of storage directories for tasks to use for backup data or logfiles. Each directory should have the following fields.

| Field | Description |
| ----- | ----------- |
| path | The path to the directory.<br><br>A `~` prefix in the path will be transformed into the home directory of the user set in the `user.name` field.
| chown | If `true` then after all the non-post tasks have run all files in this directory will be chowned to the user set in the `user.name` field.
| commit | If `true` then after all the non-post tasks have run all files in this directory will be committed with git.

Five types of storage directory are supported.

| Name | Purpose |
| ---- | ------- |
| sync | Backup files that are to be synced to other machines via some form of cloud sync service. |
| local | Backup files that are expected to remain on the local machine. |
| logs  | Logfiles that need to be saved. The main run log will be saved as a dated file here. |
| helpers | Helper scripts that need to be shared across tasks. |
| tasks | Directory where tasks can be found. Can be overridden on a per-task basis. |
| store | (Optional) Cached sessions. |

### ping

When a run completes with no errors this endpoint will be sent an HTTP GET. This can be used in combination with an uptime monitoring service like [Healthchecks][healthchecks] to ensure failing runs are not missed.

If the endpoint contains a `%s` then it will be replaced with the system hostname.

### notify

The user set in the `user.name` field will be notified with `notify-send` when the run starts and finishes.

### mail

Only one type is currently supported: `sendmail`.

An email will be sent via `sendmail` to the user set in the `user.name` field when the run starts and finishes.

The run log without filtered items (see [filters](#filters)) will be included in the email.

### filters

A list of regexes to use to filter out log lines before including them in the email. Some tasks may produce a lot of output which should still be saved to disk but would be excessive to include in an email. It can also be used to filter out sensitive lines.

## Types of task

### Built in

Built in tasks should be implemented in Go under the `tasks` directory. The directory corresponds to the task slug and should be the package name. The `GetBuiltIn` method in `task/task.go` should be modified to reference the task. Each task is a package named after the slug containing a `Run` method.

See the example in `tasks/spotify`.

### External

External tasks are executable scripts or binaries named `run` located within a directory named after the task slug in the `tasks` directory in the working directory. For example a task named `mail` would be resolved to the executable found at `tasks/mail/run`.

## Environment

Each task runs with environment variables available corresponding to certain fields in the configuration file.

| Variable | Field | Notes |
| -------- | ----- | ----- |
| `SYNC_FILES_DIR` | `storage.sync.path` | Resolved absolute path |
| `LOCAL_FILES_DIR` | `storage.local.path` | Resolved absolute path |
| `HELPERS_DIR` | `storage.helpers.path` | Resolved absolute path |
| `LOGS_DIR` | `storage.logs.path` | Resolved absolute path |
| `STORE_DIR` | `storage.store.path` | Resolved absolute path |
| `PASS_BIN` | `passBin` | Resolved absolute path |
| `TARGET_USER` |  `user.name` | |

## Contributions

Open an [issue](https://github.com/crdx/mission/issues) or send a [pull request](https://github.com/crdx/mission/pulls).

## Licence

[GPLv3](LICENCE).

[pass]: https://www.passwordstore.org
[healthchecks]: https://healthchecks.io
