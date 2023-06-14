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

### Init Vault
First, we exec into the Pod.
```
kubectl exec --stdin=true --tty=true vault-0 -n vault -- vault operator init
```

### Unseal Vault
We get a few Keys which need to be run 3 times on the following command.
```
kubectl exec --stdin=true --tty=true vault-0 -n vault -- vault operator unseal
```

### Join other pods to cluster and unseal
We need to join the other vaults to the cluster and unseal them too
```
kubectl exec -ti vault-1 -n vault -- vault operator raft join http://vault-0.vault-internal:8200
kubectl exec -ti vault-1 -n vault -- vault operator unseal
```

```
kubectl exec -ti vault-2 -n vault -- vault operator raft join http://vault-0.vault-internal:8200
kubectl exec -ti vault-2 -n vault -- vault operator unseal
```

### Check Status
Check Status of the pods
```
kubectl exec -it vault-0 -n vault -- vault status
kubectl exec -it vault-1 -n vault -- vault status
kubectl exec -it vault-2 -n vault -- vault status
```

### Check Setup
Check if all vault pods are up and running.  
```kubectl get all -n vault```

## Create an access token for your apps
The root token should only be used for setup and administrative access, so we need to create a token for our client to authenticate against vault. 

### Create ACL
First, we have to create an ACL policy what the token can do:  
Log into the vault ui and create a new policy called **basic** with the following content  
```
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
``` 
This creates a policy that grants read and write access for all paths within the secret/ prefix.  
Modify it according to your needs and paths.  

### Create Token
Now we can create our client token with this policy attached.  
```
kubectl exec --stdin=true --tty=true vault-0 -n vault -- env VAULT_TOKEN=<YOUR_ROOT_TOKEN> vault token create -policy=basic -format=json
```
We can now use this token as ENV-Var in our App. Best way is to pass it as Kubernetes-Secret.    

## Use Vault in your K8s Apps
You can use it with e.g. Python and the **hvac** module.  
The client you need to create for this will be the service-name of the vault and port 8200.  
Be careful to use the FQDN format SERVICENAME.NAMESPACE.svc.cluster.local since they are not in the same namespace, e.g. ```http://vault.vault-namespace.svc.cluster.local:8200``` is the needed URL so we specify our VAULT_SERVICE_NAME with ```vault.vault.svc.cluster.local```.

I added an example App in the /src folder which is also already containerized on **xamma/python-vault**.  
Set the ENVs e.g
```
$env:VAULT_TOKEN="dev-only-token"
$env:VAULT_SERVICE_NAME="localhost"
```
You can also try it with docker and use localhost as URL/vault-service-name:
```
docker run -p 8200:8200 -e 'VAULT_DEV_ROOT_TOKEN_ID=dev-only-token' vault
```

### Create K8s Secret for vault token
To not store the secret in the manifest we need to create it via CLI.  
```
kubectl create secret generic my-app-secrets -n <NAMESPACE> \
  --from-literal=VAULT_TOKEN=<value>
```

### Create Secrets Engine
To use/create the secrets you must first enable a new Secrets Engine via the vault CLI / UI.  
To do so, log into the UI with the root token, and click on  
```Secrets > Enable new Engine > KV > specify path (in our case, the Path is "secret" and the name of secret is "userdata")> enable```.

### Run the app
Run the ```pod.yaml``` specified in the /k8s-manifests folder and verify it wrote/read your secret.  
Make sure to create the namespace vault-test.  
If it run successfully, there will be a new entry in your secret/ when visiting the WebUI.  
You can also check the Container Logs.  

### BONUS: Go-App
I created the same App using GoLang, so you can see the differences.  
You can find it in the folder /app.  
The App is also already containerized and ready to use, find the image on my Dockerhub: ```xamma/go-vault:latest```.  
Be careful, the App uses the Path **/gosecret** now, so you have to update your ACL's and create a new secrets engine.  

Initialize: ```go mod init example.com/govault```  
Add needed packages to go.mod: ```go mod tidy```  
Run the App: ```go run main.go config.go```  
Build: ```go build```   