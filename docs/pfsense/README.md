## SETUP PFSENSE WITH LDAP

One can setup a pfsense server to do authentication with an LDAP server
The screenshot wil be using this values as example:
- DN == `dc=co,dc=my10c,dc=com`
- Query == `memberOf=cn=vpn-co,ou=groups,dc=co,dc=my10c,dc=com`

**Important note**
Query : it is advisable to have it lock to a specific group of a `memberOF` group type
  this will only allow users to login if they are member of that group!

### LDAP STEPS

- login with Administration right
- select `system` from the top bar
- select `User Manager`
- select `Authentication Servers`
- click the `Add` button

On the Form:
- Descriptive name : `a useful text`
- Hostname or IP address : 
  - if the ldap server has a proper certificate and cert, then use the `hostname` of the ldap server
  - if the ldap server does not have a proper certificate and cert then use thee `ip` of the ldap server
- Search scope : `Entire Subtree`
- Base DN : your ldap server `base dn`
- Extended query : `checked` Enable extended query
- Query : `memberOf=cn=vpn-co,ou=groups,dc=co,dc=my10c,dc=com`
- Bind credentials : `the admin dn` and the proper password
- User naming attribute : `uid`
- Group naming attribute : `cn`
- Group member attribute : `memberOf`
- Group Object Class : `posixGroup`

Go back on **Authentication containers** and click `Select a container`
```
Select LDAP containers for authentication
Containers ou=users,dc=co,dc=my10c,dc=com
           ou=groups,dc=co,dc=my10c,dc=com
           ou=SUDOers,dc=co,dc=my10c,dc=com
```
- select these 2 containers
```
ou=users,dc=co,dc=my10c,dc=com
ou=groups,dc=co,dc=my10c,dc=com
```
Then hit `save`

Go bottom of the page and hit `save`

- [screenshot top form](https://github.com/my10c/ldap-tool-go/blob/main/docs/pfsense/Authentication-Servers-form-top.png) 
- [screenshot top bottom](https://github.com/my10c/ldap-tool-go/blob/main/docs/pfsense/Authentication-Servers-form-bottom.png) 
- [screenshot LDAP container](https://github.com/my10c/ldap-tool-go/blob/main/docs/pfsense/select-LDAP-Container.png) 


### GROUP STEPS

- login with Administration right
- select `system` from the top bar
- select `User Manager`
- select `Group`
- hit `Add`

On the Form
- Group name : ** must be the same name of the Query ** in the example this will be `vpn-co`
- Scope : `Remote`
- hit `Save`


### TEST LDAP CONNECTION
- select `Diagnostics` from the top bar
- selecvt `Authentication`
- Authentication Server : select the ldap server that was created
- Username : an ldap username
- Password : the user's ldap's password


if successful it will show a `green` dialog window and show `This user is a member of groups:`
if it had groups asigned to it.
- [screenshot success test authentication](https://github.com/my10c/ldap-tool-go/blob/main/docs/pfsense/Diagnostics-Authentication-success.png) 

if not successful: check your setting
- [screenshot fail test authentication](https://github.com/my10c/ldap-tool-go/blob/main/docs/pfsense/Diagnostics-Authentication-failed.png) 

### The End
Your friendly BOFH 🦄  😈
