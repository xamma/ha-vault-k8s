import os
from models import AppConfig

"""
This is the configuration module for the app
It lets you specify predefined values and
provides the possibility to set ENVs directly
e.g in Kubernetes as secret or via CLI.
"""

# Define the environment variables and their default values
default_env_vars = {
    "VAULT_SERVICE_NAME": "localhost",
    "VAULT_TOKEN": "my-vault-token",
    # Add more environment variables here
}

# Check with Pydantic model
config_settings = AppConfig(**default_env_vars)
# Class object, so we need .dict() for dictionary and .items() to get the key-values
# print(type(config_settings)) == class 'models.AppConfig'

# Set the environment variables only if they are not already set
for key, value in config_settings.dict().items():
    os.environ.setdefault(key, value)

# Get the Vault values from the environment variables
vault_service_name = os.environ.get("VAULT_SERVICE_NAME")
vault_token = os.environ.get("VAULT_TOKEN")

# print(vault_service_name)
# print(vault_token)
