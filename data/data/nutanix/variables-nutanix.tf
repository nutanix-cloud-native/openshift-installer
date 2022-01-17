//////
// Nutanix variables
//////

variable "prism_central" {
  type        = string
  description = "This is the Prism Central for the environment."
}

variable "port" {
  type        = string
  default     = 9440
  description = "Port to connect to Prism Central."
}

variable "insecure" {
  type        = bool
  description = "Disable certificate checking when connecting to Prism Central"
}


variable "username" {
  type        = string
  description = "Prism Central user for the environment."
}

variable "password" {
  type        = string
  description = "Prism Central server password"
}

variable "nutanix_prism_element_uuid" {
  type        = string
  description = "This is the uuid of the Prism Element cluster."
}


variable "nutanix_image_filepath" {
  type        = string
  description = "This is the filepath to the image file that will be imported into Prism Central."
}

variable "nutanix_image" {
  type        = string
  description = "This is the name to the image that will be imported into Prism Central."
}

variable "nutanix_subnet_uuid" {
  type        = string
  description = "This is the uuid of the publicly accessible subnet for cluster ingress and access."
}

variable "nutanix_bootstrap_ignition_image" {
  type        = string
  description = "Path to the image containing the bootstrap ignition files"
}


///////////
// Control Plane machine variables
///////////

variable "nutanix_control_plane_memory_mib" {
  type = number
}

variable "nutanix_control_plane_disk_mib" {
  type = number
}

variable "nutanix_control_plane_num_cpus" {
  type = number
}

variable "nutanix_control_plane_cores_per_socket" {
  type = number
}
