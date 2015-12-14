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

	"github.com/Lexty/yandexdns/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rec := api.Record{}

		rec.RecordType = viper.GetString("type")
		rec.AdminMail = viper.GetString("admin-mail")
		rec.Content = viper.GetString("content")
		rec.Priority = viper.GetString("priority")
		rec.Weight = viper.GetInt("weight")
		rec.Port = viper.GetString("port")
		rec.Target = viper.GetString("target")
		rec.Subdomain = viper.GetString("subdomain")
		rec.TTL = viper.GetInt("ttl")

		resp, err := api.AddRecord(&rec, viper.GetString("domain"), viper.GetString("admin-token"))

		if err != nil {
			throwError(err)
		}

		switch viper.GetString("format") {
		case formatJson:
			fmt.Print(resp.Json)
		case formatList:
			setProps()
			fmt.Printf("Record successfully created\n\n")
			printList(filterRecords([]api.Record{resp.Record}, []string{"*"}), strings.Join([]string{propId, propType, propContent, propSubdomain, propPriority, propTTL, propFQDN}, ","))
		default:
			throwError(errors.New(fmt.Sprintf(`Unknown output format "%s".`, viper.GetString("format"))))
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("format", "f", "", fmt.Sprintf("format output (%s|%s)", formatList, formatJson))
	viper.BindPFlag("format", listCmd.Flags().Lookup("format"))
	viper.SetDefault("format", formatList)

	addCmd.Flags().StringP("type", "t", "", fmt.Sprintf("Type of record (available: %s)", strings.Join([]string{typeA, typeAAAA, typeCNAME, typeSRV, typeTXT, typeSOA, typeMX, typeNS}, ", ")))
	viper.BindPFlag("type", addCmd.Flags().Lookup("type"))

	addCmd.Flags().StringP("admin-mail", "m", "", "Email-address of the domain's administrator")
	viper.BindPFlag("admin-mail", addCmd.Flags().Lookup("admin-mail"))

	addCmd.Flags().StringP("content", "c", "", "Content of the DNS record")
	viper.BindPFlag("content", addCmd.Flags().Lookup("content"))

	addCmd.Flags().StringP("priority", "p", "", "Priority of the DNS record")
	viper.BindPFlag("priority", addCmd.Flags().Lookup("priority"))

	addCmd.Flags().IntP("weight", "w", 0, "Weight of the SRV-record relative to other SRV-records for the same domain with the same priority")
	viper.BindPFlag("weight", addCmd.Flags().Lookup("weight"))

	addCmd.Flags().StringP("port", "P", "", "TCP or UDP port of the host that is hosting the service")
	viper.BindPFlag("port", addCmd.Flags().Lookup("port"))

	addCmd.Flags().StringP("target", "T", "", "The canonical name of the host providing the service")
	viper.BindPFlag("target", addCmd.Flags().Lookup("target"))

	addCmd.Flags().StringP("subdomain", "s", "", "Name of the subdomain")
	viper.BindPFlag("subdomain", addCmd.Flags().Lookup("subdomain"))

	addCmd.Flags().IntP("ttl", "l", 0, "The lifetime of the DNS record in seconds")
	viper.BindPFlag("ttl", addCmd.Flags().Lookup("ttl"))

}
