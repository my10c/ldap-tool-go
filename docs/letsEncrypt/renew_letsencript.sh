#!/usr/bin/env bash
#
# BSD 3-Clause License
#
# Copyright (c) 2020, © Badassopc LLC
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
#
# * Redistributions of source code must retain the above copyright notice, this
#   list of conditions and the following disclaimer.
#
# * Redistributions in binary form must reproduce the above copyright notice,
#   this list of conditions and the following disclaimer in the documentation
#   and/or other materials provided with the distribution.
#
# * Neither the name of the copyright holder nor the names of its
#   contributors may be used to endorse or promote products derived from
#   this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#*
#* File			:	renew_letsencript
#*
#* Description	:	pre script to renew the lets encrypt cert
#*
#* Author	:	Luc Suryo <luc@my10c.com>
#*
#* Version	:	0.1
#*
#* Date		:	May 6, 2020
#*
#* History	:
#* 		Date:			Author:				Info:
#*		May 6, 2020		LIS					First Release
#*

_program="${0##*/}"
_author='Luc Suryo'
_copyright="Copyright 2020 - $(date "+%Y") © Badassops LLC"
_license='License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥'
_version='0.1'
_email='luc@my10c.com'
_summary='pre script to renew the lets encrypt certs'
_cancelled="OK : Process has been cancelled on your request."
_info="$_program $_version\n$_copyright\n$_license\n\nWritten by $_author <$_email>\n$_summary\n"

# defined return values
_state_ok=0 ; _state_ok_msg='OK'
_state_critical=2 ; _state_critical_msg='CRITICAL'

# Whatever script needs to be run as root and exclusive lock
_need_root=0
_need_lock=0

# color :)
# Reset
_color_off='\033[0m'       # Text Reset

# Bold
_color_black='\033[1;30m'       # Black
_color_red='\033[1;31m'         # Red
_color_green='\033[1;32m'       # Green
_color_yellow='\033[1;33m'      # Yellow
_color_blue='\033[1;34m'        # Blue
_color_purple='\033[1;35m'      # Purple
_color_cyan='\033[1;36m'        # Cyan
_color_white='\033[1;37m'       # White

function _print_it() {
	local _info_mode=$1
	shift
	case $_info_mode
	in
		main)		printf "${_color_purple}%s${_color_off}\n" "$*" 2>&1 ;;
    	info)		printf "${_color_blue}\t => %s${_color_off}\n" "$*" 2>&1 ;;
    	ok)			printf "${_color_green}\t => %s${_color_off}\n" "$*" 2>&1 ;;
    	warning)	printf "${_color_yellow}\t => %s${_color_off}\n" "$*" 2>&1 ;;
    	error)		printf "${_color_red}\t => %s${_color_off}\n" "$*" 2>&1 ;;
    	help)		printf "${_color_cyan}" ; echo -e "$*" ; printf "${_color_off}" 2>&1 ;;
		ignore)		printf "${_color_purple}\t => %s${_color_off}\n" "$*" 2>&1 ;;
    	*)			printf "${_color_cyan}\t => %s${_color_off}\n" "$*" 2>&1 ;;
	esac
	return 0
}

# working variables
_opid=$$
_hostname="${HOSTNAME%%.*}"
_work_dir=/var/tmp/"$_program"
_lockfile="$_work_dir/$_program".LOCK
_echo_flag='-e'

# commands used

# certs files. mode and owenership
_source_dir="/etc/letsencrypt/archive/co.my10c.com"
_target_dir="/etc/ldap/certs"
_target_cert_mode="0444"
_target_key_mode="0400"
_target_owner="openldap:root"
_certs_rename="cert1.pem:cert.pem chain1.pem:ca.pem"
_key_rename="privkey1.pem:key.pem"
_service="slapd"

# Set interrupt handler
trap inthandler 1 2 3 4 5 9 10 12 15 23 24 25

function inthandler() {
	clean_up
	_print_it pw "$_cancelled"
	exit $_state_ok
}

function clean_up() {
	rm -f "$_lockfile" > /dev/null 2>&1
	return 0
}

