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

	"github.com/Lexty/yandex-dns-cli-manager/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the DNS record by ID",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := api.DeleteRecordById(viper.GetInt("id"), viper.GetString("domain"), viper.GetString("admin-token"))

		if err != nil {
			throwError(err)
		}

		switch viper.GetString("format") {
		case formatJson:
			fmt.Print(resp.Json)
		case formatList:
			fmt.Printf("Record successfully deleted\n\n")
		default:
			throwError(errors.New(fmt.Sprintf(`Unknown output format "%s".`, viper.GetString("format"))))
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("format", "f", "", fmt.Sprintf("format output (%s|%s)", formatList, formatJson))
	viper.BindPFlag("format", deleteCmd.Flags().Lookup("format"))
	viper.SetDefault("format", formatList)

	deleteCmd.Flags().IntP("id", "i", 0, "ID of the record")
	viper.BindPFlag("id", deleteCmd.Flags().Lookup("id"))
}
