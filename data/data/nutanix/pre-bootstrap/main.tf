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

resource "nutanix_category_key" "ocp_category_key" {
  name        = "openshift-${var.cluster_id}"
  description = "Openshift Cluster Category Key"
}

resource "nutanix_category_value" "ocp_category_value" {
  name        = nutanix_category_key.ocp_category_key.id
  value       = "openshift-ipi-installations"
  description = "Openshift Cluster Category Value"
}

resource "nutanix_image" "rhcos" {
  name        = var.nutanix_image
  source_path = var.nutanix_image_filepath
  description = local.description

  categories {
    name  = nutanix_category_key.ocp_category_key.name
    value = nutanix_category_value.ocp_category_value.value
  }
}

resource "nutanix_virtual_machine" "vm_master" {
  count                = var.master_count
  description          = local.description
  name                 = "${var.cluster_id}-master-${count.index}"
  cluster_uuid         = var.nutanix_prism_element_uuid
  num_vcpus_per_socket = var.nutanix_control_plane_cores_per_socket
  num_sockets          = var.nutanix_control_plane_num_cpus
  memory_size_mib      = var.nutanix_control_plane_memory_mib
  boot_device_order_list = [
    "DISK",
    "CDROM",
    "NETWORK"
  ]
  disk_list {
    device_properties {
      device_type = "DISK"
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
    }
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.rhcos.id
    }
    disk_size_mib = var.nutanix_control_plane_disk_mib
  }

  categories {
    name  = nutanix_category_key.ocp_category_key.name
    value = nutanix_category_value.ocp_category_value.value
  }

  guest_customization_cloud_init_user_data = base64encode(var.ignition_master)
  nic_list {
    subnet_uuid = var.nutanix_subnet_uuid
  }

  depends_on = [nutanix_image.rhcos]
}
