# Hello, world!

This tutorial demonstrates some core features of Modal:

- You can run functions on Modal just as easily as you run them locally
- Running functions in parallel on Modal is simple and fast
- Logs and errors show up immediately, even for functions running on Modal

## Importing Modal and setting up

Start by importing `modal` and creating an `App`:

```python
import sys
import modal

app = modal.App("example-hello-world")
```

## Defining a function

Create a function that takes an input, prints logs, and returns an output:

```python
@app.function()
def f(i):
    if i % 2 == 0:
        print("hello", i)
    else:
        print("world", i, file=sys.stderr)
    
    return i * i
```

## Running our function locally, remotely, and in parallel

Demonstrate three ways to call the function:

```python
@app.local_entrypoint()
def main():
    # run the function locally
    print(f.local(1000))
    
    # run the function remotely on Modal
    print(f.remote(1000))
    
    # run the function in parallel and remotely on Modal
    total = 0
    for ret in f.map(range(200)):
        total += ret
    
    print(total)
```

## What just happened?

When calling `.remote` on `f`, the function is executed in the cloud on Modal's infrastructure.

## Why does this matter?

### Flexibility and Ease of Use

- Change code and run it again instantly
- Map over large datasets easily
- Run complex functions like:
  - Language model inference
  - Audio or image manipulation
  - Large-scale text embedding

## Getting Started

To run this example:

1. Install Modal:
```
pip install modal
```

2. Set up Modal:
```
modal setup
```

3. Clone the repository:
```
git clone https://github.com/modal-labs/modal-examples
cd modal-examples
modal run 01_getting_started/hello_world.py
```