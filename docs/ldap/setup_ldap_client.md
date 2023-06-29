## Background
This document explain how to setup a LDAP client.
Do note that a LDAP servers can also be a LDAP client

## Technologies
The following technologies are used
- servers/instances running Ubuntu 20.04+
- openLDAP

## Prerequisite
- The primary server is ldap-primary and the secondary server is ldap-secondary
- Define the LDAP domain : co.my10c.com
- The LDAP server are running properly with the users and groups data 
- Create the file /etc/ldap/ldap.conf with the below content (create directory with `mkdir -p /etc/ldap`)
```
sizelimit   12
timelimit   15
deref       never

base dc=co,dc=my10c,dc=com
uri ldap://ldap.co.my10c.com
ldap_version 3
rootbinddn cn=admin,dc=co,dc=my10c,dc=com

# normal password encryption mode
pam_password md5

sudoers_base ou=SUDOers,dc=co,dc=my10c,dc=com
sudoers_timed no
sudoers_debug 0
# ssl start_tls
index sudoUser eq

# TLS certificates (needed for GnuTLS)
tls_cacert /etc/ldap/certs/ca.pem
tls_cert /etc/ldap/certs/cert.pem
tls_key /etc/ldap/certs/key.pem
tls_cipher_suite SECURE256:SECURE:-VERS-TLS1.1:-VERS-TLS1.0:-VERS-SSL3.0:-MD5:-SHA1:-ARCFOUR-128
tls_protocol_min 3.3
tls_reqcert never
create the ldap password files
```
- Create the symlink
```
ln -s /etc/ldap/ldap.conf /etc/ldap.conf
```
- installed make sure the /etc/nsswitch.conf content is (after the packages has been installed)


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


### Install the client packages
```
apt install ldap-utils nscd nslcd sudo-ldap
```

#### Configure the client authentication
```
dpkg-reconfigure ldap-auth-config
    LDAP server Uniform Resource Identifier:    ldap://ldap://ldap.co.my10c.com
    Distinguished name of the search base:      dc=co,dc=my10c,dc=com
    LDAP version to use:                        3
    Should debconf manage LDAP configuration?   yes
    Make local root Database admin:             yes
    Does the LDAP database require login?       No
    LDAP account for root:                      cn=admin,dc=co,dc=my10c,dc=com
    LDAP root account password:                 <password>
```

## The End
Congraculation you should be all set now : 🦄👏
