# Secrets

Securely provide credentials and other sensitive information to your Modal Functions with Secrets.

## Deploy Secrets from the Modal Dashboard

You can create and edit Secrets via:
- The [dashboard](/secrets)
- Command line interface (`modal secret`)
- Programmatically from Python code (`modal.Secret`)

## Use Secrets in your Modal Apps

Add secrets to a function by using the `secrets=[...]` argument:

```python
@app.function(secrets=[modal.Secret.from_name("secret-keys")])
def some_function():
    import os
    
    secret_key = os.environ["MY_PASSWORD"]
    ...
```

You can inject multiple secrets:

```python
@app.function(secrets=[
    modal.Secret.from_name("my-secret-name"),
    modal.Secret.from_name("other-secret"),
])
def other_function():
    ...
```

## Create Secrets Programmatically

Create secrets directly in your script:

```python
import os

if modal.is_local():
    local_secret = modal.Secret.from_dict({"FOO": os.environ["LOCAL_FOO"]})
else:
    local_secret = modal.Secret.from_dict({})

@app.function(secrets=[local_secret])
def some_function():
    import os
    print(os.environ["FOO"])
```

You can also use `Secret.from_dotenv()` with a `.env` file:

```python
@app.function(secrets=[modal.Secret.from_dotenv()])
def some_other_function():
    print(os.environ["USERNAME"])
```

## Interact with Secrets from the Command Line

View secrets:
```
modal secret list
```

Create a new secret:
```
modal secret create database-secret PGHOST=uri PGPORT=5432 PGUSER=admin PGPASSWORD=hunter2
```

Remove a secret:
```
modal secret delete database-secret
```