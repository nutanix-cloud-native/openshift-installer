locals {
  description = "Created By OpenShift Installer"
}

provider "nutanix" {
  wait_timeout = 60
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_prism_central_address
  port         = var.nutanix_prism_central_port
}

resource "nutanix_image" "rhcos" {
  name        = var.nutanix_image
  source_path = var.nutanix_image_filepath
  description = local.description
}

resource "nutanix_category_key" "ocp_category_key" {
  name = "openshift-${var.cluster_id}"
  description = "Openshift Cluster Category Key"
}

resource "nutanix_category_value" "ocp_category_value" {
  name  = nutanix_category_key.ocp_category_key.id
  value = "openshift-ipi-installations"
  description = "Openshift Cluster Category Value"
}

resource "nutanix_image" "bootstrap_ignition" {
  name        = var.nutanix_bootstrap_ignition_image
  source_path = var.nutanix_bootstrap_ignition_image_filepath
  description = local.description
  categories {
    name  = nutanix_category_key.ocp_category_key.id
    value = nutanix_category_value.ocp_category_value.id
  }
}
