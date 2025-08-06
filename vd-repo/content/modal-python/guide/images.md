# Images in Modal

This guide explains how to define container environments for Modal Functions.

## Overview

Containers are lightweight virtual machines that isolate programs, making execution environments more reproducible. Modal uses the sandboxed gVisor container runtime for added security.

By default, Modal Functions run in a Debian Linux container with a Python installation matching your local Python version.

## Typical Image Definition

The standard approach involves method chaining to customize a base image:

```python
import modal

image = (
    modal.Image.debian_slim(python_version="3.10")
    .apt_install("git")
    .pip_install("torch==2.6.0")
    .env({"HALT_AND_CATCH_FIRE": "0"})
    .run_commands("git clone https://github.com/modal-labs/agi && echo 'ready to go!'")
)
```

## Adding Python Packages

Use `uv_pip_install` to add Python packages:

```python
datascience_image = (
    modal.Image.debian_slim(python_version="3.10")
    .uv_pip_install("pandas==2.2.0", "numpy")
)
```

## Adding Local Files

You can add local files using `add_local_dir` and `add_local_file`:

```python
image = modal.Image.debian_slim().add_local_dir("/user/erikbern/.aws", remote_path="/root/.aws")
```

## Running Shell Commands

Use `.run_commands` to execute shell commands during image build:

```python
image_with_model = (
    modal.Image.debian_slim().apt_install("curl").run_commands(
        "curl -O https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalcatface.xml",
    )
)
```

## Running Python Functions During Build

Use `.run_function` to execute Python functions during image setup:

```python
def download_models() -> None:
    import diffusers
    model_name = "segmind/small-sd"
    pipe = diffusers.StableDiffusionPipeline.from_pretrained(model_name)
```