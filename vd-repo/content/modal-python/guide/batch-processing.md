# Batch Processing in Modal

Modal is optimized for large-scale batch processing, enabling functions to scale to thousands of parallel containers with minimal configuration.

## Background Execution with `.spawn_map`

The fastest way to submit multiple jobs asynchronously is using `.spawn_map`. Here's an example of submitting 100,000 video embeddings:

```python
import modal

app = modal.App("batch-processing-example")
volume = modal.Volume.from_name("video-embeddings", create_if_missing=True)

@app.function(volumes={"/data": volume})
def embed_video(video_id: int):
    # Business logic:
    # - Load the video from the volume
    # - Embed the video
    # - Save the embedding to the volume
    ...

@app.local_entrypoint()
def main():
    embed_video.spawn_map(range(100_000))
```

This approach works best for jobs that store results externally, such as in a Modal Volume or cloud bucket.

## Parallel Processing with `.map`

`.map` allows offloading expensive computations to powerful machines while gathering results. Example of parallel video similarity queries:

```python
import modal
import itertools

app = modal.App("gather-results-example")

@app.function(gpu="L40S")
def compute_video_similarity(query: str, video_id: int) -> tuple[int, int]:
    # Embed video with GPU acceleration & compute similarity with query
    return video_id, score

@app.local_entrypoint()
def main():
    queries = itertools.repeat("Modal for batch processing")
    video_ids = range(100_000)

    for video_id, score in compute_video_similarity.map(queries, video_ids):
        # Process results (e.g., extract top 5 most similar videos)
        pass
```

## Integration with Existing Systems

The recommended method is using [deployed function invocation](/docs/guide/trigger-deployed-functions):

```python
def external_function(inputs):
    compute_similarity = modal.Function.from_name(
        "gather-results-example",
        "compute_video_similarity"
    )
    results = compute_similarity.remote_gen(inputs)
    return list(results)
```