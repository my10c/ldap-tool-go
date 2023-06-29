#==============================================================================
# LTB Self Service Password
#
# Copyright (C) 2009 Clement OUDOT
# Copyright (C) 2009 LTB-project.org
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# GPL License: http://www.gnu.org/licenses/gpl.txt
#
#==============================================================================

#==============================================================================
# All the default values are kept here, you should not modify it but use
# config.inc.local.php file instead to override the settings from here.
#==============================================================================

#==============================================================================
# Configuration
#==============================================================================

# Debug mode
$debug = false;

# LDAP
$ldap_url = "ldap://localhost";
$ldap_starttls = false;
$ldap_binddn = ""; <---- ADJUST THIS
$ldap_bindpw = ""; <---- ADJUST THIS
$ldap_base = ""; <---- ADJUST THIS
$ldap_login_attribute = "uid";
$ldap_fullname_attribute = "cn";
$ldap_filter = "(&(objectClass=person)($ldap_login_attribute={login}))";

# Active Directory mode
$ad_mode = false;

# Samba mode
$samba_mode = false;

# Shadow options - require shadowAccount objectClass
# Update shadowLastChange
$shadow_options['update_shadowLastChange'] = true;
$shadow_options['update_shadowExpire'] = true;
# Default to -1, never expire
$shadow_options['shadow_expire_days'] = 90;

# Hash mechanism for password:
# SSHA, SSHA256, SSHA384, SSHA512
# SHA, SHA256, SHA384, SHA512
# SMD5
# MD5
# CRYPT
# clear (the default)
# auto (will check the hash of current password)
# This option is not used with ad_mode = true
$hash = "SSHA";

# Prefix to use for salt with CRYPT
$hash_options['crypt_salt_prefix'] = "$6$";
$hash_options['crypt_salt_length'] = "6";

# Rate limiting
$use_ratelimit = false;
$max_attempts_per_user = 10;
$max_attempts_per_ip = 20;
$max_attempts_block_seconds = "60";

# Header to use for client IP (HTTP_X_FORWARDED_FOR ?)
$client_ip_header = 'REMOTE_ADDR';

# Local password policy
# This is applied before directory password policy
$pwd_min_length = 16;
$pwd_max_length = 25;
$pwd_min_lower = 2;
$pwd_min_upper = 2;
$pwd_min_digit = 2;
$pwd_min_special = 1;
$pwd_special_chars = "@#$%^&*?.><";
$pwd_no_reuse = true;
$pwd_diff_login = true;
$pwd_diff_last_min_chars = 8;
$pwd_forbidden_words = array();
$pwd_forbidden_ldap_fields = array();
$pwd_complexity = 0;
$use_pwnedpasswords = false;
$pwd_show_policy = "always";
$pwd_show_policy_pos = "above";
$pwd_no_special_at_ends = false;

# Who changes the password?
$who_change_password = "manager";
$use_change = true;

# SSH Key Change
$change_sshkey = true;
$change_sshkey_attribute = "sshPublicKey";
$who_change_sshkey = "manager";
$notify_on_sshkey_change = false;

# Questions/answers
$use_questions = true;
$multiple_answers = true;
$multiple_answers_one_str = true;
$answer_objectClass = "extensibleObject";
$answer_attribute = "info";
$crypt_answers = true;

$messages['questions']['ice'] = "What is your favorite ice cream flavor?";
$messages['questions']['food'] = "What is your favorite food?";
$messages['questions']['nfl'] = "Who is your favorite NFL team?";
$messages['questions']['nba'] = "Who is your favorite NBA team?";
$messages['questions']['animal'] = "What is your favorite animal?";
$messages['questions']['pet'] = "What is the name of your first pet?";
$messages['questions']['highschool'] = "In which city did you went to highschool?";
$messages['questions']['college'] = "In which city did you went to college?";
$messages['questions']['mother'] = "What is your mother's name?";
$messages['questions']['father'] = "What is your father's name?";
$messages['questions']['coffee'] = "How do you like your coffee?";
$messages['questions']['tea'] = "How do you like your tea?";
$messages['questions']['wine'] = "What is your favorite wine?";
$messages['questions']['liqoure'] = "What is your favorite liqoure?";
$questions_count = 2;
$question_populate_enable = true;

