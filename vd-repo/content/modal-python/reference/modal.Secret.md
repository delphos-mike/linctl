# modal.Secret

A class for managing secrets in Modal, providing a secure way to handle environment variables and credentials.

## Overview

`Secret` is a class that allows you to securely add sensitive information to containers running your functions.

## Methods

### `hydrate()`
- Synchronizes the local object with its server-side identity
- Rarely needs explicit calling
- Replaces deprecated `.resolve()` method

### `from_dict(env_dict: dict)`
- Creates a secret from a dictionary of environment variables
- Example usage:
  ```python
  @app.function(secrets=[modal.Secret.from_dict({"FOO": "bar"})])
  def run():
      print(os.environ["FOO"])
  ```

### `from_local_environ(env_keys: list)`
- Creates secrets from local environment variables

### `from_dotenv(path=None, filename=".env")`
- Creates secrets from a .env file
- Can specify custom file location and name

### `from_name(name: str, environment_name: Optional[str] = None, required_keys: list = [])`
- References a Secret by its name
- Requires pre-provisioning from Dashboard

### `info()`
- Returns information about the Secret object

## Properties

### `name`
- Returns the optional name of the secret

## Key Features
- Secure credential management
- Flexible secret creation methods
- Integration with Modal's function deployment system