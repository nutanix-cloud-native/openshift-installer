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

resource "nutanix_virtual_machine" "vm_bootstrap" {
  name                 = "${var.cluster_id}-bootstrap"
  description          = local.description
  cluster_uuid         = var.nutanix_prism_element_id
  num_vcpus_per_socket = 4
  num_sockets          = 1
  memory_size_mib      = 16384
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
      uuid = var.image_id
    }
    disk_size_bytes = var.nutanix_control_plane_disk_mib * 1024 * 1024 * 1024
  }

  disk_list {
    device_properties {
      device_type = "CDROM"
      disk_address = {
        adapter_type = "IDE"
        device_index = 0
      }
    }
    data_source_reference = {
      kind = "image"
      uuid = var.bootstrap_ignition_image_id
    }
  }
  categories {
    name  = var.ocp_category_key_id
    value = var.ocp_category_value_id
  }


  nic_list {
    subnet_uuid = var.subnet_id
  }
}

