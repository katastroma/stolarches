# Stolarches

Katastroma's orderer. Implements the
[diataxis](https://github.com/katastroma/diataxis) client API.

Given manifests, stolarches ensures they are in a safe apply order.

## Communication

The renderer streams a YAML manifest blob via the diataxis `Order` RPC
(client-streaming gRPC). The orderer type can be passed as gRPC metadata — if
absent, it defaults to helm.

After receiving the full stream, stolarches responds to the renderer
immediately. Ordering and forwarding to the provisioner happen asynchronously —
the renderer does not wait for ordering to complete. This decouples the
renderer's lifecycle from ordering time.

Stolarches opens a client-streaming connection to the provisioner's katartismos
`Provision` RPC and streams the ordered resources as JSON. The provisioner
address is configured via `PROVISIONER_ADDR`.

## Backend Dispatch

Stolarches registers ordering backends at startup. The orderer type from gRPC
metadata selects which backend processes the request.

Each backend implements `Receive` (parse the YAML manifest stream) and `Order`
(sort manifests into a safe apply sequence). Adding a new ordering backend means
implementing that interface and registering it.

## Ordering Backends

- **Helm** — uses `releaseutil.SortManifests` with `InstallOrder`. Battle-tested
  and widely understood. Carries helm as a transitive dependency.
- **CLI-Utils** — uses `sigs.k8s.io/cli-utils/pkg/ordering`, sorts by GVK →
  namespace → name. No helm dependency.

Both produce a sequence where prerequisites (Namespaces, CRDs, RBAC) appear
before the resources that depend on them.

## Assumptions

- Manifests arrive as a valid YAML blob. Invalid YAML fails at the backend's
  receive/decode step.
- The orderer type in gRPC metadata, if present, matches a registered backend.
  Unknown types fall back to helm.
- The provisioner is reachable at `PROVISIONER_ADDR`. Forwarding failures are
  logged but not retried — pharos handles failure recovery via replay at the
  source handler.

## Design Direction

**Per-source-target orderer selection.** The orderer type is currently
determined by gRPC metadata from the renderer (defaulting to helm). In the
future, the orderer should be configurable at the source target (the ConfigMap),
so tenants can specify which ordering strategy applies to their source.
