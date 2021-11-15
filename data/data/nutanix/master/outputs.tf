output "control_plane_ips" {
  value = nutanix_virtual_machine.vm_master.*.nic_list_status.0.ip_endpoint_list.0.ip
}
