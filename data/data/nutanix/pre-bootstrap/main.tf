locals {
  description = "Created By OpenShift Installer"
}

provider "nutanix" {
  wait_timeout = 60
  username     = var.username
  password     = var.password
  endpoint     = var.prism_central
  insecure     = var.insecure
  port         = var.port
}

#TODO: Upload instead of query
resource "nutanix_image" "rhcos" {
  name        = var.nutanix_image
  source_path  = var.nutanix_image_filepath
  description = local.description
}

//for dev
# data "nutanix_image" "rhcos" {
#   image_name = "rhcos-manual"
# }

resource "nutanix_category_key" "ocp_category_key" {
  name = "openshift-${var.cluster_id}"
}

resource "nutanix_category_value" "ocp_category_value" {
  name  = nutanix_category_key.ocp_category_key.id
  value = "openshift-ipi-installations"
}

resource "nutanix_image" "bootstrap_ignition" {
  name        = var.nutanix_bootstrap_ignition_image
  source_path = var.nutanix_bootstrap_ignition_image
  description = local.description
  categories {
    name  = nutanix_category_key.ocp_category_key.id
    value = nutanix_category_value.ocp_category_value.id
  }
}
