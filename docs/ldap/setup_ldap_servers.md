## Background
This document explain how to setup a LDAP infrastructure,.

## Technologies
The following technologies are used

- servers/instances running Ubuntu 20.04
- 2 Instances, one for the primary and one for fallback (replica)
- OpenLDAP as the LDAP software
- Let’s Encrypt to automatically getting a public SSL certificate
- Nginx, web server to host the LDAP admin and self service password
- Optional LDAP TOTP module for 2FA

## Prerequisite
- The primary server is ldap-primary and the secondary server is ldap-secondary
- Define the LDAP domain : co.my10c.com
- Make sure a R53 private zone is create and attached to the VPC
- create the DNS records in the DNS zone
- The admin DN : cn=admin,dc=co,dc=my10c,dc=com
- generate a password 25 length (no special charter) (indicated with <password> further in this document)
```
< /dev/urandom tr -dc A-Za-z0-9 | head -c${1:-25}
```
- generate a password 16 length (no special charters) for replication (indicated with <replica-password> further in this document)
```
< /dev/urandom tr -dc A-Za-z0-9 | head -c${1:-16}
```
- Obtain the SUDO and OPENSSH Schema
- Obtain the LDIF to setup memberOf in LDAP server
- Optional compiled the OpenLDAP TOTP module, so the src package to obtain the sources

The domain  co.my10c.com is chosen to make getting cert from Let’s Encrypt easier

# Note
- becuase we doing TLS, we must use name, not IP and the name must match the certificate used by the primary.
- is you are using subdomain such in the example, you **must** set the hostname to the subdomain!
  en edit the host file
```
hostnamectl set-hostname ldap-main.co.my10c.com
```


## Setup the Primary and Secondaries servers
To make things easier create the ldap configuration (ldap.conf)
already has the setup for SUDO LDAP
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

```
mkdir -p /etc/ldapscripts
echo -n $(cat ldap.pass) > /etc/ldap.secret
echo -n $(cat ldap.pass)> /etc/ldapscripts/ldapscripts.passwd
chmod 0400 /etc/ldap.secret /etc/ldapscripts/ldapscripts.passwd
```
 

### Install the base packages
```
apt install slapd ldap-utils nscd nslcd nslcd-utils ldapscripts
```

#### Configure the server
```
sudo dpkg-reconfigure slapd
    Omit OpenLDAP configuration Yes/No : No
    DNS domain name : co.my10c.com
    Organization name : Badassops LLC
    Administrator Password <password>
    Confirm Admin Password <password>
    Database purged if slapd is purged: Yes
    Move old database: Yes
```

adjust the ldapscript file /etc/ldapscripts/ldapscripts.conf with this content
```
SERVER=localhost
BINDDN='cn=admin,dc=co,dc=my10c,dc=com'
BINDPWDFILE="/etc/ldapscripts/ldapscripts.passwd"
SUFFIX='dc=co,dc=my10c,dc=com'
GSUFFIX='ou=groups'
USUFFIX='ou=users'
MSUFFIX='ou=Computers'
GIDSTART=8000
UIDSTART=8000
MIDSTART=1000
```
 

#### Setup password encryption salt format to SHA-512
Create the sha512.ldif file with the following content
```
dn: cn=config
replace: olcPasswordCryptSaltFormat
olcPasswordCryptSaltFormat: $%.16
```

```
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f sha512.ldif
```
 

#### Setup LDAP to support SUDO and OPENSSH

Copy the SUDO and OPENSSH schema and ldif files under /etc/ldap/schema then activate the schema’s and then restart the LDAP service (slapd)


