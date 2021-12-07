
output "image_id" {
  # value = nutanix_image.rhcos.id
  value = data.nutanix_image.rhcos.id
}

output "bootstrap_ignition_image_id" {
  value = nutanix_image.bootstrap_ignition.id
}

output "cluster_domain" {
  value = var.cluster_domain
}

output "cluster_id" {
  value = var.cluster_id
}

output "ocp_category_key_id" {
  value = nutanix_category_key.ocp_category_key.id
}

output "ocp_category_value_id" {
  value = nutanix_category_value.ocp_category_value.id
}
