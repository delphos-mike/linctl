# Run OpenAI-compatible LLM Inference with LLaMA 3.1-8B and vLLM

## Overview

This example demonstrates how to run a vLLM server in OpenAI-compatible mode on Modal, using Meta's LLaMA 3.1 8B Instruct model. The guide covers setting up a container, downloading model weights, configuring vLLM, and deploying a serverless LLM inference endpoint.

## Key Components

### Container Image Setup
```python
vllm_image = (
    modal.Image.debian_slim(python_version="3.12")
    .pip_install(
        "vllm==0.9.1",
        "huggingface_hub[hf_transfer]==0.32.0",
        "flashinfer-python==0.2.6.post1",
        extra_index_url="https://download.pytorch.org/whl/cu128",
    )
    .env({"HF_HUB_ENABLE_HF_TRANSFER": "1"})
)
```

### Model Configuration
- Model: "RedHatAI/Meta-Llama-3.1-8B-Instruct-FP8"
- Uses 8-bit floating point quantization
- Cached using Modal Volumes

### vLLM Server Deployment
```python
@app.function(
    image=vllm_image,
    gpu=f"B200:{N_GPU}",
    volumes={
        "/root/.cache/huggingface": hf_cache_vol,
        "/root/.cache/vllm": vllm_cache_vol,
    },
)
@modal.web_server(port=VLLM_PORT)
def serve():
    # Launches vLLM server with OpenAI-compatible endpoints
```

## Deployment and Interaction

1. Deploy server: `modal deploy vllm_inference.py`
2. Interact via Swagger UI or OpenAI library
3. Test with included `client.py` script

## Performance Considerations

- Configurable boot and compilation strategies
- Supports dynamic JIT compilation for optimization