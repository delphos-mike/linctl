# API Reference

This is the API reference for the [`modal`](https://pypi.org/project/modal/) Python package, which allows you to run distributed applications on Modal.

## Application Construction

- [`App`](/docs/reference/modal.App): The main unit of deployment for code on Modal
- [`App.function`](/docs/reference/modal.App#function): Decorator for registering a function with an App
- [`App.cls`](/docs/reference/modal.App#cls): Decorator for registering a class with an App

## Serverless Execution

- [`Function`](/docs/reference/modal.Function): A serverless function backed by an autoscaling container pool
- [`Cls`](/docs/reference/modal.Cls): A serverless class supporting parametrization and lifecycle hooks

## Extended Function Configuration

### Class Parametrization
- [`parameter`](/docs/reference/modal.parameter): Used to define class parameters, akin to a Dataclass field

### Lifecycle Hooks
- [`enter`](/docs/reference/modal.enter): Decorator for method executed during container startup
- [`exit`](/docs/reference/modal.exit): Decorator for method executed during container shutdown
- [`method`](/docs/reference/modal.method): Decorator for exposing a method as an invokable function

### Web Integrations
- [`fastapi_endpoint`](/docs/reference/modal.fastapi_endpoint): Decorator for exposing a simple FastAPI-based endpoint
- [`asgi_app`](/docs/reference/modal.asgi_app): Decorator for constructing an ASGI web application
- [`wsgi_app`](/docs/reference/modal.wsgi_app): Decorator for constructing a WSGI web application
- [`web_server`](/docs/reference/modal.web_server): Decorator for constructing an HTTP web server

### Function Semantics
- [`batched`](/docs/reference/modal.batched): Decorator enabling dynamic input batching
- [`concurrent`](/docs/reference/modal.concurrent): Decorator enabling input concurrency

### Scheduling
- [`Cron`](/docs/reference/modal.Cron): Schedule running based on cron syntax
- [`Period`](/docs/reference/modal.Period): Schedule running based on time intervals

## Sandboxed Execution
- [`Sandbox`](/docs/reference/modal.Sandbox): Create and manage sandboxed execution environments
- [`ContainerProcess`](/docs/reference/modal.ContainerProcess): Process management in containers
- [`FileIO`](/docs/reference/modal.FileIO): File I/O operations in sandboxes

## Container Configuration
- [`Image`](/docs/reference/modal.Image): Container image definition and customization
- [`Secret`](/docs/reference/modal.Secret): Secure credential management

## Data Primitives
### Persistent Storage
- [`Volume`](/docs/reference/modal.Volume): Distributed file system
- [`CloudBucketMount`](/docs/reference/modal.CloudBucketMount): Cloud storage mounting
- [`NetworkFileSystem`](/docs/reference/modal.NetworkFileSystem): Network-attached storage

### In-memory Storage
- [`Dict`](/docs/reference/modal.Dict): Distributed key-value store
- [`Queue`](/docs/reference/modal.Queue): Distributed task queue

## Networking
- [`Proxy`](/docs/reference/modal.Proxy): Network proxy configuration
- [`forward`](/docs/reference/modal.forward): Port forwarding utilities