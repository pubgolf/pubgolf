#!/bin/bash

apt-get update
apt-get upgrade

apt-get install fail2ban

useradd deployer
mkdir /home/deployer
mkdir /home/deployer/.ssh
chmod 700 /home/deployer/.ssh
usermod -s /bin/bash deployer
# Contents of ./keys/pubgolf_rsa.pub
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDe3Le4/wUjFZTgbleG/QgYHxXnEa0mwK4TvoGhlbstbDXKr4Mzmf5jWla88hUUXOZfGzns7e1igyN7KQCd8np2MGFORRCSYLhgF/Uf/+gqeclrUtLvo7s0H5DRre6itlcHqKaNQ+S9ndJdVY0Q86Px+nrDL1PUXDPUj6n99otfgKnljHRuy9RrnsRAXvPeoec5BW+w/zlYhyDeJQMcV0vipyJrRffsBvY/MOEU1mDfqlMo5M+XwSIPFyz9c4ywOV2EZnWtUJLneeNvf7GDStuc/iYZdQAm6kordmbqCnx66aFl+G7iNQAxsPWWdFpIoF/qEjs5OI0HY6YzibzjSLEF ericmorris@spacegray-storm-cloud.local" >> /home/deployer/.ssh/authorized_keys
chmod 400 /home/deployer/.ssh/authorized_keys
chown deployer:deployer /home/deployer -R
echo 'deployer ALL=(ALL) NOPASSWD: ALL' | sudo EDITOR='tee -a' visudo

nano /etc/ssh/sshd_config
# PermitRootLogin no
# PasswordAuthentication no

ufw allow 22
ufw allow 80
ufw allow 443
ufw allow 2376
ufw enable

apt-get install unattended-upgrades
echo 'APT::Periodic::Update-Package-Lists "1";' > /etc/apt/apt.conf.d/10periodic
echo 'APT::Periodic::Download-Upgradeable-Packages "1";' >> /etc/apt/apt.conf.d/10periodic
echo 'APT::Periodic::AutocleanInterval "7";' >> /etc/apt/apt.conf.d/10periodic
echo 'APT::Periodic::Unattended-Upgrade "1";' >> /etc/apt/apt.conf.d/10periodic
