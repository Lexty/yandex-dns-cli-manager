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

	"errors"
	"strings"

	"github.com/lexty/yandex-dns-cli-manager/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		rec := api.Record{}

		rec.RecordId = viper.GetInt("id")
		rec.RecordType = viper.GetString("type")
		rec.AdminMail = viper.GetString("admin-mail")
		rec.Content = viper.GetString("content")
		rec.Priority = viper.GetString("priority")
		rec.Weight = viper.GetInt("weight")
		rec.Port = viper.GetInt("port")
		rec.Target = viper.GetString("target")
		rec.Subdomain = viper.GetString("subdomain")
		rec.TTL = viper.GetInt("ttl")
		rec.Refresh = viper.GetInt("refresh")
		rec.Retry = viper.GetInt("retry")
		rec.Expire = viper.GetInt("expire")
		rec.NegCache = viper.GetInt("neg-cache")

		resp, err := api.EditRecord(&rec, viper.GetString("domain"), viper.GetString("admin-token"))

		if err != nil {
			throwError(err)
		}

		switch viper.GetString("format") {
		case formatJson:
			fmt.Print(resp.Json)
		case formatList:
			setProps()
			fmt.Printf("Record successfully changed\n\n")
			printList(filterRecords([]api.Record{resp.Record}, []string{"*"}), strings.Join([]string{propId, propType, propContent, propSubdomain, propPriority, propTTL, propFQDN}, ","))
		default:
			throwError(errors.New(fmt.Sprintf(`Unknown output format "%s".`, viper.GetString("format"))))
		}
	},
}

func init() {
	RootCmd.AddCommand(editCmd)

	editCmd.Flags().StringP("format", "f", "", fmt.Sprintf("format output (%s|%s)", formatList, formatJson))
	viper.BindPFlag("format", editCmd.Flags().Lookup("format"))
	viper.SetDefault("format", formatList)

	editCmd.Flags().IntP("id", "i", 0, "ID of the record")
	viper.BindPFlag("id", editCmd.Flags().Lookup("id"))

	editCmd.Flags().StringP("admin-mail", "m", "", "email-address of the domain's administrator")
	viper.BindPFlag("admin-mail", editCmd.Flags().Lookup("admin-mail"))

	editCmd.Flags().StringP("content", "c", "", "content of the DNS record")
	viper.BindPFlag("content", editCmd.Flags().Lookup("content"))

	editCmd.Flags().StringP("priority", "p", "", "priority of the DNS record")
	viper.BindPFlag("priority", editCmd.Flags().Lookup("priority"))

	editCmd.Flags().IntP("weight", "w", 0, "weight of the SRV-record relative to other SRV-records for the same domain with the same priority")
	viper.BindPFlag("weight", editCmd.Flags().Lookup("weight"))

	editCmd.Flags().StringP("port", "P", "", "TCP or UDP port of the host that is hosting the service")
	viper.BindPFlag("port", editCmd.Flags().Lookup("port"))

	editCmd.Flags().StringP("target", "T", "", "the canonical name of the host providing the service")
	viper.BindPFlag("target", editCmd.Flags().Lookup("target"))

	editCmd.Flags().StringP("subdomain", "s", "", "name of the subdomain")
	viper.BindPFlag("subdomain", editCmd.Flags().Lookup("subdomain"))

	editCmd.Flags().IntP("ttl", "l", 0, "the lifetime of the DNS record in seconds")
	viper.BindPFlag("ttl", editCmd.Flags().Lookup("ttl"))

	editCmd.Flags().IntP("refresh", "r", 0, "time between updates")
	viper.BindPFlag("refresh", editCmd.Flags().Lookup("refresh"))

	editCmd.Flags().IntP("retry", "R", 0, "the time between attempts to obtain records")
	viper.BindPFlag("retry", editCmd.Flags().Lookup("retry"))

	editCmd.Flags().IntP("expire", "e", 0, "time limit")
	viper.BindPFlag("expire", editCmd.Flags().Lookup("expire"))

	editCmd.Flags().IntP("neg-cache", "n", 0, "caching time")
	viper.BindPFlag("neg-cache", editCmd.Flags().Lookup("neg-cache"))

}
