# terraform-provider-lovi: [terraform](https://www.terraform.io/) provider for lovi-cloud

## Usage

### Authentication

terraform-provider-lovi need to set a API endpoint.

- `SATELIT_API_ENDPOINT`
    - example: 192.0.2.1:9262

### Quick Start

Example of creating virtual machine

```tf
### Network definition ###

resource "lovi_subnet" "life" {
  name            = "life"
  vlan_id         = 1000
  network         = "192.0.2.0/24"
  start           = "192.0.2.100"
  end             = "192.0.2.200"
  gateway         = "192.0.2.254"
  dns_server      = "8.8.8.8"
  metadata_server = "192.0.2.1"
}

resource "lovi_bridge" "life-1" {
  name    = "life-1" # UNIQUE in lovi_bridge and lovi_internal_bridge
  vlan_id = 1000     # same value in lovi_subnet.vlan_id

  depends_on = [lovi_subnet.life]
}

resource "lovi_internal_bridge" "life-1" {
  name = "life-1-in"
}

# lovi_cpu_pinning_group is Pinning group in vCPU. (optional)
resource "lovi_cpu_pinning_group" "life-1" {
  name            = "life-1"
  count_of_core   = 2
  hypervisor_name = "hv0001"
}


### Virtual Machine definition ###

resource "lovi_address" "life-1" {
  subnet_id = lovi_subnet.life.id
  fixed_ip  = "192.0.2.101" # (optional)

  depends_on = [lovi_subnet.life]
}

resource "lovi_lease" "life-1" {
  address_id = lovi_address.life-1.id

  depends_on = [lovi_address.life-1]
}

resource "lovi_virtual_machine" "vm-1" {
  name                = "life-${format("%03d", 1)}"
  vcpus               = 2
  memory_kib          = 1 * 1024 * 1024
  root_volume_gb      = 30
  source_image_id     = "00000000-0000-0000-0000-000000000000"
  hypervisor_name     = "hv0001"
  europa_backend_name = "europa001"

  read_bytes_sec  = 1 * 1000 * 1000 * 1000 // 1Gbps
  write_bytes_sec = 1 * 1000 * 1000 * 1000 // 1Gbps
  read_iops_sec   = 800
  write_iops_sec  = 800

  cpu_pinning_group_name = lovi_cpu_pinning_group.life-1.name  # (optional)

  depends_on = [
    lovi_cpu_pinning_group.life-1
  ]
}

resource "lovi_interface_attachment" "vm-1" {
  virtual_machine_id = lovi_virtual_machine.vm-1.id
  bridge_id = lovi_bridge.life-1.id
  average = 125000 // NOTE: 1Gbps
  name = "life-1-${format("%03d", 1)}"
  lease_id = lovi_lease.life-1.id

  depends_on = [
    lovi_virtual_machine.vm-1,
    lovi_lease.life-1
  ]
}
```
