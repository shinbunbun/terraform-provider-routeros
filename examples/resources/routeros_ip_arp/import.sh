#The ID can be found via API or the terminal
#The command for the terminal is -> :put [/ip/arp get [print show-ids]]
terraform import routeros_ip_arp.ap "*1"
#Or you can import a resource using one of its attributes
terraform import routeros_ip_arp.ap "address=192.168.88.2"
