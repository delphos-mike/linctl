# modal.Volume

A writeable volume for sharing files between Modal functions.

## Overview

`Volume` is a class that allows persistent file storage and sharing across Modal functions. Key characteristics:
- Explicitly reload to see changes
- Explicitly commit changes to persist them
- Not ideal for concurrent file modifications
- Supports read-only mounting

## Class Definition

```python
class Volume(modal.object.Object):
    # Methods and properties
```

## Methods

### `hydrate()`
Synchronize local object with Modal server metadata.

### `read_only()`
Configure Volume to mount as read-only. Prevents file system write operations.

### `from_name(name, create_if_missing=False)`
Reference a Volume by name, optionally creating if it doesn't exist.

### `ephemeral()`
Create a temporary volume within a context manager.

### `info()`
Return information about the Volume object.

### `commit()`
Persist changes made to the volume.

### `reload()`
Make latest committed state available in the running container.

### `iterdir(path, recursive=True)`
Iterate over files in a directory.

### `listdir(path, recursive=False)`
List files under a path prefix.

### `read_file(path)`
Read a file from the volume.

### `remove_file(path, recursive=False)`
Remove a file or directory from the volume.

### `copy_files(src_paths, dst_path, recursive=False)`
Copy files within the volume.

### `batch_upload(force=False)`
Initiate a batched upload to a volume.

### Static Methods
- `delete(name)`
- `rename(old_name, new_name)`

## Example Usage

```python
import modal

volume = modal.Volume.from_name("my-volume", create_if_missing=True)

@app.function(volumes={"/data": volume})
def example_function():
    with open("/data/file.txt", "w") as f:
        f.write("hello")
    volume.commit()
```