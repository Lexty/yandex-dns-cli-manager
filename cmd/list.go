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
	"fmt"

	"os"

	"strconv"

	"strings"

	"errors"

	"github.com/Lexty/yandexdns/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	propId        string = "id"
	propSubdomain string = "sub"
	propType      string = "type"
	propContent   string = "content"
	propPriority  string = "prior"
	propTtl       string = "ttl"
)

var props map[string]string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "The list of records in the domain zone",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("admin-token") {
			fmt.Println("Error: --admin-token is not set")
			os.Exit(-1)
		}
		if !viper.IsSet("domain") {
			fmt.Println("Error: --domain is not set")
			os.Exit(-1)
		}
		list, err := api.GetList(viper.GetString("domain"), viper.GetString("admin-token"))
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(-1)
		}
		setProps()
		printTable(list.Records, viper.GetString("props"))

		var data [][]string

		for _, r := range list.Records {
			data = append(data, []string{strconv.Itoa(r.RecordId), r.Subdomain, r.RecordType, r.Content, fmt.Sprintf("%v", r.Priority)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "Subdomain", "Type", "Content", "Priority"})

		for _, v := range data {
			table.Append(v)
		}
		//		table.Render()
	},
}

func parseProps (raw string) []string {
	parts := strings.Split(raw, ",")
	for i, part := range parts {
		parts[i] = strings.Trim(part, " ")
	}
	return parts
}

func setProps() {
	props = make(map[string]string, 6)
	props[propId] = "Id"
	props[propSubdomain] = "Subdomain"
	props[propType] = "Type"
	props[propContent] = "Content"
	props[propTtl] = "Ttl"
	props[propPriority] = "Priority"
}

func printTable(r []api.Record, props string) {

}

func getHeader(prop string) (string, error) {
	if title, ok := props[strings.ToLower(prop)]; ok {
		return title, nil
	} else {
		return "", errors.New(fmt.Sprintf(`Unknown record property %s.`, prop))
	}
}
func getValue(prop string, r *api.Record) (string, error) {
	var val string
	switch prop {
	case propId:
		val = strconv.Itoa(r.RecordId)
	case propSubdomain:
		val = r.Subdomain
	case propType:
		val = r.RecordType
	case propContent:
		val = r.Content
	case propTtl:
		val = strconv.Itoa(r.Ttl)
	case propPriority:
		val = fmt.Sprintf("%v", r.Priority)
	default:
		return "", errors.New(fmt.Sprintf(`Unknown record property %s.`, prop))
	}
	return val, nil
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("json", "j", false, "set output in JSON format")
	viper.BindPFlag("json", listCmd.Flags().Lookup("json"))
	viper.SetDefault("json", false)

	listCmd.Flags().StringP("admin-token", "a", "", "admin's token")
	viper.BindPFlag("admin-token", listCmd.Flags().Lookup("admin-token"))

	listCmd.Flags().StringP("domain", "d", "", "domain name")
	viper.BindPFlag("domain", listCmd.Flags().Lookup("domain"))

	listCmd.Flags().StringP("format", "f", "", "format output (table|list)")
	viper.BindPFlag("format", listCmd.Flags().Lookup("props"))
	viper.SetDefault("format", "table")

	listCmd.Flags().StringP("props", "p", "", "record poroperties for display (does not work for json format)")
	viper.BindPFlag("props", listCmd.Flags().Lookup("props"))
	viper.SetDefault("props", strings.Join([]string{propId, propSubdomain, propType, propContent, propPriority}, ","))
}
