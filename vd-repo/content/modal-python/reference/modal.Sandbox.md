# modal.Sandbox

A `Sandbox` class for interacting with running sandboxes, similar to Python's `asyncio.subprocess.Process`.

## Methods

### `hydrate()`
Synchronize the local object with its server identity. Rarely needs explicit calling.

### `create()`
Create a new Sandbox to run untrusted code. Key parameters include:
- `app`: Associate sandbox with an app
- `image`: Container image to run
- `secrets`: Environment variables
- `timeout`: Maximum execution time
- `gpu`: GPU configuration
- `network_file_systems`: File system mounts

### `from_name()`
Retrieve a running Sandbox by name from a specific app.

### `from_id()`
Construct a Sandbox using its unique identifier.

### `set_tags()`
Add key-value tags to the Sandbox for filtering.

### `snapshot_filesystem()`
Create an `Image` from the current Sandbox filesystem.

### `wait()`
Wait for Sandbox execution to complete.

### `tunnels()`
Retrieve tunnel metadata for the sandbox.

### `reload_volumes()`
Reload all Volumes mounted in the Sandbox.

### `terminate()`
Stop Sandbox execution.

### `poll()`
Check Sandbox running status, returning exit code if completed.

### `exec()`
Execute a command within the Sandbox.

### `open()`
Open a file in the Sandbox.

### `ls()`
List directory contents in the Sandbox.

### `mkdir()`
Create a new directory in the Sandbox.

### `rm()`
Remove a file or directory in the Sandbox.

### `watch()`
Monitor a file or directory for changes.

### `list()`
List Sandboxes, optionally filtered by app or tags.

## Properties
- `stdout`: Stream reader for standard output
- `stderr`: Stream reader for standard error
- `stdin`: Stream writer for standard input
- `returncode`: Process exit code