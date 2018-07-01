## Running with Minikube
For this, you don't need any tools other than Docker and minikube.

## Running locally
Running locally is useful for testing REST APIs and gRPC interaction.
However, by default, we assume RabbitMQ is not running locally.
That is, by default, services won't publish and read messages.

Before you get started, you'll need to install some tools.

### Prerequisites
#### Postgres
All our services use Postgres. I have a script to do this in `./scripts`.

#### Golang
Install golang.

#### Protocol Buffers
You need to be able to compile protocol buffers with `protoc`: 
```bash
make install-proto
```

### Building and Running
To build all Go files, run: 
```bash
make
```

To run your services, run each Make target in a different tab:
```bash
make run-auth
make run-password-reset
```

To both build and run a specific service, you can run:
```bash
make auth
make password-reset
```
