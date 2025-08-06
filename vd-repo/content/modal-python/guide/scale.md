# Scaling out

Modal makes it easy to scale compute across thousands of containers, automatically managing scaling without manual intervention.

## How does autoscaling work on Modal?

Every Modal Function corresponds to an autoscaling pool of containers. The autoscaler:
- Spins up new containers when no capacity is available
- Spins down containers when resources are idling
- Scales to zero by default when no inputs are being processed

Autoscaling decisions are made quickly to support batch jobs and deployed Apps with variable traffic.

## Configuring autoscaling behavior

Modal provides several configuration settings for the autoscaler:

- `max_containers`: Upper limit on containers for a Function
- `min_containers`: Minimum number of warm containers
- `buffer_containers`: Container buffer to reduce input queuing
- `scaledown_window`: Maximum idle duration before container shutdown

These settings help balance cost and latency.

## Dynamic autoscaler updates

You can update autoscaler settings dynamically using `Function.update_autoscaler()`:

```python
f = modal.Function.from_name("my-app", "f")
f.update_autoscaler(max_containers=100)
```

A common pattern is adjusting pool size based on time of day:

```python
@app.function(schedule=modal.Cron("0 6 * * *", timezone="America/New_York"))
def increase_warm_pool():
    inference_server.update_autoscaler(min_containers=4)
```

## Parallel execution of inputs

Use `Function.map()` to run function calls in parallel:

```python
@app.local_entrypoint()
def main():
    inputs = list(range(100))
    for result in evaluate_model.map(inputs):
        ...
```

## Scaling Limits

Modal enforces these limits per function:
- 2,000 pending inputs
- 25,000 total inputs
- 1,000 concurrent inputs per `.map()` invocation

For async jobs, up to 1 million pending inputs are allowed.

## Additional Features

- Supports asynchronous programming
- GPU acceleration available
- Flexible scaling strategies