# Script format

Name each scenario `<command>_<scenario>.txt` under
`cmd/layerctl/testdata/script`.
The runner starts each script in its own temporary directory, available as
`$WORK`, and puts `layerctl` on `PATH` using the Go test binary.
Git starts with a test author and committer, fixed dates for both, and no
system or user config.
Use `at 2000-01-02T03:04:05Z` to set both dates for subsequent commands
when a scenario depends on commit timing.
Choose explicit branch names when they matter to the scenario.

Testscript reads a command per line.
Use `#` comments for the promise or phases whose purpose needs explaining.
`exec` runs a program and expects success, and `! exec` expects failure.
Use `git` for Git commands in the script's working directory.
Use `! git` when the command should fail.
Single quotes group arguments and suppress variable expansion.
Shell pipes, redirections, and substitutions require an explicitly invoked
shell.
Use built-in commands such as `cd`, `cp`, `exists`, and `cmp` for fixture
operations.
Lines of the form `-- path --` at the end of the script introduce files
unpacked under `$WORK` before the first command.

Put input files under `give/` and expected output files under `want/`.
Use `give/repo/` for the initial contents of a Git repository and
`give/layer/` for its layer manifest, patches, and other layer files.
Include the patch itself when its contents are part of the setup.
When a command writes layer files, put the expected manifest and patches
under `want/layer/`.
Omit the expected manifest when the command leaves it unchanged.

This fragment checks the path printed for a projection named `demo` that the
preceding setup created at `$WORK/workspace`.

```text
exec layerctl projection path demo
cmpenv stdout $WORK/want/path

-- want/path --
$WORK/workspace
```

`cmp` compares file contents and accepts `stdout` and `stderr` from the last
program as inputs.
`cmpenv` expands variables in the expected file.
Use `cp stdout saved` to keep output for a later comparison.
`stdout` and `stderr` take regular expressions.
Use `${WORK@R}` outside single-quoted segments to escape regex metacharacters
in the work directory path.

Run the full script suite or select one script by its filename stem.

```sh
go test ./cmd/layerctl -run '^TestScript$'
go test ./cmd/layerctl -run '^TestScript$/^help$' -v
```

See the [testscript package documentation](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript)
for the full command reference.
