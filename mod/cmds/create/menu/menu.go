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

	createGroup "my10c.ldap/cmds/create/group"
	createSudo "my10c.ldap/cmds/create/sudo"
	createUser "my10c.ldap/cmds/create/user"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

func CreateMenu(conn *ldap.Connection, funcs *vars.Funcs) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Create", 20, true))
	fmt.Printf("\tModify (%s)ser, (%s)roup, (%s)udo rule or (%s)uit?\n\t(default to User)? choice: ",
		funcs.P.MessageGreen("U"),
		funcs.P.MessageGreen("G"),
		funcs.P.MessageBlue("S"),
		funcs.P.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		createUser.Create(conn, funcs)
	case "group", "g":
		createGroup.Create(conn, funcs)
	case "sudo", "s":
		createSudo.Create(conn, funcs)
	case "quit", "q":
		funcs.P.PrintRed("\n\t\tOperation cancelled\n")
		fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 40))
		break
	default:
		createUser.Create(conn, funcs)
	}
}
