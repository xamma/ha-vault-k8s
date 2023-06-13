import os
import hvac
import config

# Get the Vault service name from the environment variable
vault_service_name = os.environ.get("VAULT_SERVICE_NAME")

# Configure the Vault client
client = hvac.Client(
    url=f"http://{vault_service_name}:8200",
    token=config.MY_VAULT_TOKEN
    )

# Write a secret to Vault
data = {
    'username': 'myuser',
    'password': 'mypassword'
}
client.secrets.kv.v2.create_or_update_secret(path='secret/path', secret=data)

# Read a secret from Vault
secret = client.secrets.kv.v2.read_secret_version(path='secret/path')
print(secret['data'])