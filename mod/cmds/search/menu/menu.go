//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	searchGroup "my10c.ldap/cmds/search/group"
	searchSudo "my10c.ldap/cmds/search/sudo"
	searchUser "my10c.ldap/cmds/search/user"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

func SearchMenu(conn *ldap.Connection, funcs *vars.Funcs) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Search", 20, true))
	fmt.Printf("\tSearch (%s)ser, (%s)ll Users, (%s)roup, all Group(%s)\n",
		funcs.P.MessageGreen("U"),
		funcs.P.MessageGreen("A"),
		funcs.P.MessageGreen("G"),
		funcs.P.MessageGreen("S"),
	)
	fmt.Printf("\t\t(%s)sudo role, (%s)all sudos role or (%s)uit?\n\t(default to User)? choice: ",
		funcs.P.MessageBlue("X"),
		funcs.P.MessageBlue("Z"),
		funcs.P.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		searchUser.User(conn, funcs)
	case "users", "a":
		searchUser.Users(conn, funcs)
	case "group", "g":
		searchGroup.Group(conn, funcs)
	case "groups", "s":
		searchGroup.Groups(conn, funcs)
	case "sudo", "x":
		searchSudo.Sudo(conn, funcs)
	case "sudos", "z":
		searchSudo.Sudos(conn, funcs)
	case "quit", "q":
		funcs.P.PrintRed("\n\t\tOperation cancelled\n")
		fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 40))
		vars.QuitTool = true
		break
	default:
		searchUser.User(conn, funcs)
	}
}
