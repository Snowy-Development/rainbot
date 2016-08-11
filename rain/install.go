// Copyright (C) 2015  Rodolfo Castillo-Valladares & Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// Send any inquiries you may have about this program to: rcvallada@gmail.com

package rain

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func doInstall(install string, name string) error {
	// Install Core Dependencies
	err := GetDepends(CoreDepends)
	if err != nil {
		return err
	}

	// Install Specific Dependencies
	switch install {
	case "neverfree":
		err = GetDepends(NeverfreeDepends)
	}

	if err != nil {
		return err
	}

	var path string
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		path = gopath + "/bin/" + name
	} else if goroot := os.Getenv("GOROOT"); goroot != "" {
		path = goroot + "/bin/" + name
	} else {
		fmt.Println(" No GOROOT or GOPATH set, aborting install")
		return errors.New("No GOROOT or GOPATH set")
	}

	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}

	output, err := exec.Command("go", "build", "-o", path,
		"github.com/raindevteam/rain/install/"+install).CombinedOutput()
	if err != nil {
		fmt.Printf(" Internal command error\n ----\n%s\n -----\n\n", string(output[:]))
		return err
	}

	return nil
}
