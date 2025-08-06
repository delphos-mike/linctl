# Create your own music samples with MusicGen

## Overview
MusicGen is an open-source music generation model from Meta that allows users to create custom music samples using AI. This example demonstrates how to run MusicGen models on Modal GPUs with a Gradio UI.

## Key Components

### Dependencies
The example uses:
- Debian Slim Python 3.11
- Hugging Face Hub
- Audiocraft library
- Torch
- FastAPI
- Gradio

### Main Features
- Generate music from text prompts
- Supports variable duration clips
- Web UI for interactive music generation
- Runs on NVIDIA L40S GPU

## Code Highlights

### Model Setup
```python
image = (
    modal.Image.debian_slim(python_version="3.11")
    .apt_install("git", "ffmpeg")
    .pip_install(
        "huggingface_hub[hf_transfer]==0.27.1",
        "torch==2.1.0",
        "numpy<2",
        "git+https://github.com/facebookresearch/audiocraft.git@v1.3.0"
    )
)
```

### Music Generation Method
The `generate` method creates music clips by:
- Supporting prompts up to 30 seconds
- Handling longer durations through auto-regressive generation
- Allowing custom duration and overlap

### Web UI
A Gradio interface lets users:
- Input text prompts
- Set music duration
- Choose audio format
- Generate and play music samples

## Deployment
Users can deploy the app using:
```bash
modal deploy musicgen.py
```

## Running the Example
```bash
modal run musicgen.py --prompt="Baroque boy band, Bachstreet Boys, basso continuo" --duration=60
```

This example showcases an accessible, powerful approach to AI-generated music using Modal's cloud infrastructure.