```
cp openssh.ldif openssh.schema  sudo.ldif sudo.schema  /etc/ldap/schema/
chmod 0444 /etc/ldap/schema/{openssh.ldif,openssh.schema,sudo.ldif,sudo.schema}
chown root:root /etc/ldap/schema/{openssh.ldif,openssh.schema,sudo.ldif,sudo.schema}
```

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f openssh.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f sudo.ldif
systemctl restart slapd
```

#### Setup logging

Create the file log.ldif for logging then load the configuration

```
dn: cn=config
changetype: modify
replace: olcLogLevel
olcLogLevel: stats shell sync
```

```
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f log.ldif
```
 
#### Setup Cert to enable TLS

We are assume that Let’s encrypt was installed and out renew cron/script were create and enabled.
(beyond this documentation)
The certicates were installed under /etc/ldap/certs (mkdir -p /etc/ldap/certs)

```
dn: cn=config
changetype: modify
replace: olcTLSCACertificateFile
olcTLSCACertificateFile: /etc/ldap/certs/ca.pem
-
replace: olcTLSCertificateFile
olcTLSCertificateFile: /etc/ldap/certs/cert.pem
-
replace: olcTLSCertificateKeyFile
olcTLSCertificateKeyFile: /etc/ldap/certs/key.pem
```

```
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f cert.ldif
```
 

#### Setup rsyslog

create the file, 40-ldap.conf with the following content
```
local4.* /var/log/ldap.log
local4.* stop 
```

Copy the file to the rsyslog configuration and set owner and permission of the file and setup the initial rsyslog file and restart the services, rsyslog and slapd

```
cp 40-ldap.conf /etc/rsyslog.d/40-ldap.conf
chmod 0444 /etc/rsyslog.d/40-ldap.conf
chown root:root /etc/rsyslog.d/40-ldap.conf
touch /var/log/ldap.log
chown syslog:adm /var/log/ldap.log
chmod 0640 /var/log/ldap.log
systemctl restart rsyslog
systemctl restart slapd
```

#### Setup logrotate 

configuration for ldap by  /etc/logrotate.d/ldap with this content
```
/var/log/ldap.log
{
    rotate 7
    weekly
    missingok
    notifempty
    compress
    delaycompress
    sharedscripts
}
```
Set ownership and permission

```
chmod 0444 /etc/logrotate.d/ldap 
chown root:root /etc/logrotate.d/ldap 
``` 

#### Setup sudo

For sudo to work with sudo a special package need to be installed and a configuration create

```
export SUDO_FORCE_REMOVE=yes
apt install sudo-ldap
```
the sudo configuration can be add to the main ldap configuration, so we need to delete the configuration and symbolic the configuration with the ldap main configuration 

```
mv /etc/sudo-ldap.conf /etc/sudo-ldap.conf.dpkg
ln -s /etc/ldap/ldap.conf /etc/sudo-ldap.conf
```

# THESE SHOULD ONLY BE EXECUTED ON THE PRIMARY LDAP SERVER ONLY

create the group and sudoers LDAP configuration records
content of the group.ldif file 
These is are example where we use the user ** and the group &vpn
however we do need the group *nogroup*

```
dn: ou=users,dc=co,dc=my10c,dc=com
objectClass: organizationalUnit
ou: Users

dn: ou=groups,dc=co,dc=my10c,dc=com
objectClass: organizationalUnit
ou: Groups

dn: cn=ansible,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: ansible
gidNumber: 400

dn: cn=devops,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: devops
gidNumber: 777

dn: cn=sre,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: sre
gidNumber: 778

dn: cn=engineering,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: engineering
gidNumber: 2000

dn: cn=datascientist,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: datascientist
gidNumber: 2010

dn: cn=qa,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: qa
gidNumber: 2020

dn: cn=nogroup,ou=groups,dc=co,dc=my10c,dc=com
objectClass: posixGroup
cn: nogroup
gidNumber: 10666

```

content of sudoers.ldif file

```
dn: ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: organizationalUnit
ou: SUDOers

dn: cn=defaults,ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: sudoRole
cn: defaults
sudoOption: ignore_dot
sudoOption: shell_noargs
sudoOption: mail_badpass
sudoOption: mail_no_user
sudoOption: mail_no_host
sudoOption: mail_no_perms
sudoOption: insults
sudoOption: log_year
sudoOption: env_keep+="HOME PS1 PS2 PROMPT"
sudoOption: secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
sudoOrder: 1

dn: cn=root,ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: sudoRole
cn: root
sudoUser: root
sudoHost: ALL
sudoCommand: ALL
sudoOption: !authenticate
sudoOrder: 2

dn: cn=ALL,ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: sudoRole
cn: ALL
sudoUser: ALL
sudoHost: ALL
sudoCommand: /usr/bin/sudo -l
sudoOption: !authenticate
sudoOrder: 2

dn: cn=%devops,ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: sudoRole
cn: %devops
sudoUser: %devops
sudoHost: ALL
sudoRunAsUser: ALL
sudoCommand: ALL
sudoOption: !authenticate
sudoOrder: 3

