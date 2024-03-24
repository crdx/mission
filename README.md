# mission

**mission** is a task runner designed for a (daily) system and cloud backup. It loads a set of tasks from a configuration file and executes them in order.

## Features

- Tasks can be executable scripts, binaries, or built in (written in Go).
- Schedule tasks to run as part of a system cronjob, or manually.
- Robust [failure handling](#failure-handling) ensures failing tasks are not missed.
- Post-tasks which run after other tasks. Useful for when filesystem backups need to run after cloud backups.
- Configurable storage directories available to each task via its environment.
- Handle [chowning](#chown) and committing files in storage directories.
- Ping (via HTTP) a remote endpoint when a run completes without errors.
- Send an email (via sendmail) with the log when a run completes.
- Show a desktop notification (via notify-send) when a run completes.
- Verbose timestamped logs that can be sent to stdout, saved to disk, and/or emailed.
- Filter out potentially sensitive lines from logs sent over email.
- Fetch credentials from [pass][pass] (or alternatives with the same API).

## Installation

```sh
go install crdx.org/mission@latest
```

## CLI

```
Usage:
    mission [options] init
    mission [options] run [--task SLUG...] [--headless]
    mission [options] list [--verbose]
    mission [options] check
    mission [options] dump

Commands:
    init     Create a sample configuration file
    run      Run all tasks or specific tasks
    list     List all available tasks
    check    Validate configuration
    dump     Dump parsed configuration as JSON

Options:
    --headless           Run headlessly
    -c, --config PATH    Configuration file
    -t, --task SLUG      One or more tasks to run
    -q, --quiet          Be quiet
    -v, --verbose        Be verbose
    -C, --no-color       Disable colours
```

## Commands

### init

Create a sample configuration file in the working directory.

The contents of the file will be the as `mission.sample.json` found in this repository.

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

If the path is not passed via `-c/--config` then it will be assumed to be `mission.config.json` in the working directory.

A prefixed `~` in any path will be transformed into the home directory of the user set in the `user.name` field.

See `mission.sample.json` for a sample configuration file, or carry on reading for a more detailed description of each section.

### tasks

The most important section: the list of tasks. Each task should have the following fields.

| Field | Description |
| ----- | ----------- |
| slug | A unique lowercase label identifying this task.<br><br>If this is an `exec` task then this is the name of the directory within the `tasks` directory.<br><br>If this is a `builtin` task then this is the value used by the `GetBuiltInAction` method in `config/task.go`. |
| name | A user-friendly name for this task. |
| type | Type of the task. Should be `exec` or `builtin`. See [Task types](#task-types) for more information. |
| scheduled | Whether this task should run if called without any specified tasks with `-t/--task`.<br><br>Tasks that are not scheduled can still be run manually. |
| post | Whether this task should run after all other tasks have run and after commit and [chown](#chown) operations are complete. |
| entrypoint | (Optional) Path to executable entrypoint. Overrides the lookup done via `storage.tasks.path`, and only applies to exec tasks. |

### user

User-related configuration.

| Field | Description |
| ----- | ----------- |
| name | The username of the system's unprivileged user. Used as the file owner and the recipient of notifications. |
| email | Address used as the recipient of emails. Only mandatory if `mail.enabled` is `true`. |

### passBin

Path to the `pass` binary.

This does not _have_ to be [pass][pass], but it should be a binary with an API compatible with your tasks. For built in tasks there is a `GetPassValue` helper available that will run it with a single argument, and for exec tasks it's passed as an environment variable so it depends entirely on how your tasks invoke it.

### storage

A set of storage directories for tasks to use for backup data or logfiles. Each directory should have the following fields.

| Field | Description |
| ----- | ----------- |
| path | The path to the directory.
| chown | If `true` then after all the non-post tasks have run all files in this directory will be [chowned](#chown) to the user set in the `user.name` field.
| commit | If `true` then after all the non-post tasks have run all files in this directory will be committed with git.

Two storage directories are mandatory.

| Name | Purpose |
| ---- | ------- |
| tasks | Directory where tasks are kept. Can be overridden on a task-by-task basis via the task definition. |
| logs | Logfiles that need to be saved. The main run log will be saved as a timestamped file here. |

Other storage directories depend on your workflow and needs. Some examples:

| Name | Purpose |
| ---- | ------- |
| sync | Backup files that are to be synced via some form of cloud sync service. |
| local | Backup files that are expected to remain on the local machine. |
| helpers | Helper scripts that need to be shared across tasks. |
| sessions | Cached sessions. |

All storage directories are available to tasks via [environment variables](#environment).

### ping

When a run completes with no errors this endpoint will be sent an HTTP GET. This can be used in combination with an uptime monitoring service like [Healthchecks][healthchecks] to ensure failing runs are not missed.

If the endpoint contains a `%s` then it will be replaced with the system hostname.

### notify

The user set in the `user.name` field will be notified with `notify-send` when the run starts and finishes.

### mail

Only one type is currently supported: `sendmail`.

An email will be sent via `sendmail` to the email set in the `user.email` field when the run starts and finishes.

The run log without filtered items (see [filters](#filters)) will be included in the email.

### filters

A list of regexes to use to filter out log lines before including them in the email. Some tasks may produce a lot of output which should still be saved to disk but would be excessive to include in an email. It can also be used to filter out sensitive lines.

## Failure handling

Ensuring that failures are brought to the user's attention is important. If a task exits with a non-zero exit code it will be considered to have failed. All other tasks will still be run to prevent a single task's failure affecting the rest of the backup, however, once a task fails the overall run's behaviour changes as follows:

- The process will exit with a non-zero exit code (1) and print the error count to stderr. The convention is that a cronjob that outputs to stderr has something important to say.
- The email subject will be prefixed with "[FAILURE]".
- The ping endpoint will NOT be pinged.
- Errors will be summarised in a section at the bottom of the run log.

## Chown?

It may not be obvious why chowning files would be the job of a task runner at all.

If a backup task needs to run as root to access the resources it needs, it may also need to create files (containing the backed up data) which will by default be owned by root. If they should ultimately end up owned by the unprivileged user then there are two options: ensure the task creates the files with the right owner, or let the task create all the root-owned files it likes and let the runner handle the chowning job at a later stage in the pipeline.

This workflow, though slightly unconventional, is entirely optional. It's flexible enough that you have the choice of running as root and letting specific tasks drop privileges, or running unprivileged and letting specific tasks request higher privileges.

Note that the run log will be owned by the user that runs **mission**.

## Task types

### Exec

Exec tasks are executable scripts or binaries named `run` located within a directory named after the task slug in the `tasks` storage directory. For example, a task named `mail` would be resolved to the executable found at `$TASKS_DIR/mail/run`.

### Built in

Built in tasks should be implemented in Go under the `tasks` directory. The directory corresponds to the task slug and should be the package name. The `GetBuiltInAction` method in `config/task.go` should be modified to reference the task. Each task is a package named after the slug containing a `Run` method. If you want to do this then you probably want to maintain a fork of this repository for yourself.

See the example in `tasks/spotify`.

## Environment

Each task runs with environment variables available corresponding to certain fields in the configuration file.

| Variable | Field |
| -------- | ----- |
| `LOGS_DIR` | `storage.logs.path` |
| `TASKS_DIR` | `storage.tasks.path` |
| `PASS_BIN` | `passBin` |
| `TARGET_USER` |  `user.name` |

Any additional storage directories will also be available in the environment following the same convention as above.

If the directory is defined as

```json
"storage": {
    "foo": { "path": "~/foo" },
}
```

then the corresponding environment variable will be `FOO_DIR`.

## Contributions

Open an [issue](https://github.com/crdx/mission/issues) or send a [pull request](https://github.com/crdx/mission/pulls).

## Licence

[GPLv3](LICENCE).

[pass]: https://www.passwordstore.org
[healthchecks]: https://healthchecks.io
