# Scheduling Remote Cron Jobs in Modal

## Basic Scheduling

To schedule a function to run periodically, create a Modal App and use the `@app.function` decorator with a schedule parameter:

```python
# heavy.py
import modal

app = modal.App()

@app.function(schedule=modal.Period(days=1))
def perform_heavy_computation():
    ...

if __name__ == "__main__":
   app.deploy()
```

Deploy the app using the CLI:
```
modal deploy --name daily_heavy heavy.py
```

## Monitoring Scheduled Runs

- View execution logs in the [Apps](https://modal.com/apps) section
- Schedules cannot be paused
- Manually start schedules using the "run now" button

## Schedule Types

### Period Scheduling
Specify intervals between function calls:

```python
# runs once every 5 hours
@app.function(schedule=modal.Period(hours=5))
def perform_heavy_computation():
    ...
```

### Cron Scheduling
Use cron syntax for precise scheduling:

```python
# runs at 8 am (UTC) every Monday
@app.function(schedule=modal.Cron("0 8 * * 1"))
def perform_heavy_computation():
    ...

# runs daily at 6 am (New York time)
@app.function(schedule=modal.Cron("0 6 * * *", timezone="America/New_York"))
def send_morning_report():
    ...
```

**Note:** When redeploying with `modal.Period`, the schedule resets. For consistent timing, use `modal.Cron`.