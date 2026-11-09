resource "routeros_ip_arp" "ap" {
  address     = "192.168.88.2"
  interface   = "bridge"
  mac_address = "02:00:00:00:00:02"
  comment     = "access point"
}
