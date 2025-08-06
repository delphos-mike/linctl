# modal.App

A Modal App is a group of functions and classes deployed together, serving three primary purposes:
- Unit of deployment
- Syncing identities across processes
- Managing log collection

## Class Definition

```python
class App(object):
    def __init__(
        self,
        name: Optional[str] = None,
        image: Optional[Image] = None,
        secrets: Sequence[Secret] = [],
        volumes: dict[Union[str, PurePosixPath], Volume] = {},
        include_source: bool = True
    )
```

## Key Properties

- `name`: User-provided name of the App
- `is_interactive`: Whether the app is running in interactive mode
- `app_id`: Unique identifier for the running or stopped app
- `description`: App's name or fallback identifier

## Key Methods

### `run()`
Context manager to run an ephemeral app on Modal
```python
with app.run():
    some_modal_function.remote()
```

### `deploy()`
Deploy the App persistently
```python
app.deploy(name="my-app")
```

### `function()`
Decorator to register a new Modal Function
```python
@app.function(
    secrets=[modal.Secret.from_name("some_secret")],
    schedule=modal.Period(days=1)
)
def foo():
    pass
```

### `local_entrypoint()`
Decorator for CLI entrypoint functions
```python
@app.local_entrypoint()
def main():
    some_modal_function.remote()
```

### `include()`
Include another App's objects
```python
app_a.include(app_b)
```

## Registered Components

- `registered_functions`: Modal Function objects
- `registered_classes`: Modal Cls objects
- `registered_entrypoints`: Local CLI entrypoints
- `registered_web_endpoints`: Webhook function names