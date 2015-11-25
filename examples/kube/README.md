# kube example

A sample kubernetes pod with a cfssl sidecar that continuously rotates certificates according to an interval. The items in secret/ must be added as a secret to the cluster with the appropriate credentials and urls.

Secrets could be created separately for each signing profile, or for each consumer. Alternatively, csr could be mounted as a hostpath or git repo - it isn't really secret.
