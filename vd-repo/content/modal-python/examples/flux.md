# Run Flux fast on H100s with `torch.compile`

## Overview
This guide demonstrates how to run Flux image generation quickly on NVIDIA H100 GPUs using Modal and `torch.compile` optimization techniques.

## Key Components

### Dependencies and Image Setup
- Uses CUDA toolkit 12.4.0
- Installs dependencies like PyTorch 2.5.0, Diffusers, and other ML libraries
- Configures a custom Docker image with necessary tools

### Optimization Techniques
- Uses `torch.compile()` to optimize compute graphs
- Applies memory layout optimizations
- Caches compilation artifacts to reduce startup time

## Code Structure

### Main Components
- Modal App configuration
- GPU-accelerated inference class
- Parameterized model with compilation option
- Local entrypoint for running inference

### Compilation Optimization Example
```python
def optimize(pipe, compile=True):
    # Fuse QKV projections
    pipe.transformer.fuse_qkv_projections()
    pipe.vae.fuse_qkv_projections()

    # Optimize memory layout
    pipe.transformer.to(memory_format=torch.channels_last)
    pipe.vae.to(memory_format=torch.channels_last)

    if compile:
        # Additional torch compilation configurations
        pipe.transformer = torch.compile(
            pipe.transformer, mode="max-autotune", fullgraph=True
        )
```

## Running the Example
```bash
modal run flux.py --compile
```

## Performance Notes
- First inference can take up to 20 minutes for compilation
- Subsequent inferences are significantly faster
- Achieves image generation in approximately 700-1200 milliseconds

## Additional Resources
- [Hugging Face Diffusers Guide](https://huggingface.co/docs/diffusers/en/tutorials/fast_diffusion)
- [PyTorch Compilation Tutorial](https://pytorch.org/tutorials/recipes/torch_compile_caching_tutorial.html)