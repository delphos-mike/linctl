# Introduction to Modal

Modal is a cloud function platform that enables developers to:

- Run code remotely within seconds
- Define container environments in code
- Scale horizontally to thousands of containers
- Attach GPUs easily
- Serve functions as web endpoints
- Deploy scheduled jobs
- Store data in distributed dictionaries and queues

## Getting Started

To begin using Modal:

1. Create an account at [modal.com](https://modal.com)
2. Install the Modal Python package: `pip install modal`
3. Authenticate: `modal setup` (or `python -m modal setup`)

### Quick Examples

- [Hello, world!](/docs/examples/hello_world)
- [A simple web scraper](/docs/examples/web-scraper)

## How Does It Work?

Modal takes your code and executes it in a cloud container environment. Key benefits include:

- Solving infrastructure challenges
- No need to manage Kubernetes, Docker, or cloud accounts
- Currently Python-only, with potential future language support

### Key Features

- Full serverless execution
- Per-second pricing
- Zero configuration
- Code-first approach

**Note**: Everything is defined directly in code, eliminating the need for complex configuration files like YAML.

You can also explore Modal interactively through their [code playground](/playground).