locals {
  description = "Created By OpenShift Installer"
}

provider "nutanix" {
  wait_timeout = 60
  username     = var.username
  password     = var.password
  endpoint     = var.prism_central
  port         = var.port
  insecure     = true
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
      uuid = var.image_id
    }
    disk_size_mib = var.nutanix_control_plane_disk_mib
  }

  categories {
    name  = var.ocp_category_key_id
    value = var.ocp_category_value_id
  }

  guest_customization_cloud_init_user_data = base64encode(var.ignition_master)
  nic_list {
    subnet_uuid = var.nutanix_subnet_uuid
  }
}
