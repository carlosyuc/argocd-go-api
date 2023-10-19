resource "azurerm_resource_group" "cyucra-sandbox" {
  name     = "cyucra-sandbox"
  location = "East US"
}

resource "azurerm_container_registry" "cyucraacr" {
  name                = "cyucraacr"
  resource_group_name = azurerm_resource_group.cyucra-sandbox.name
  location            = azurerm_resource_group.cyucra-sandbox.location
  sku                 = "Basic"
  admin_enabled       = false
#   georeplications {
#     location                = "East US"
#     zone_redundancy_enabled = true
#     tags                    = {}
#   }
#   georeplications {
#     location                = "North Europe"
#     zone_redundancy_enabled = true
#     tags                    = {}
#   }
}

provider "azurerm" {
  features { }
}