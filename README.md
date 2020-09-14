# Kubernetes Mailbox Manager

_N.B.: This application has been created for educational purposes and has not been fully tested in production yet._

Kubernetes Mailbox Manager is a small application allowing mailboxes to be created in a Kubernetes cluster using [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
Any 'Mailbox' known in namespace will automatically be transformed into a postfix accounts config file.

This file can then be supplied as a ConfigMap to any postfix-running setup such as [tomav/docker-mailserver](https://github.com/tomav/docker-mailserver).

## Mailbox

_Make sure the [Mailbox Custom Resource Definition](kubernetes/crd.yaml) is applied._

```yaml
apiVersion: "k8smailman.mikedeheij.nl/v1"
kind: Mailbox
metadata:
  name: some-client
spec:
  emailAddress: "some@client.local"
  passwordHash: "{SHA512-CRYPT}$6$something"
```

# Usage

Example deployment of application:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubernetes-mailbox-manager
spec:
  selector:
    matchLabels:
      app: mail
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mail
    spec:
    #   imagePullSecrets:                       # Optional
    #   - name: <private registry secret>
      containers:
      - image: docker.pkg.github.com/mdeheij/kubernetes-mailbox-manager/app:latest
        name: controller
```

Example of using a generated configmap in a pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: k8smailman-test-pod
spec:
  containers:
    - name: test-container
      image: k8s.gcr.io/busybox
      command: [ "/bin/sh", "-c", "cat /tmp/docker-mailserver/postfix-accounts.cf" ]
      volumeMounts:
      - name: k8smailman-volume
        mountPath: /tmp/docker-mailserver
  volumes:
    - name: k8smailman-volume
      configMap:
        # Provide the name of the ConfigMap created by the mailbox manager
        name: kubernetes-mailbox-manager
  restartPolicy: Never
```

# Development

This application requires Go version 1.11 or higher. 

## Configuration

`K8SMAILMAN_KUBE_CONFIG` can be set to contain a path to a Kubernetes client configuration file (e.g. `~/.kube/config`).

## Building and usage

```bash
go install ./...
K8SMAILMAN_KUBE_CONFIG="${HOME}/.kube/config" kubernetes-mailbox-manager
```