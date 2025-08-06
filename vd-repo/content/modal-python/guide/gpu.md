# GPU Acceleration

Modal makes it easy to run your code on GPUs.

## Quickstart

Here's a simple example of a Function running on an A100 in Modal:

```python
import modal

image = modal.Image.debian_slim().pip_install("torch")
app = modal.App(image=image)

@app.function(gpu="A100")
def run():
    import torch

    assert torch.cuda.is_available()
```

## Specifying GPU Type

Modal supports the following GPU types:
- T4
- L4
- A10
- A100
- A100-40GB
- A100-80GB
- L40S
- H100/H100!
- H200
- B200

Example: `@app.function(gpu="B200")`

## Specifying GPU Count

You can specify multiple GPUs per container:

```python
@app.function(gpu="H100:8")
def run_llama_405b_fp8():
    ...
```

Note: B200, H200, H100, A100, L4, T4, and L40S support up to 8 GPUs.

## Picking a GPU

For neural network inference, the L40S is recommended, offering 48 GB of GPU RAM with a good cost-performance balance.

## B200 GPUs

B200s are NVIDIA's flagship data center chips for the Blackwell architecture:

```python
@app.function(gpu="B200:8")
def run_deepseek():
    ...
```

## H200 and H100 GPUs

H200s and H100s are top-of-the-line data center chips from NVIDIA based on the Hopper architecture.

### Automatic Upgrades to H200s

Modal may automatically upgrade H100 requests to H200 at no additional cost.

## A100 GPUs

A100s come in 40 GB and 80 GB RAM versions:

```python
@app.function(gpu="A100")  # May upgrade to 80 GB
def qwen_7b():
    ...

@app.function(gpu="A100-80GB")
def llama_70b():
    ...
```