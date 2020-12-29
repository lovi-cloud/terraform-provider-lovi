# terraform-provider-lovi

[Terraform](https://www.terraform.io/) provider for lovi-cloud

## Usage

### Authentication

terraform-provider-lovi need to set a API endpoint.

- `SATELIT_API_ENDPOINT`
    - example: 192.0.2.1:9262

### Quick Start

Example of creating virtual machine

```tf
resource "lovi_subnet" "life" {
  name = "life"
  vlan_id = 1000
  network = "192.0.2.0/24"
  start = "192.0.2.100"
  end = "192.0.2.200"
  gateway = "192.0.2.254"
  dns_server = "8.8.8.8"
  metadata_server = "192.0.2.1"
}

resource "lovi_bridge" "life" {
  name = "life"
  vlan_id = 1000

  depends_on = [lovi_subnet.life]
}

resource "lovi_internal_bridge" "life" {
  name = "life-in"
}

resource "lovi_cpu_pinning_group" "life" {
  name = "life"
  count_of_core = 8
  hypervisor_name = "hv0001"
}

resource "lovi_address" "problem-life-1" {
  subnet_id = lovi_subnet.life.id
  fixed_ip = "192.0.2.101"

  depends_on = [lovi_subnet.life]
}

resource "lovi_lease" "problem-life-1" {
  address_id = lovi_address.problem-life-1.id

  depends_on = [lovi_address.problem-life-1]
}

resource "lovi_virtual_machine" "problem-life-1" {
  name = "life-${format("%03d", 1)}"
  vcpus = 2
  memory_kib = 1 * 1024 * 1024
  root_volume_gb = 30
  source_image_id = "00000000-0000-0000-0000-000000000000"
  hypervisor_name = "hv0001"
  europa_backend_name = "europa001"

  read_bytes_sec = 1 * 1000 * 1000 * 1000  // 1Gbps
  write_bytes_sec = 1 * 1000 * 1000 * 1000 // 1Gbps
  read_iops_sec = 800
  write_iops_sec = 800

  cpu_pinning_group_name = lovi_cpu_pinning_group.life.name

  depends_on = [
    lovi_cpu_pinning_group.life
  ]
}
```
