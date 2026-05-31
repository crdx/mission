# Changelog

## [1.1.0] - 2023-01-20

- Add `init` command.

## [1.0.5] - 2023-01-20

- Upgrade dependencies.
- Fix tests.

## [1.0.4] - 2023-01-19

- Switch internal testing library to testify.
- Upgrade dependencies.

## [1.0.3] - 2023-01-18

- Update import paths to use a vanity domain.

## [1.0.1] - 2023-01-17

- Move code out of `src/` as Go does not support such a directory structure. Put it in `internal/` instead.

## [1.0.0] - 2023-01-16

- Initial release.

## [0.4.1] - 2023-01-12

- Various code cleanups.

## [0.4.0] - 2022-12-28

- Search for `notify-send` and `sendmail` in `PATH`.
- Send emails to address defined in config file.
- Update error messages.

## [0.3.0] - 2022-12-17

- Change default config filename (yes, again) from `mission.json` to `mission.config.json`.
- Add storage directory validation to spotify task.
- Make passBin optional.
- Make storage directories fully configurable.
- Use a `type` string key to specify task type instead of a bool.
- Don't show headers for sections (chown and commit) when there is nothing to do.
- Various code cleanups.

## [0.2.0] - 2022-12-15

- Change default config filename from `config.json` to `mission.json`.

## [0.1.0] - 2022-12-15

- Initial release.
