package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func main() {
	var kubeContext, namespace, podName, container string

	flag.StringVar(&kubeContext, "context", "", "Kubernetes context")
	flag.StringVar(&namespace, "namespace", "default", "Namespace")
	flag.StringVar(&podName, "pod", "", "Pod name (required)")
	flag.StringVar(&container, "container", "", "Container name (optional)")
	flag.Parse()

	command := flag.Args()
	if len(command) == 0 {
		fmt.Fprintf(os.Stderr, "Error: command is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -pod <pod-name> [-context <context>] [-namespace <namespace>] [-container <container>] <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	if podName == "" {
		fmt.Fprintf(os.Stderr, "Error: pod name is required\n")
		os.Exit(1)
	}

	// Load kubeconfig
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if env := os.Getenv("KUBECONFIG"); env != "" {
		kubeconfig = env
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: kubeContext},
	).ClientConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating clientset: %v\n", err)
		os.Exit(1)
	}

	// Create exec request
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	// Execute
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating executor: %v\n", err)
		os.Exit(1)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	// Print stderr and stdout
	if stderr.Len() > 0 {
		fmt.Fprint(os.Stderr, stderr.String())
	}
	if stdout.Len() > 0 {
		fmt.Print(stdout.String())
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
