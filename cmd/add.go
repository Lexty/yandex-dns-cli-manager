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

	"strings"

	"errors"

	"github.com/lexty/yandex-dns-cli-manager/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := api.AddRecord(&rec, viper.GetString("domain"), viper.GetString("admin-token"))

		if err != nil {
			throwError(err)
		}

		switch viper.GetString("format") {
		case formatJson:
			fmt.Print(resp.Json)
		case formatList:
			setProps()
			fmt.Print("Record successfully created\n\n")
			printList(filterRecords([]api.Record{resp.Record}, []string{"*"}), strings.Join([]string{propId, propType, propContent, propSubdomain, propPriority, propTTL, propFQDN}, ","))
		default:
			throwError(errors.New(fmt.Sprintf(`Unknown output format "%s".`, viper.GetString("format"))))
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("format", "f", "", fmt.Sprintf("format output (%s|%s)", formatList, formatJson))
	viper.BindPFlag("format", addCmd.Flags().Lookup("format"))
	viper.SetDefault("format", formatList)

	addCmd.Flags().StringVarP(&rec.RecordType, "type", "t", "", fmt.Sprintf("type of record (available: %s)", strings.Join([]string{typeA, typeAAAA, typeCNAME, typeSRV, typeTXT, typeSOA, typeMX, typeNS}, ", ")))
	addCmd.Flags().StringVarP(&rec.AdminMail, "admin-mail", "m", "", "email-address of the domain's administrator")
	addCmd.Flags().StringVarP(&rec.Content, "content", "c", "", "content of the DNS record")
	//addCmd.Flags().StringVarP(rec.Priority.(*string), "priority", "p", "", "priority of the DNS record")
	addCmd.Flags().IntVarP(&rec.Weight, "weight", "w", 0, "weight of the SRV-record relative to other SRV-records for the same domain with the same priority")
	addCmd.Flags().IntVarP(&rec.Port, "port", "P", 0, "TCP or UDP port of the host that is hosting the service")
	addCmd.Flags().StringVarP(&rec.Target, "target", "T", "", "the canonical name of the host providing the service")
	addCmd.Flags().StringVarP(&rec.Subdomain, "subdomain", "s", "", "Name of the subdomain")
	addCmd.Flags().IntVarP(&rec.TTL, "ttl", "l", 0, "the lifetime of the DNS record in seconds")
}
