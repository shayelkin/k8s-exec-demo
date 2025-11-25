# k8s-exec-demo

A minimal Go application that replicates `kubectl exec` functionality without shelling out to the kubectl executable.

## Overview

This tool uses the Kubernetes client-go library to execute commands in pods. Given a Kubernetes context, namespace, and pod name, it executes a command and returns its output.

## Features

- Direct Kubernetes API communication (no kubectl subprocess)
- Supports custom kubeconfig contexts
- Configurable namespace and pod
- Optional container specification
- Captures both stdout and stderr

## Building

```bash
cd k8s-exec-demo
go mod download
go build -o k8s-exec-demo
```

## Usage

```bash
# Basic usage
./k8s-exec-demo -pod <pod-name> <command> [args...]

# With namespace
./k8s-exec-demo -namespace <namespace> -pod <pod-name> <command> [args...]

# With context
./k8s-exec-demo -context <context> -namespace <namespace> -pod <pod-name> <command> [args...]

# With specific container
./k8s-exec-demo -pod <pod-name> -container <container-name> <command> [args...]
```

## Examples

```bash
# Run ls in a pod
./k8s-exec-demo -pod my-pod ls /app

# Execute in specific namespace
./k8s-exec-demo -namespace production -pod my-pod env

# Use specific context and container
./k8s-exec-demo -context staging -pod my-pod -container main-container ps aux
```

## Flags

- `-context`: Kubernetes context to use (defaults to current context)
- `-namespace`: Namespace of the pod (defaults to "default")
- `-pod`: Pod name (required)
- `-container`: Container name within the pod (optional, uses first container if not specified)

## How it works

The tool:

1. Loads kubeconfig from `~/.kube/config` or `$KUBECONFIG`
2. Creates a Kubernetes clientset using client-go
3. Uses the Pod Exec API with SPDY executor to run commands
4. Streams stdout/stderr from the pod
5. Returns the output and exits with appropriate status code

## Dependencies

- k8s.io/client-go: Kubernetes Go client
- k8s.io/api: Kubernetes API types

No external kubectl binary required.
