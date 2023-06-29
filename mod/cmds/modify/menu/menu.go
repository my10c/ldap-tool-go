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

	modifyGroup "my10c.ldap/cmds/modify/group"
	modifySudo "my10c.ldap/cmds/modify/sudo"
	modifyUser "my10c.ldap/cmds/modify/user"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

func ModifyMenu(conn *ldap.Connection, funcs *vars.Funcs) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Modify", 20, true))
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
		modifyUser.Modify(conn, funcs)
	case "group", "g":
		modifyGroup.Modify(conn, funcs)
	case "sudo", "s":
		modifySudo.Modify(conn, funcs)
	case "quit", "q":
		funcs.P.PrintRed("\n\tOperation cancelled\n")
		fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 40))
		vars.QuitTool = true
		break
	default:
		modifyUser.Modify(conn, funcs)
	}
}
