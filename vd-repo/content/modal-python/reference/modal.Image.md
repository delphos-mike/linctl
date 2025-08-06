# modal.Image Class Reference

## Overview
A base class for container images to run functions in Modal. Do not construct directly; use static factory methods.

## Static Factory Methods

### `debian_slim(python_version: Optional[str] = None)`
Creates a default image based on official Python Docker images.

### `from_registry(tag: str)`
Builds an image from a public or private image registry.

### `from_dockerfile(path: Union[str, Path])`
Builds an image from a local Dockerfile.

## Key Methods

### `pip_install(*packages: str)`
Install Python packages using pip.

**Example:**
```python
image = modal.Image.debian_slim().pip_install("click", "httpx~=0.23.3")
```

### `add_local_file(local_path: str, remote_path: str, copy: bool = False)`
Adds a local file to the image.

### `add_local_dir(local_path: str, remote_path: str, copy: bool = False)`
Adds a local directory's contents to the image.

### `apt_install(*packages: str)`
Install Debian packages using apt.

**Example:**
```python
image = modal.Image.debian_slim().apt_install("git")
```

### `run_commands(*commands: str)`
Extend an image with shell commands.

### `env(vars: dict[str, str])`
Sets environment variables in the image.

**Example:**
```python
image = (
    modal.Image.debian_slim()
    .env({"HF_HUB_ENABLE_HF_TRANSFER": "1"})
)
```

## Additional Features
- Support for private repository installations
- Poetry and micromamba package management
- Cloud registry integrations (GCP, AWS)
- Custom Dockerfile commands
- Environment and working directory configuration