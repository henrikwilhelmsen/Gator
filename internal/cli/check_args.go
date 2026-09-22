// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import "fmt"

// checkArgs checks if the given cmd args matches the wanted number of arguments and
// argument names.
func checkArgs(wantLen int, wantNames []string, cmd Command) error {
	if wantLen == 0 && len(cmd.Args) != 0 {
		return fmt.Errorf(
			"command '%s' expects no arguments, got %d",
			cmd.Name,
			len(cmd.Args),
		)
	}
	if len(cmd.Args) != wantLen {
		return fmt.Errorf(
			"command '%s' expects %d arguments (%v), got %d",
			cmd.Name, wantLen, wantNames, len(cmd.Args),
		)
	}

	return nil
}
