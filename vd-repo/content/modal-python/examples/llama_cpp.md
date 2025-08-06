# Run Large and Small Language Models with llama.cpp (DeepSeek-R1, Phi-4)

## Overview

This Modal example demonstrates how to run small (Phi-4) and large (DeepSeek-R1) language models using [`llama.cpp`](https://github.com/ggerganov/llama.cpp).

## GPU Requirements

### DeepSeek-R1
- 671B total parameters
- Requires over 100GB storage
- Needs four L40S GPUs (192 GB memory)

### Phi-4
- 14B total parameters
- Approximately 5 GB when quantized
- Can run comfortably on a CPU

## Key Features

- Command-line inference triggering
- Model weight storage using Modal Volumes
- CUDA-compiled llama.cpp
- Flexible model and prompt selection

## Example Usage

Run default inference:
```bash
modal run llama_cpp.py
```

Run Phi-4 model:
```bash
modal run llama_cpp.py --model="phi-4"
```

Custom prompt and token limit:
```bash
modal run llama_cpp.py --prompt="Your custom prompt" --n-predict=1024
```

## Technical Components

- Modal App configuration
- Model download function
- CUDA-enabled container image
- Inference function with GPU/CPU support
- Output collection and storage

## Default Configurations

### DeepSeek-R1 Arguments
```python
DEFAULT_DEEPSEEK_R1_ARGS = [
    "--cache-type-k", "q4_0",
    "--threads", "12",
    "-no-cnv", "--prio", "2",
    "--temp", "0.6",
    "--ctx-size", "8192"
]
```

### Phi-4 Arguments
```python
DEFAULT_PHI_ARGS = [
    "--threads", "16",
    "-no-cnv",
    "--ctx-size", "16384"
]
```

## Getting Started

1. Install Modal: `pip install modal`
2. Set up Modal: `modal setup`
3. Clone examples repository
4. Run the example