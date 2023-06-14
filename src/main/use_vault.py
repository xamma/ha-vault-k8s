import hvac
import config

"""
This is an example implemenation of using the Vault
with the python hvac library. We create a client with the
URL and token for authentication and can then create the
secret or read secrets.
"""

# Configure the Vault client
client = hvac.Client(
    url=f"http://{config.vault_service_name}:8200",
    token=config.vault_token
    )

# Check connection
if client.is_authenticated():
    print("Successfully authenticated with Vault.")

# Write a secret to Vault
data = {
    'username': 'myuser',
    'password': 'mypassword'
}
# Path is the name of the secret
client.secrets.kv.v2.create_or_update_secret(path='secret/path', secret=data)

# Read a secret from Vault
secret = client.secrets.kv.v2.read_secret_version(path='secret/path')
print(secret['data'])