function _print_it() {
	local _info_mode=$1
	shift
	case $_info_mode
	in
		main)       printf "${_color_purple}%s${_color_off}\n" "$*" 2>&1 ;;
		info)       printf "${_color_blue}\t => %s${_color_off}\n" "$*" 2>&1 ;;
		ok)         printf "${_color_green}\t => %s${_color_off}\n" "$*" 2>&1 ;;
		warning)    printf "${_color_yellow}\t => %s${_color_off}\n" "$*" 2>&1 ;;
		error)      printf "${_color_red}\t => %s${_color_off}\n" "$*" 2>&1 ;;
		help)       printf "${_color_cyan}" ; echo -e "$*" ; printf "${_color_off}" 2>&1 ;;
		ignore)     printf "${_color_purple}\t => %s${_color_off}\n" "$*" 2>&1 ;;
	*)          printf "${_color_cyan}\t => %s${_color_off}\n" "$*" 2>&1 ;;
	esac
return 0
}

function help() {
	trap 1 2 3 4 5 9 10 12 15 23 24 25
	echo $_echo_flag "$_info"
	echo $_echo_flag "Usage : $_program [-h] [ option ... ]"
	echo $_echo_flag " Options:"
	echo $_echo_flag " -h this help page."
	echo $_echo_flag " -v Show version."
	echo $_echo_flag "\n Notes:"
	clean_up
	exit $1
}

function get_given_options() {
	while [[ -n "$1" ]]
	do
		case "$1" in
			'-v'|'--version')		_print_it mi "$_version" ; exit $_state_ok ; ;;
			*)						help 0 ;; 	# Which includes -h and --help
		esac
		shift
	done
	return 0
}

function isRoot() {
	if (($(id -u) != 0)) ; then
		_print_it pe "$_program: this script must be run as the user root"
		return 1
	fi
	return 0
}

function check_running() {
	if [[ -f "$_lockfile" ]]; then
		/bin/ps -p $(cat "$_lockfile") > /dev/null 2>&1
		if (( $? == 0 )) ; then
			_print_it pw "There is already a $_program running, execution has been terminated"
			_print_it pw "If this is an error please remove the lock file: $_lockfile"
			exit $_state_ok
		else
			_print_it pw "Lock file found and deleted since there is no process with that pid"
			rm -rf "$_lockfile" > /dev/null 2>&1
		fi
	fi
	if ! mkdir -p "$_work_dir" > /dev/null 2>&1;then
		_print_it pi "Unable to create file working directory $_work_dir"
		exit $_state_unknown
	fi
	echo "$_opid" > "$_lockfile"
	return 0
}

function _install_new_cert() {
	local _result=0
	for _cert in $_certs_rename
	do
		_org_name=${_cert%%:*}
		_dest_name=${_cert##*:}

		cp $_source_dir/$_org_name $_target_dir/$_dest_name
		_result=$(($_result + $?))

		chmod $_target_cert_mode $_target_dir/$_dest_name
		_result=$(($_result + $?))

		chown $_target_owner $_target_dir/$_dest_name
		_result=$(($_result + $?))
	done
	for _key in $_key_rename
	do
		_org_name="${_key%%:*}"
		_dest_name="${_key##*:}"

		cp $_source_dir/$_org_name $_target_dir/$_dest_name
		_result=$(($_result + $?))

		chmod $_target_key_mode $_target_dir/$_dest_name
		_result=$(($_result + $?))

		chown $_target_owner $_target_dir/$_dest_name
		_result=$(($_result + $?))
	done
	return $_result
}

function renew_it() {
	local _result=0
	_install_new_cert
	_result=$?
	if (( $_result == 0 )) ; then
		systemctl restart $_service
		return $?
	fi
	return $_result
}

function main() {
	local _var_exit=$_state_ok

	get_given_options $@
	if (( $_need_root == 1 )) ; then
		 isRoot
		(( $? != 0 )) && echo $_echo_flag "$_info" && exit 255
	fi
	(( $_need_lock == 1 )) && check_running

	renew_it
	_var_exit=$?

	clean_up
	case $_var_exit
	in
		0)	_print_it info "new certificate has installed uccessfully" ;;
		*)	_print_it error "fail to install new certificated" ;;
	esac
	trap 1 2 3 4 5 9 10 12 15 23 24 25
	exit $_var_exit
}
main $@
