# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  config.vm.box = "generic/ubuntu2204"


  # provides rabbitmq ui on localhost:15672
  #config.vm.network "forwarded_port", guest: 15672, host: 15672, host_ip: "127.0.0.1"

  # provides neo4j ui on localhost:7474
  #config.vm.network "forwarded_port", guest: 7474, host: 7474, host_ip: "127.0.0.1"
  #
  config.vm.network "public_network", bridge: "enp0s31f6"

  config.vm.synced_folder ".", "/vagrant", type: "sshfs",
	rsync__exclude: ".vagrant/"

  config.vm.provision "shell", path: "docker.sh"
  config.vm.provision "shell", privileged: false, inline: <<-SHELL
	sudo usermod -aG docker $USER
	newgrp docker
  SHELL
  config.vm.provision "shell", path: "golang.sh", privileged: false
  config.vm.provision "shell", inline: <<-SHELL
  sudo apt install -y build-essential
  SHELL
end
