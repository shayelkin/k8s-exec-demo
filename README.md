# k8s-exec-demo

Replicates `kubectl exec` functionality without shelling out to the kubectl executable.

## Overview

This tool uses the Kubernetes client-go library to execute commands in pods. Given a Kubernetes
context, namespace, and pod name, it executes a command and returns its output.

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

## Command line flags

- `-context`: Kubernetes context to use (defaults to current context)
- `-namespace`: Namespace of the pod (defaults to "default")
- `-pod`: Pod name (required)
- `-container`: Container name within the pod (optional, uses first container if not specified)

## Usage examples

```bash
# Run ls in a pod
./k8s-exec-demo -pod my-pod ls /app

# Execute in specific namespace
./k8s-exec-demo -namespace production -pod my-pod env

# Use specific context and container
./k8s-exec-demo -context staging -pod my-pod -container main-container ps aux
```

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

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
