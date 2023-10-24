resource "azurerm_resource_group" "cyucra-sandbox" {
  name     = "cyucra-sandbox"
  location = "East US"
}

resource "azurerm_container_registry" "cyucraacr" {
  name                = "cyucraacr"
  resource_group_name = azurerm_resource_group.cyucra-sandbox.name
  location            = azurerm_resource_group.cyucra-sandbox.location
  sku                 = "Basic"
  admin_enabled       = true
}

locals {
  current_user_id = coalesce(var.msi_id, data.azurerm_client_config.current.object_id)
}


data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "key_vault" {
  name                        = "cyucra-sandbox-kv"
  location                    = azurerm_resource_group.cyucra-sandbox.location
  resource_group_name         = azurerm_resource_group.cyucra-sandbox.name
  #enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = local.current_user_id
    # object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create",
    ]

    secret_permissions = [
      "Get", "List", "Delete", "Set", 
    ]

    storage_permissions = [
      "Get", "List",
    ]
  }
}

resource "azurerm_user_assigned_identity" "managed_identity" {
  name                = "cyucra-sandbox-mi"
  location            = azurerm_resource_group.cyucra-sandbox.location
  resource_group_name = azurerm_resource_group.cyucra-sandbox.name
}

# to show all definitions https://learn.microsoft.com/en-us/azure/key-vault/general/rbac-migration?WT.mc_id=AZ-MVP-5004151
resource "azurerm_role_assignment" "assign_identity_storage_blob_data_contributor" {
  scope                = azurerm_key_vault.key_vault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.managed_identity.principal_id
}

provider "azurerm" {
  features { }
}