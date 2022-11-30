# Watchexec Investigation

[`watchexec`](https://github.com/watchexec/watchexec) is a tool that watches the filesystem for modifications and runs a command when it detects them. The [Reloadable Process Types RFC](https://github.com/paketo-buildpacks/rfcs/blob/2492ad84fd3628a2a332010abcc49d4379d4da49/text/0032-reloadable-process-types.md) proposed that the Paketo project add a buildpack that installs `watchexec`. Ultimately, other Paketo buildpacks will likely use `watchexec` to implement the `BP_LIVE_RELOAD_ENABLED` feature proposed in that RFC. To begin to understand how buildpacks can use `watchexec`, we have investigated how `watchexec` itself behaves.

`watchexec` is a standalone binary that takes a command as its argument. For example:
```bash
watchexec -- echo hello
```

When the above runs, `watchexec` will run `echo` whenever files in the current directory (or subdirectories) change.

## Process Management

There are some nuances under-the-hood to how `watchexec` accomplishes process reloading. 

### Process Groups
First, when `watchexec` is called, it creates a process group and runs the desired command within that process group. `watchexec` is the parent of that process group. As a result, if `watchexec` receives a signal (e.g. `SIGTERM`), it will pass that signal to any and all processes created by the command it's managing. A user can disable the process-group feature with `--no-process-group`. 

### Shell Invocation
Also, `watchexec` runs the command in a shell by default. In other words, running `watchexec -- echo hello` results in a process tree like:
```bash
watchexec -- echo hello
  \_ sh -c echo hello
      \_ echo hello
```
`watchexec` invokes commands directly when the `--shell=none` flag is set. This will be necessary when `watchexec` is used in run images without shell (e.g. Paketo Tiny).

### Interrupting the Child Process
By default, if `watchexec` detects file modifications while the managed command is still running, it will queue another invocation of the command. The queued invocation will start once the previous one completes. But most buildpacks would probably use `watchexec` to run commands more like `node server.js`, where `server.js` is a long-running server, not a one-off task. To achieve "hot reloading" of the server
when files change, buildpacks would need to use `watchexec`'s `--restart` flag.
The flag tells `watchexec` to interrupt the process it's managing
and start a new one. To stop the current child, it sends `SIGTERM` (by
default). It creates a new child process once it receives `SIGCHLD` from the
previous one.

Users can set `--signal` to configure the signal (e.g. `SIGINT`) that `watchexec` sends to terminate its child when it detects file changes. This may be useful if the process `watchexec` is managing has certain signal
handlers that a user does/doesn't want to invoke.


## Filtering File Change Events
By default, `watchexec` triggers on modifications to any files in its current working directory (and subdirectories). Several flags
enable users to fine-tune which files should be watched.

The `--watch` flag `watchexec` overrides the default behaviour of watching the
current directory. Instead, `watchexec` will watch the specified files/directories (and their subdirectories). This flag can be included multiple times. For instance,
```bash
watchexec --watch ./src --watch ./node-modules node server.js.
```
will watch for changes in the `src` and `node-modules` subdirectories.

The flags `--exts` and `--filter` tell `watchexec` can further narrow the set of watched files. `--exts` takes a comma-separated list of file
extensions (e.g. `--exts js,txt`). Given that flag, `watchexec` will only trigger on changes to files with those extensions. `--filter` is similar, but accepts file glob patterns.

Taken together, these options can be useful in a buildpacks context for watching only some files for changes. For instance,
```bash
watchexec --restart --watch /workspace --exts js node server.js
``` 
will only restart the server when `*.js` files inside the app working directory change.

The `--ignore` flag is the negation of `--filter`. Instead of being selected, files that match provided globs will be ignored. Adding to the previous example, something like:
```bash
watchexec --restart --watch /workspace --exts js --ignore **/node_modules -- echo hi
``` 
will _not_ trigger if `*.js` files inside `/workspace/node_modules` change.

By default, watchexec also excludes a set of
commonly-ignored glob patterns.
TODO: What are these patterns?
The `--no-default-ignore` flag disables this behaviour. Watchexec will also exclude
patterns which appear inside a `.gitignore` or `.ignore` file in the project
directory. This can be disabled using the `--no-vcs-ignore flag`.


## Unresolved Questions
* Are there reasons why buildpacks users wouldn't want `watchexec` to manage their app's process in a process group?
* What are the implications of `watchexec` wrapping commands in a shell for signal handling within app containers?

