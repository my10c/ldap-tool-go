# ADDING NEW SCHEMA FOR openssh and sudo
The openssh schema will allow one to store SSH public for a user in LDAP 
the sudo schema will allow to store SUDO rules in LDAP 
Note
- the LDAP sudo schema has limited feature
- all basic sudo rule will work
- read about how to set this up at the [SUDO web site](https://www.sudo.ws/docs/man/1.8.17/sudoers.ldap.man)

### files openssh.schema and sudo.schema

```
cp openssh.ldif sudo.ldif /etc/ldap/schema
chmod 0444 /etc/ldap/schema/{openssh.schema,sudo.schema}
chown root:root /etc/ldap/schema/{openssh.schema,sudo.schema}
```
