<?php

$config->custom->appearance['tree'] = 'AJAXTree';
$config->custom->appearance['hide_template_warning'] = true;
$config->custom->appearance['minimalMode'] = true;
$config->custom->session['blowfish'] = 'XXXXX'; <--- REPLACE THIS

$config->custom->appearance['friendly_attrs'] = array(
    'facsimileTelephoneNumber' => 'Fax',
    'gid'                      => 'Group',
    'mail'                     => 'Email',
    'telephoneNumber'          => 'Telephone',
    'uid'                      => 'User Name',
    'userPassword'             => 'Password'
);

$config->custom->appearance['attr_display_order'] = array(
   'uid',
   'sn',
   'cn',
   'givenName',
   'displayName',
   'uidNumber',
   'gidNumber',
   'homeDirectory',
   'loginShell',
   'userPassword'
 );

/*********************************************
 * Define your LDAP servers in this section  *
 *********************************************/
$servers = new Datastore();
$servers->newServer('ldap_pla');
$servers->setValue('server','name','XXXX'); <--- REPLACE THIS
$servers->setValue('server','host','127.0.0.1'); <--- ADJUST THIS IF NOT LOCALHOST
$servers->setValue('server','base',array('XXXXX')); <--- REPLACE THIS
$servers->setValue('login','auth_type','session');
$servers->setValue('login','bind_id','XXXXX'); <--- REPLACE THIS
/* Use TLS (Transport Layer Security) to connect to the LDAP server. */
// $servers->setValue('server','tls',false);
$servers->setValue('unique','attrs',array('mail','uid','uidNumber'));

?>
