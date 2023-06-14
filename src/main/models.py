from pydantic import BaseModel

class AppConfig(BaseModel):
    """
    This is the configuration Class for the App.
    It uses pydantics BaseModel to declare the Types
    and what happens, when the entry is not defined.
    """
    VAULT_SERVICE_NAME : str | None = None
    VAULT_TOKEN : str | None = None
