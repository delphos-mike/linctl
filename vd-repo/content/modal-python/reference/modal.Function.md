# modal.Function

A class representing serverless functions in Modal, typically created using the `App.function()` decorator.

## Methods

### hydrate
```python
def hydrate(self, client: Optional[_Client] = None) -> Self
```
Synchronizes the local object with its server-side identity. Rarely needs explicit calling.

### update_autoscaler
```python
def update_autoscaler(
    self,
    *,
    min_containers: Optional[int] = None,
    max_containers: Optional[int] = None,
    buffer_containers: Optional[int] = None,
    scaledown_window: Optional[int] = None,
) -> None
```
Overrides current autoscaler behavior for the function.

### from_name
```python
@classmethod
def from_name(
    cls,
    app_name: str,
    name: str,
    *,
    environment_name: Optional[str] = None,
) -> "_Function"
```
References a Function from a deployed App by its name.

### get_web_url
```python
def get_web_url(self) -> Optional[str]
```
Returns the URL of a Function running as a web endpoint.

### remote
```python
def remote(self, *args: P.args, **kwargs: P.kwargs) -> ReturnType
```
Calls the function remotely with given arguments.

### local
```python
def local(self, *args: P.args, **kwargs: P.kwargs) -> OriginalReturnType
```
Calls the function locally in the current environment.

### spawn
```python
def spawn(self, *args: P.args, **kwargs: P.kwargs) -> "_FunctionCall[ReturnType]"
```
Calls the function without waiting for results, returning a `FunctionCall` object.

### map
```python
def map(
    self,
    *input_iterators,
    kwargs={},
    order_outputs: bool = True,
    return_exceptions: bool = False
) -> AsyncOrSyncIterable
```
Performs parallel mapping over input iterators.

### starmap
```python
def starmap(
    self,
    input_iterator,
    kwargs={},
    order_outputs: bool = True,
    return_exceptions: bool = False
) -> AsyncOrSyncIterable
```
Like `map` but unpacks arguments from tuples.