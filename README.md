# cfssl-sidecar

A sidecar utility for managing certificates using a remote cfssl signing server.

## usage

Grab the release from github releases, or use the docker release `roboll/cfssl-sidecar`. The sidecar will generate a key and csr and boot, and will immediately request a signed cert from the specified remote. It will then wake every interval and request a new certificate based on the existing csr.

It accepts the following options:
```
  -certname string
    	name of resulting cert and key (${certname}.pem, ${certname}-key.pem)
  -certpath string
    	path of resulting cert and key
  -csrfile string
    	path to csr json file
  -hostname string
    	hostname for the cert, comma separated
  -interval duration
    	repeat interval
  -label string
    	signing config label
  -perms uint
    	permissions for resulting dirs/files (default 384)
  -profile string
    	signing config profile
  -remote string
    	remote cfssl server
```

## examples

Example usage as a [kubernetes pod sidecar container](examples/kube).
