# Stolarches

Katastroma's orderer. Implements the
[diataxis](https://github.com/katastroma/diataxis) interface.

Given manifests, stolarches ensures they are in a safe apply order.

## Implementation

Candidate libraries for the ordering logic:

- `helm.sh/helm/v3/pkg/releaseutil` — Helm's `InstallOrder`, battle-tested and
  widely understood. Carries helm as a transitive dependency.
- `sigs.k8s.io/cli-utils/pkg/ordering` — from kubernetes-sigs, sorts by
  GVK → namespace → name. No helm dependency.
