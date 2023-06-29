// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:	0.1
//

package main

import (
	"fmt"
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	// local
	"my10c.ldap/configurator"
	"my10c.ldap/initializer"
	"my10c.ldap/ldap"
	"my10c.ldap/logs"
	"my10c.ldap/vars"

	// on github
	"github.com/my10c/packages-go/lock"
	"github.com/my10c/packages-go/spinner"

	// the menus
	createMenu "my10c.ldap/cmds/create/menu"
	deleteMenu "my10c.ldap/cmds/delete/menu"
	limit "my10c.ldap/cmds/limit"
	modifyMenu "my10c.ldap/cmds/modify/menu"
	searchMenu "my10c.ldap/cmds/search/menu"
)

func main() {
	var enteredChoice string
	LockPid := os.Getpid()
	progName, _ := os.Executable()
	progBase := filepath.Base(progName)

	s := spinner.New(10)

	config := configurator.Configurator()

	// get given parameters
	config.InitializeArgs()

	// get the configuration
	config.InitializeConfigs()

	// initialize the user data dictionary
	funcs := initializer.Init(config)

	// make sure the configuration file has the proper settings
	runningUser, _ := funcs.I.IsRunningUser()
	if !funcs.I.IsInList(config.AuthValues.AllowUsers, runningUser) {
		funcs.P.PrintRed(fmt.Sprintf("The program has to be run as these user(s): %s or use sudo, aborting..\n",
			strings.Join(config.AuthValues.AllowUsers[:], ", ")))
		os.Exit(0)
	}
	ownerInfo, ownerOK := funcs.I.IsFileOwner(config.ConfigFile, config.AuthValues.AllowUsers)
	if !ownerOK {
		funcs.P.PrintRed(fmt.Sprintf("%s,\nAborting..\n", ownerInfo))
		os.Exit(0)
	}
	permInfo, permOK := funcs.I.IsFilePermission(config.ConfigFile, config.AuthValues.AllowMods)
	if !permOK {
		funcs.P.PrintRed(fmt.Sprintf("%s,\nAborting..\n", permInfo))
		os.Exit(0)
	}

	go s.Run()
	// initialize the logger system
	LogConfig := &logs.LogConfig{
		LogsDir:       config.LogValues.LogsDir,
		LogFile:       config.LogValues.LogFile,
		LogMaxSize:    config.LogValues.LogMaxSize,
		LogMaxBackups: config.LogValues.LogMaxBackups,
		LogMaxAge:     config.LogValues.LogMaxAge,
	}

	logs.InitLogs(LogConfig)
	logs.Log("System all clear", "INFO")

	// prevent a race
	time.Sleep(1 * time.Second)
	s.Stop()

	// check if server was set to allow read-write
	if config.ServerValues.ReadWrite == false {
		funcs.P.PrintRed(
			fmt.Sprintf("\n\tThe server %s is set to be ready only.\n\tOnly the Search option is available...\n",
				config.ServerValues.Server))
		funcs.P.PrintGreen("\tPress enter to continue to search: ")
		fmt.Scanln()
		config.Cmd = "search"
	}

	// create the lock file to prevent an other script is running/started
	lockPtr := lock.New(config.DefaultValues.LockFile)
	config.LockPID = LockPid
	if config.Cmd != "search" {
		// check lock file; lock file should not exist
		if _, fileExist, _ := funcs.I.IsExist(config.DefaultValues.LockFile, "file"); fileExist {
			lockPid, _ := lockPtr.LockGetPid()
			if progRunning, _ := funcs.I.IsRunning(progBase, lockPid); progRunning {
				funcs.P.PrintRed(fmt.Sprintf("\nError there is already a process %s running, aborting...\n", progBase))
				os.Exit(0)
			}
		}
		// save to create new or overwrite the lock file
		if err := lockPtr.LockIt(LockPid); err != nil {
			funcs.P.PrintRed(fmt.Sprintf("\nError creating the lock file, error %s, aborting..\n", err.Error()))
			os.Exit(0)
		}
	}

	// start the LDAP connection
	conn := ldap.New(config)

	for true {
		// semi-hardcoded
		if config.ServerValues.Admin != "cn=admin," +config.ServerValues.BaseDN {
			switch config.Cmd {
			case "search":
				limit.UserRecord(conn, funcs)
			case "modify":
				limit.ModifyUserPassword(conn, funcs)
			default:
				funcs.P.PrintRed("\n\tThis command is only available for admin...\n\n")
			}
		} else {
			switch config.Cmd {
			case "search":
				searchMenu.SearchMenu(conn, funcs)
			case "create":
				createMenu.CreateMenu(conn, funcs)
			case "modify":
				modifyMenu.ModifyMenu(conn, funcs)
			case "delete":
				deleteMenu.DeleteMenu(conn, funcs)
			}
		}
		if vars.QuitTool == true {
			break
		}
		fmt.Printf("\tEnter %s(uit), or %s(earch) %s(reate), %s(odify), %s(elete) to switch mode\n",
			funcs.P.MessageRed("Q"), funcs.P.MessageYellow("S"), funcs.P.MessageYellow("C"),
			funcs.P.MessageYellow("M"), funcs.P.MessageYellow("D"))
		fmt.Printf("\tPress enter to performance an other %s : ", config.Cmd)
		reader := bufio.NewReader(os.Stdin)
		enteredChoice, _ = reader.ReadString('\n')
	    enteredChoice = strings.TrimSuffix(enteredChoice, "\n")
		enteredChoice = strings.ToUpper(enteredChoice)
		if enteredChoice == "Q" {
			fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
			break
		}
		switch enteredChoice {
			case "S": config.Cmd = "search"
			case "C": config.Cmd = "create"
			case "M": config.Cmd = "modify"
			case "D": config.Cmd = "delete"
			default: break
		}
	}

	if config.Cmd != "search" {
		lockPtr.LockRelease()
	}

	funcs.P.TheEnd()
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	logs.Log("System Normal shutdown", "INFO")
	os.Exit(0)
}
