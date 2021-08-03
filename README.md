# vault-k8s-unsealer

A simple utility to automate the unsealing of a Vault instance.

To run this application simply provide two environment variables:
* `VAULT_ADDR` with the address of the Vault instance you would like to keep unsealed.
* `UNSEAL_KEY` with hte name of the file to read in the secret key to provide.

Ideally the `UNSEAL_KEY` would be contained a `secret` mount within a Kubernetes pod but totally not required.

## Docker

Images are published under `meschbach/vault-k8s-unsealer` with the SCM tag corresponding to the image tag.
