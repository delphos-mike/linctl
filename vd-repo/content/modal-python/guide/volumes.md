# Volumes

Modal Volumes provide a high-performance distributed file system for Modal applications. They are designed for write-once, read-many I/O workloads, like creating machine learning model weights and distributing them for inference.

## Creating a Volume

Create a Volume using the CLI command:

```bash
% modal volume create my-volume
Created volume 'my-volume' in environment 'main'.
```

## Using a Volume on Modal

Attach an existing Volume to a Modal Function:

```python
vol = modal.Volume.from_name("my-volume")

@app.function(volumes={"/data": vol})
def run():
    with open("/data/xyz.txt", "w") as f:
        f.write("hello")
    vol.commit()  # Persist changes
```

You can also browse Volumes in a Modal Shell:

```bash
% modal shell --volume my-volume --volume another-volume
```

### Creating Volumes Lazily

Create Volumes programmatically:

```python
vol = modal.Volume.from_name("my-volume", create_if_missing=True)
```

## Using a Volume from Outside Modal

### Local Code Interaction

```python
vol = modal.Volume.from_name("my-volume")

with vol.batch_upload() as batch:
    batch.put_file("local-path.txt", "/remote-path.txt")
    batch.put_directory("/local/directory/", "/remote/directory")
    batch.put_file(io.BytesIO(b"some data"), "/foobar")
```

### Command Line Interface

Use `modal volume` for various operations:

```bash
% modal volume
# Shows available subcommands for volume management
```

## Volume Commits and Reloads

Volumes require explicit reloading to see changes. Changes must be committed to become visible outside the current container.

### Background Commits

Modal automatically commits Volume changes periodically during function execution.

## Model Serving and Checkpointing

Volumes are useful for storing and serving ML models:

```python
import modal

app = modal.App()
volume = modal.Volume.from_name("model-store")
model_store_path = "/vol/models"

@app.function(volumes={model_store_path: volume})
def save_model():
    # Save model to volume
    pass

@app.function(volumes={model_store_path: volume})
def load_model():
    # Load model from volume
    pass
```