dn: cn=%sre,ou=SUDOers,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: sudoRole
cn: %sre
sudoUser: %sre
sudoHost: ALL
sudoRunAsUser: ALL
sudoCommand: ALL
sudoOption: !authenticate
sudoOrder: 3
```

Add the groups and sudoers configuration
```
ldapadd -c -x -W -D cn=admin,dc=co,dc=my10c,dc=com -f group.ldif
ldapadd -c -x -W -D cn=admin,dc=co,dc=my10c,dc=com -f sudoers.ldif
```
 

# THESE SHOULD BE EXECUTED ON THE ALL LDAP SERVER

### Setup the memberOf configuration

This should be done on a new installation, for existing installation adjustment will need to be done by finding the right internal LDAP db, and is beyond the scope of this document
Add the memberof plugin, create the file memberof-add-plugin.ldif 

The 2 steps (memberof-add-plugin.ldif and memberof-configuration.ldif must to be loaded on the replicas

#### Configure the memberof overlay, create file memberof-configuration.ldif

```
dn: cn=module,cn=config
cn: module
objectclass: olcModuleList
objectclass: top
olcmoduleload: memberof.la
olcmodulepath: /usr/lib/ldap
```

```
dn: olcOverlay=memberof,olcDatabase={1}mdb,cn=config
changetype: add
objectClass: olcMemberOf
objectClass: olcOverlayConfig
objectClass: olcConfig
objectClass: top
olcOverlay: memberof
olcMemberOfDangling: ignore
olcMemberOfRefInt: TRUE
olcMemberOfGroupOC: groupOfNames
olcMemberOfMemberAD: member
olcMemberOfMemberOfAD: memberOf
```

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof-add-plugin.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof-configuration.ldif
```

# THESE SHOULD ONLY BE EXECUTED ON THE PRIMARY LDAP SERVER ONLY

We need to create an temporary user and then add to the special group VPN (as memberOf), we need these 2 files

initial-user.ldif
```
dn: uid=initial-user,ou=users,dc=co,dc=my10c,dc=com
objectClass: top
objectClass: person
objectClass: organizationalPerson
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
objectClass: ldapPublicKey
uid: initial-user
sn: LDAP
givenName: Initial-user
cn: Initial-user LDAP
mail: initial-user@my10c.com
displayName: Initial-user LDAP
uidNumber: 666
gidNumber: 66
userPassword: verySecretReally
gecos: Initial-user LDAP
departmentNumber: LDAP
loginShell: /bin/zsh
homeDirectory: /home/initial-user
shadowLastChange: 18322
shadowExpire: -1
shadowMax: 180
shadowWarning: 14
sshPublicKey: none
```

memberof-groups.ldif
```
dn: cn=vpn,ou=groups,dc=co,dc=my10c,dc=com
objectClass: groupOfNames
cn: vpn
member: uid=initial-user,ou=users,dc=co,dc=my10c,dc=com

dn: cn=admins,ou=groups,dc=co,dc=my10c,dc=com
objectClass: groupOfNames
cn: admins
member: uid=initial-user,ou=users,dc=co,dc=my10c,dc=com

```

load the file into the ldap server
```
ldapadd  -c -x -W -D cn=admin,dc=co,dc=my10c,dc=com -f initial-user.ldif
ldapadd  -c -x -W -D cn=admin,dc=co,dc=my10c,dc=com -f memberof-groups.ldif
systemctl restart slapd
```
 

# THESE SHOULD BE EXECUTED ON THE ALL LDAP SERVER (OPTIONAL)

## Setup index (optional)
This will make lookup faster, we index on uid if required, so we check first, normally the uid on idex is always create.


```
ldapsearch -Q -LLL -Y EXTERNAL -H ldapi:/// -b cn=config '(olcDatabase={1}mdb)' olcDbIndex
```
ouput might look like this
```
dn: olcDatabase={1}mdb,cn=config
olcDbIndex: objectClass eq
olcDbIndex: cn,uid eq
olcDbIndex: uidNumber,gidNumber eq
olcDbIndex: member,memberUid eq
olcDbIndex: entryCSN eq
olcDbIndex: entryUUID eq

dn: olcDatabase={2}mdb,cn=config
olcDbIndex: default eq
olcDbIndex: entryCSN,objectClass,reqEnd,reqResult,reqStart
```
We see there is an index on uid, if there were none then we create the uid_index.ldif with this content

uid_index.ldif
```
dn: olcDatabase={1}mdb,cn=config
add: olcDbIndex
olcDbIndex: uid eq
```

This to add index on memberof
memberof_index.ldif
```
dn: olcDatabase={1}mdb,cn=config
add: olcDbIndex
olcDbIndex: memberof eq
```

Then add the configuration

```
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f uid_index.ldif
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f memberof_index.ldif
```

## The End
Congraculation you should be all set now : 🦄👏
 
