// Copyright © 2015 Alexandr Medvedev <alexandr.mdr@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"io/ioutil"
	"log"
	"os/user"
	"strings"

	"os"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newAdminToken string
var newDomain string
var newProps string

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show or change settings",
	Run: func(cmd *cobra.Command, args []string) {
		if "" == newAdminToken && "" == newDomain && "" == newProps {
			printSettings()
		} else {
			saveSettings()
		}
	},
}

func printSettings() {
	props := viper.GetString("props")
	if props == "" {
		props = propsDefault
	}

	fmt.Printf(`Settings:
	admin-token %s
	domain      %s
	props       %s
`, viper.GetString("admin-token"), viper.GetString("domain"), props)
}

func saveSettings() {
	if "" == newAdminToken && viper.IsSet("admin-token") {
		newAdminToken = viper.GetString("admin-token")
	}
	if "" == newDomain && viper.IsSet("domain") {
		newDomain = viper.GetString("domain")
	}

	var lines []string
	if newAdminToken != "" {
		lines = append(lines, `    "admin-token": "`+newAdminToken+`"`)
	}
	if newDomain != "" {
		lines = append(lines, `    "domain": "`+newDomain+`"`)
	}
	if newProps != "" {
		lines = append(lines, `    "props": "`+newProps+`"`)
	}

	fmt.Println(newAdminToken)

	cfgContent := "{\n" + strings.Join(lines, ",\n") + "\n}"

	if cfgFile == "" {
		cfgFile = getDefaultCfgFilepath()
	}
	if err := ioutil.WriteFile(cfgFile, []byte(cfgContent), 0600); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Settings successfully changed in \"%s\"\n", cfgFile)
}

func getDefaultCfgFilepath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + string(os.PathSeparator) + cfgFileName + "." + cfgFileType
}

func init() {
	RootCmd.AddCommand(settingsCmd)
	settingsCmd.Flags().StringVarP(&newAdminToken, "admin-token", "a", "", "set your admin token")
	settingsCmd.Flags().StringVarP(&newDomain, "domain", "d", "", "set domain name")
	settingsCmd.Flags().StringVarP(&newProps, "props", "p", "", "set default output record properties")
}
