//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package delete

import (
	"fmt"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

// delete a ldap record
func Delete(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Delete Group", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.GroupWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.GroupWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.GroupDisplayFieldID
	if common.GetObjectRecord(conn, true, "group", funcs) {
		common.DeleteObjectRecord(conn, vars.SearchResultData.SearchResult, "group", funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
