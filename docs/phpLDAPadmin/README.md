## phpLDAPAdmin

### Background
This document lay out how to install phpLDAPAdmin, a web interface to manage ldap users

### Technology
- running Ubuntu 20.04 or newer
- OpenLDAP as the LDAP software
- Nginx, web server to host the LDAP admin web interface

### Installation
- Install the web admin 
```
apt install phpldapadmin
```

- add more valid shells
```
vi /etc/phpldapadmin/templates/creation/posixAccount.xml
   and add the bash and zsh
```

- adjust the main configuration file
```
cp /etc/phpldapadmin/config.php /etc/phpldapadmin/config.php
chmod 0000 /etc/phpldapadmin/config.php
chown root:root /etc/phpldapadmin/config.php
adjust the values in /etc/phpldapadmin/config.php
```
  you can find and example of my configuration [here](https://github.com/my10c/ldap-tool-go/blob/main/docs/phpLDAPadmin/example-config.php) (minus secrets 😈 )

- setup with proper permissions and owner
```
chmod 0444 /etc/phpldapadmin/config.php
chown root:root /etc/phpldapadmin/config.php
```

- Install the Web server
  This should be done if nginx will be replacing the apache web server, 
  setup php-fpm and nginx to serve the phpLDAPadmin application is beyond this scope
  you can find and example of my nginx configuration [here](https://github.com/my10c/ldap-tool-go/blob/main/docs/phpLDAPadmin/example-nginx-phpldapadmin.conf)
  you can find and example of my ssl configuration [here](https://github.com/my10c/ldap-tool-go/blob/main/docs/phpLDAPadmin/example-nginx-ssl.conf)
```
apt remove apache2 apache2-bin apache2-data
apt install nginx
```

#### Note
- sometime you **might need to replace the file /usr/share/phpldapadmin/lib/functions.php** that comes with the phpldapadmin packages due bug.
 you can find the one I use [here](https://github.com/my10c/ldap-tool-go/blob/main/docs/phpLDAPadmin/fixed-functions.php) 
- the web interface **should** be installed on the primary LDAP server as it needs **read-write**

### The End
Your friendly BOFH 🦄  😈

