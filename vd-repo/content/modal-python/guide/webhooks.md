# Web Endpoints in Modal

Modal provides multiple ways to expose Python functions as web endpoints for different use cases.

## Simple Endpoints

The easiest method is using the `@modal.fastapi_endpoint` decorator:

```python
image = modal.Image.debian_slim().pip_install("fastapi[standard]")

@app.function(image=image)
@modal.fastapi_endpoint()
def f():
    return "Hello world!"
```

### Developing with `modal serve`

Run your script with:

```bash
modal serve server_script.py
```

This creates a temporary public URL for testing your endpoint.

### Deploying with `modal deploy`

Use `modal deploy` to create a persistent web endpoint in the cloud.

### Passing Arguments to an Endpoint

You can pass query parameters or JSON bodies:

```python
@app.function(image=image)
@modal.fastapi_endpoint(method="POST")
def square(item: dict):
    return {"square": item['x']**2}
```

## How Web Endpoints Run in the Cloud

- Endpoints only run when needed
- First request triggers container startup
- Modal keeps containers alive briefly for subsequent requests
- Can create multiple parallel containers for high traffic

## Serving ASGI and WSGI Apps

### ASGI Apps (FastAPI, Starlette)

```python
@app.function(image=image)
@modal.concurrent(max_inputs=100)
@modal.asgi_app()
def fastapi_app():
    from fastapi import FastAPI

    web_app = FastAPI()

    @web_app.post("/echo")
    async def echo(request: Request):
        body = await request.json()
        return body

    return web_app
```

### WSGI Apps (Django, Flask)

```python
@app.function(image=image)
@modal.wsgi_app()
def flask_app():
    from flask import Flask, request

    web_app = Flask(__name__)

    @web_app.post("/echo")
    def echo():
        return request.json

    return web_app
```

## Non-ASGI Web Servers

For servers without ASGI/WSGI support, use the `@modal.web_server` decorator to expose raw network ports.