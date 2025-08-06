# Apps, Functions, and entrypoints

An `App` is the object that represents an application running on Modal. All functions and classes are associated with an `App`.

When you `run` or `deploy` an `App`, it creates an ephemeral or a deployed `App`, respectively.

## Ephemeral Apps

An ephemeral App is created when you use the `modal run` CLI command, or the `app.run` method. This creates a temporary App that only exists for the duration of your script.

Ephemeral Apps are stopped automatically when the calling program exits, or when the server detects that the client is no longer connected. You can use `--detach` to keep an ephemeral App running even after the client exits.

Example of running an App:

```python
def main():
    with app.run():
        some_modal_function.remote()
```

To enable output, use `modal.enable_output()`:

```python
def main():
    with modal.enable_output():
        with app.run():
            some_modal_function.remote()
```

## Deployed Apps

A deployed App is created using the `modal deploy` CLI command. The App is persisted indefinitely until you delete it via the web UI. Functions in a deployed App with an attached schedule will run on that schedule. Otherwise, you can invoke them manually using web endpoints or Python.

Deployed Apps are named via the `App` constructor. Re-deploying an existing `App` will update it in place.

## Entrypoints for Ephemeral Apps

The code that runs first when you `modal run` an App is called the "entrypoint".

You can register a local entrypoint using the `@app.local_entrypoint()` decorator.

### Argument Parsing

If your entrypoint function takes arguments with primitive types, `modal run` automatically parses them as CLI options:

```python
@app.local_entrypoint()
def main(foo: int, bar: str):
    some_modal_function.remote(foo, bar)
```

This can be called with `modal run script.py --foo 1 --bar "hello"`.

For custom argument parsing, you can use libraries like `argparse`:

```python
import argparse

@app.local_entrypoint()
def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--foo", type=int)
    args = parser.parse_args()
    some_modal_function.remote(args.foo)
```