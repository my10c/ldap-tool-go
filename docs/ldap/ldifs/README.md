# MAKE SURE TO ADJUST THE FILES BEFORE USE !!!

## openssh.ldif and sudo.ldif
```
cp openssh.ldif sudo.ldif /etc/ldap/schema
chmod 0444 /etc/ldap/schema/{openssh.ldif,sudo.ldif}
chown root:root /etc/ldap/schema/{openssh.ldif,sudo.ldif}
```
