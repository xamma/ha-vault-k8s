# Setup an HA Vault on K8s
This Repo demonstrates how to setup an highly-availabe (HA) **HashiCorp Vault**.  
Remember to set TLS correctly if you want to use this in an Production-Environment.  
Also you should use auto-unseal with e.g AWS KMS (Create KMS Key, create K8s secret).  

## Configure Vault Environment
You need to configure your Vault with the ```override-values.yaml```.  
The provided file in this Repo sets up an HA vault cluster with 3 replicas but withoud TLS.  
See ```override-values-alloptions.yaml``` to see what specifications can be done or go to the official documentation [HashiCorp Vault](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide).  

## Install Vault via Helm Repo
This installs the vault Helm-Repo and overrides the default values with the configuration
you did before.

```
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm install vault hashicorp/vault \
    --namespace vault \
    --create-namespace \
    -f override-values.yaml
```

## Unseal Vault
Because we did not activate auto-unseal we have to manually unseal the vaults using only the vault-0.  
Vault is sealed by default!  

First, we exec into the Pod.
```
kubectl exec --stdin=true --tty=true vault-0 -n vault -- vault operator init
```

We get a few Keys which need to be run 3 times on the following command.
```
kubectl exec --stdin=true --tty=true vault-0 -n vault -- vault operator unseal
```

We need to join the other vaults to the cluster and unseal them too
```
kubectl exec -ti vault-1 -n vault -- vault operator raft join http://vault-0.vault-internal:8200
kubectl exec -ti vault-1 -n vault -- vault operator unseal
```

```
kubectl exec -ti vault-2 -n vault -- vault operator raft join http://vault-0.vault-internal:8200
kubectl exec -ti vault-2 -n vault -- vault operator unseal
```

Check Status of the pods
```
kubectl exec -it vault-0 -n vault -- vault status
kubectl exec -it vault-1 -n vault -- vault status
kubectl exec -it vault-2 -n vault -- vault status
```

## Check Setup
Check if all vault pods are up and running.  
```kubectl get all -n vault```

## Use Vault in your K8s Apps
You can use it with e.g. Python and the **hvac** module.  
The client you need to create for this will be the service-name of the vault and port 8200.  

I added an example App in the /app folder.