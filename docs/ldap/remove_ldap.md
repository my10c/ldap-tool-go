
## Background
This document explain how to remove LDAP on a server

## Technologies
The following technologies are used

- servers/instances running Ubuntu 20.04
- OpenLDAP as the LDAP software


#### Note
The steps should be done on all servers primary and secondaries (replicas)

### Remove the packages
```
apt remove --purge -y slapd ldap-utils nscd nslcd nslcd-utils ldapscripts
```

### Replace SUDO ldap
```
export SUDO_FORCE_REMOVE=yes
apt remove -y sudo-ldap
apt install -y sudo
```
Note:
```
make sure to adjust the sudo setting
use visudo
```

### Reset the authentication
A quick way to disable LDAP authentication is to remove any *ldap* sources from the 
file */etc/nsswitch.conf*
below is an example **with** LDAP
```
passwd: compat systemd ldap
group: compat systemd ldap
shadow: compat ldap

gshadow:        files

hosts:          files dns
networks:       files

protocols:      db files
services:       db files
ethers:         db files
rpc:            db files

netgroup:       nis
sudoers:        files ldap
```
below is an exmple **without** LDAP
```
passwd: compat systemd
group: compat systemd
shadow: compat

gshadow:        files

hosts:          files resolve [!UNAVAIL=return] dns
networks:       files

protocols:      db files
services:       db files
ethers:         db files
rpc:            db files

netgroup:       nis
```


### Delete the left over files and packages (optional)
To **permanently** delete
```
rm -rf /var/lib/slapd
```

edit these file and remove any mention of the ldap settings
```
vi /etc/libnss-ldap.conf /etc/pam_ldap.conf
```

To **save** the ldap files (db files) in case you need it back
```
mv /var/lib/slapd /var/lib/slapd-$(date "+%Y-%m-%d-%H-%M")
```
Remove the no longer needed packages
```
apt autoremove
apt autoclean
```

## The End
Congraculation you should be all set now : 🦄👏
 
