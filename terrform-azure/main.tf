# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }

  required_version = ">= 1.1.0"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "goeats"
  location = "southeastasia"
}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "goeats-k8scluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "goeatsk8scluster"
  http_application_routing_enabled = true

  default_node_pool {
    name       = "default"
    node_count = "1"
    vm_size    = "Standard_B2s"
  }

  identity {
    type = "SystemAssigned"
  }
}