# Token
$use_tokens = true;
$crypt_tokens = true;
$token_lifetime = "3600";

# SMTP settings
$mail_attribute = "mail";
$mail_address_use_ldap = false;
$mail_from = ""; <---- ADJUST THIS
$mail_from_name = "Self Service Password";
$mail_signature = ""; <---- ADJUST THIS
$notify_on_change = false;
$mail_sendmailpath = '/usr/sbin/sendmail';
$mail_protocol = 'smtp';
$mail_smtp_debug = 5;
$mail_debug_format = 'error_log';
$mail_smtp_host = 'localhost';
$mail_smtp_auth = false;
$mail_smtp_user = '';
$mail_smtp_pass = '';
$mail_smtp_port = 25;
$mail_smtp_timeout = 30;
$mail_smtp_keepalive = false;
$mail_smtp_secure = '';
$mail_smtp_autotls = false;
$mail_contenttype = 'text/plain';
$mail_wordwrap = 0;
$mail_charset = 'utf-8';
$mail_priority = 3;
$mail_newline = PHP_EOL;

# SMS
$use_sms = false;

# Encryption, decryption keyphrase, required if $crypt_tokens = true
# Please change it to anything long, random and complicated, you do not have to remember it
# Changing it will also invalidate all previous tokens and SMS codes
$keyphrase = ""; <---- ADJUST THIS

# Misc
$show_help = true;
$lang = "en";
$allowed_lang = array();
$show_menu = true;
$logo = "images/"; <---- ADJUST THIS
$background_image = "images/"; <---- ADJUST THIS
$custom_css = "";
$display_footer = false;
$reset_request_log = "/var/log/self-service-password";
$login_forbidden_chars = "*()&|";

# Captcha from Google
$use_recaptcha = true;
$recaptcha_publickey = ""; <---- ADJUST THIS
$recaptcha_privatekey = ""; <---- ADJUST THIS
$recaptcha_theme = "light";
$recaptcha_type = "image";
$recaptcha_size = "normal";
$recaptcha_request_method = null;

# API
$use_restapi = false;

# Cache directory
$smarty_compile_dir = "/var/cache/self-service-password/smarty_compile";
$smarty_cache_dir = "/var/cache/self-service-password/smarty_cache";

# Allow to override current settings with local configuration
if (file_exists (__DIR__ . '/config.inc.local.php')) {
    require __DIR__ . '/config.inc.local.php';
}

# Smarty
if (!defined("SMARTY")) {
    define("SMARTY", "/usr/share/php/smarty3/Smarty.class.php");
}

# Set preset login from HTTP header $header_name_preset_login
$presetLogin = "";
if (isset($header_name_preset_login)) {
    $presetLoginKey = "HTTP_".strtoupper(str_replace('-','_',$header_name_preset_login));
    if (array_key_exists($presetLoginKey, $_SERVER)) {
        $presetLogin = preg_replace("/[^a-zA-Z0-9-_@\.]+/", "", filter_var($_SERVER[$presetLoginKey], FILTER_SANITIZE_STRING, FILTER_FLAG_STRIP_HIGH));
    }
}

# Allow to override current settings with an extra configuration file, whose reference is passed in HTTP_HEADER $header_name_extra_config
if (isset($header_name_extra_config)) {
    $extraConfigKey = "HTTP_".strtoupper(str_replace('-','_',$header_name_extra_config));
    if (array_key_exists($extraConfigKey, $_SERVER)) {
        $extraConfig = preg_replace("/[^a-zA-Z0-9-_]+/", "", filter_var($_SERVER[$extraConfigKey], FILTER_SANITIZE_STRING, FILTER_FLAG_STRIP_HIGH));
        if (strlen($extraConfig) > 0 && file_exists (__DIR__ . "/config.inc.".$extraConfig.".php")) {
            require  __DIR__ . "/config.inc.".$extraConfig.".php";
        }
    }
}
