## LDAP-TOOL
<span style=color:red;>INFO</span>
Simple tool written in Go to manage OpenLDAP users and groups

### Background
The script is based on a certain LDAP settings
- OpenLDAP
- the memberOf ldap plugin
- the SSH schema
- the SUDO schema
- password length and use of special charachter is in the config file
- the config file is toml formatted
- Runs on OSX or Linux (Ubuntu 20.04 or newer)

### History
Using an UI interface such a phpLDAPadmin is not always possible, and so I decide 
to build this tools. Orignally it was written in **bash** using the ldap CLI's and then
in **Python**, but it was not always working on someone else's laptop do the specific command
(*bash*) or module (*Python*), so I decided to write it in Go

### Capabilities
The script to be able manage OpenLDAP users and groups:
- add, modify, search and delete an user or a group
- if the admin define is an user, it will be able to get the user's record (search) 
  and modify the user's password (modify)

## Usage
```
usage: ldap-tool [-h|--help] [-c|--configFile "<value>"] [-s|--server
                 "<value>"] [-C|--command (create|modify|delete|search)]
                 [-d|--debug] [-i|--info] [-v|--version]

                 Simple script to manage LDAP users, groups and SUDO rules

Arguments:

  -h  --help        Print help information
  -c  --configFile  Path to the configuration file to be use. Default:
                    /usr/local/etc/ldap-tool/ldap-tool.ini
  -s  --server      Server profile name. Default: default
  -C  --command     commands: search, create,
                    modify, delete. Default: search
  -d  --debug       Enable debug. Default: false
  -i  --info        Show information
  -v  --version     Show version
```

### Build or run the code the code
To build the code into a single binary as simple as
```
go build ldap-tools.go
```
If everything is well, then this will produce a binary called **ldap-tools** 

To run the code
```
go run ldap-tools.go -c <your-config-file> -C <command>
```

### Screen shots

[screenshots](https://github.com/my10c/ldap-tool-go/blob/main/docs/screenshots) 


### TODO / *wishlist*
 - redis for **delete of user only**
	- start : add key in redis
	- end   : delete key in redis
	- safe guard : run a crontab that will GET KEYS
		then for every KEY, delete the user in ldap
		and DEL KEY in redis once the user has been deleted
		so in case one loses the network connection, we have a safe guard
		that the user will be deleted
	- if redis it not available write the key on disk ?

#### IDEA:
 let me know your request 👻  and I *might* add it 😎

### The End
Your friendly BOFH 🦄 😈          
