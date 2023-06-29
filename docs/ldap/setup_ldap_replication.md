# Background
This document explain how to setup a LDAP replica

# Technologies
The following technologies are used

- AWS EC2 instances running Ubuntu 20.04
- 2 Instances, one for the primary and one for fallback (replica)
- OpenLDAP as the LDAP software

# Setup
NOTE:
- becuase we doing TLS, we must use name, not IP and the name must match the certificate used by the primary.
- is you are using subdomain such in the example, you **must** set the hostname to the subdomain!
  en edit the host file
```
hostnamectl set-hostname  ldap-main.co.my10c.com
```


### Create a password for the replicator user
```
< /dev/urandom tr -dc A-Za-z0-9 | head -c${1:-16} > replicator.pass
chmod 0400 replicator.pass
chown root:root replicator.pass
```

### Create the provider file, this will be use on the primary server
File provider.ldif content
```
dn: cn=module,cn=config
objectClass: olcModuleList
cn: module
olcModulePath: /usr/lib/ldap
olcModuleLoad: syncprov

dn: olcOverlay=syncprov,olcDatabase={1}mdb,cn=config
objectClass: olcOverlayConfig
objectClass: olcSyncProvConfig
olcOverlay: syncprov
olcSpSessionLog: 100
On the primary load the configuration and restart the service
```

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f provider.ldif
``` 


### Create replicator user on the primary

create variable in memory to make thing easier
```
_ldap_pass="$(cat /etc/ldap.secret)"
_dn_admin="cn=admin,dc=co,dc=my10c,dc=com"
```

Create these 2 files with the content show below

file replicator.ldif 
note that the password is not set, this will be change later
```
dn: cn=replicator,dc=co,dc=my10c,dc=com
objectClass: simpleSecurityObject
objectClass: organizationalRole
cn: replicator
description: Replication user
userPassword: {CRYPT}x
```

file replicator_acl.ldif
```

dn: olcDatabase={1}mdb,cn=config
changetype: modify
add: olcAccess
olcAccess: {0}to *
  by dn.exact="cn=replicator,dc=co,dc=my10c,dc=com" read
  by * break
-
add: olcLimits
olcLimits: dn.exact="cn=replicator,dc=co,dc=my10c,dc=com"
  time.soft=unlimited time.hard=unlimited
  size.soft=unlimited size.hard=unlimited
then execute this commands
```

These will:
- setup the user *replicator*
- change the password of the user *replicator*
- setup the access for the *replicator* users (from the file replicator.pass)
```
ldapadd -w $_ldap_pass -x -ZZ -D $_dn_admin -f replicator.ldif

ldappasswd -H ldap://localhost -w $_ldap_pass -x -D $_dn_admin \
        -T replicator.pass cn=replicator,dc=co,dc=my10c,dc=com

ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f replicator_acl.ldif
```

### Create the consumer file
this will be use on the secondaries servers (replicas)

file consumer.ldif content
```
dn: cn=module{0},cn=config
changetype: modify
add: olcModuleLoad
olcModuleLoad: syncprov

dn: olcDatabase={1}mdb,cn=config
changetype: modify
add: olcDbIndex
olcDbIndex: entryUUID eq
-
add: olcSyncrepl
olcSyncrepl: rid=0
  provider=ldap://ldap-primary.ops.my10c.com
  bindmethod=simple
  binddn="cn=replicator,dc=co,dc=my10c,dc=com"
  credentials=XXXXXXX <- change to the replicatot password!
  searchbase="dc=co,dc=my10c,dc=com"
  schemachecking=on
  type=refreshAndPersist
  starttls=critical tls_reqcert=demand
  retry="60 +"
-
add: olcUpdateRef
olcUpdateRef: ldap://ldap-primary.co.my10c.com
On the secondary load the configuration and restart the service
```

```
ldapadd -c -Q -Y EXTERNAL -H ldapi:/// -f consumer.ldif
systemctl restart slapd
```

### Check replication

Check replication by this command on both the primary and secondary, 
the result should show same values once the servers are in sync

```
ldapsearch -z1 -LLLQY EXTERNAL -H ldapi:/// -s base -b dc=co,dc=my10c,dc=com contextCSN
```

check TLS
```
ldapwhoami -H ldap://localhost -x -ZZ
```

## The End
Congraculation you should be all set now : 🦄 👏
