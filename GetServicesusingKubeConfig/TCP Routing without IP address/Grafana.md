# Traefik Configuration for Exposing a Service on Port 3000

This guide will walk you through configuring Traefik to expose a service (e.g., Grafana) on port 3000. We'll set up an entry point, configure the service, and route traffic using an `IngressRoute`.

## Step 1: Add Traefik Entry Point for Port 3000

You need to define a custom entry point in the Traefik configuration for port 3000. Add the following configuration to your Traefik deployment YAML under `additionalArguments`:

```yaml
additionalArguments:
  - "--entryPoints.custom3000.address=:3000" # for Grafana
```

This command creates a new entry point `custom3000` that listens on port 3000.

## Step 2: Configure the Service with the Desired Port

Now, specify the port for your service in Traefik, making sure it's exposed and accessible. Add the following configuration under the service specification:

```yaml
custom3000:
  port: 3000
  expose:
    default: true
  # The exposed port for this service
  exposedPort: 3000
  # The port protocol (TCP/UDP)
  protocol: TCP
```

This will expose port 3000 for the service using the TCP protocol.

## Step 3: Create a Traefik `IngressRoute` to Route Traffic to Port 3000

Now, define the routing rules to direct traffic on port 3000 to the correct service (e.g., Grafana). Add the following `IngressRoute` configuration:

```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: grafana-ingress-3000
  namespace: default
spec:
  entryPoints:
    - custom3000 # Use the custom3000 entry point to route traffic on port 3000
  routes:
    - match: HostRegexp('.*') # Replace with your public IP or domain
      kind: Rule
      services:
        - name: grafana
          port: 3000 # Internally route to Grafana service on port 3000
```

### Key Components:

- **entryPoints.custom3000:** This ensures traffic on port 3000 is routed using the custom entry point.
- **match:** Replace `67.207.68.98` with your public IP or domain name to define the host that will route traffic.
- **services:** Define the service name (`grafana` in this case) and ensure the correct internal port (`3000`) is specified.

---

With these configurations, Traefik will expose your Grafana service on port 3000, route traffic from a specific host, and allow external access.
