# tokenbridge-k8s-sidecar

`tokenbridge-k8s-sidecar` is a Kubernetes sidecar container that exchanges a projected Kubernetes ServiceAccount token for a custom access token using a remote tokenbridge service. The sidecar writes the resulting token to a shared volume, making it available to your main application container.

## Features

- Securely exchanges Kubernetes ServiceAccount tokens for custom tokens.
- Periodically refreshes the token before expiration.
- Shares the token with your main application via a writable volume.
- Designed for use as a sidecar in Kubernetes Pods.

## Usage

### 1. Use the Public Docker Image

A prebuilt image is available on Docker Hub:

```sh
docker pull hupe1980/tokenbridge-k8s-sidecar:latest
```

You can reference this image directly in your Kubernetes manifests.

### 2. Example Kubernetes Manifest

See [`config/deployment.yaml`](config/deployment.yaml) for an example manifest:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tokenbridge-sa
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tokenbridge-example
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tokenbridge-example
  template:
    metadata:
      labels:
        app: tokenbridge-example
    spec:
      serviceAccountName: tokenbridge-sa
      volumes:
        - name: token-vol
          projected:
            sources:
              - serviceAccountToken:
                  path: sa-token
                  expirationSeconds: 3600
                  audience: tokenbridge
        - name: tokenbridge-vol
          emptyDir: {}
      containers:
        - name: main-app
          image: alpine:3.21
          command: ["/bin/sh"]
          args: ["-c", "sleep infinity"]
          volumeMounts:
            - name: tokenbridge-vol
              mountPath: /run/secrets/tokenbridge
              readOnly: true
          env:
            - name: TOKENBRIDGE_TOKEN_FILE
              value: /run/secrets/tokenbridge/access-token

        - name: sidecar
          image: hupe1980/tokenbridge-k8s-sidecar:latest
          volumeMounts:
            - name: token-vol
              mountPath: /var/run/secrets/tokens
              readOnly: true
            - name: tokenbridge-vol
              mountPath: /run/secrets/tokenbridge
          env:
            - name: SA_TOKEN_PATH
              value: /var/run/secrets/tokens/sa-token
            - name: OUTPUT_TOKEN_PATH
              value: /run/secrets/tokenbridge/access-token
            - name: EXCHANGE_URL
              value: https://your-tokenbridge-service/exchange
            - name: REFRESH_INTERVAL
              value: 1h
            - name: AUDIENCE
              value: tokenbridge
```

### 3. Environment Variables

The sidecar expects the following environment variables:

- `SA_TOKEN_PATH` (required): Path to the projected ServiceAccount token.
- `OUTPUT_TOKEN_PATH` (required): Path to write the exchanged token.
- `EXCHANGE_URL` (required): URL of the tokenbridge service.
- `REFRESH_INTERVAL` (optional): How often to refresh the token (default: `1h`).
- `AUDIENCE` (optional): Audience for the token exchange.


## Development

- Build: `make build`
- Lint: `make lint`
- Test: `make test`

## Related Projects

- [**TokenBridge**](https://github.com/hupe1980/tokenbridge): The main project for TokenBridge, providing core functionality and documentation.
- [**TokenBridge GitHub Action**](https://github.com/hupe1980/tokenbridge-action): Automate your workflows with TokenBridge using GitHub Actions.
- [**TokenBridge Backend Example**](https://github.com/hupe1980/tokenbridge-backend-example): A practical example of how to create a TokenBridge backend application.

## License

MIT License. See [LICENSE](LICENSE) for details.