output "image_id" {
  value = nutanix_image.rhcos.id
}

output "ocp_category_key_id" {
  value = nutanix_category_key.ocp_category_key.id
}

output "ocp_category_value_id" {
  value = nutanix_category_value.ocp_category_value.id
}

output "control_plane_ips" {
  value = nutanix_virtual_machine.vm_master[*].nic_list_status[0].ip_endpoint_list[0].ip
}
