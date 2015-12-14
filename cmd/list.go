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

	"log"

	"github.com/Lexty/yandexdns/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	propAll       string = "*"
	propId        string = "id"
	propSubdomain string = "subdomain"
	propType      string = "type"
	propContent   string = "content"
	propPriority  string = "priority"
	propTTL       string = "ttl"
	propFQDN      string = "fqdn"
	propAdminMail string = "admin_mail"
	propRetry     string = "retry"
	propRefresh   string = "refresh"
	propExpire    string = "expire"
	propMinTTL    string = "minttl"

	propsDefault string = propId + "," + propSubdomain + "," + propType + "," + propContent + "," + propPriority

	typeAll   string = "*"
	typeA     string = "A"
	typeAAAA  string = "AAAA"
	typeCNAME string = "CNAME"
	typeMX    string = "MX"
	typeNS    string = "NS"
	typeSOA   string = "SOA"
	typeTXT   string = "TXT"
	typeSRV   string = "SRV"

	formatList  string = "list"
	formatTable string = "table"
	formatJson  string = "json"
)

var props map[string]string

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
			throwError(err)
		}

		var types []string
		if viper.GetString("types") == propAll {
			types = []string{typeA, typeAAAA, typeCNAME, typeSRV, typeTXT, typeSOA, typeMX, typeNS}
		} else {
			types = parseCommaSep(viper.GetString("types"))
		}

		props := viper.GetString("props")
		if props == propAll {
			props = strings.Join([]string{propId, propType, propContent, propSubdomain, propPriority, propTTL, propFQDN, propAdminMail, propRetry, propRefresh, propExpire, propMinTTL}, ",")
		}
		setProps()

		printResponse(list, viper.GetString("format"), props, types)
	},
}

func throwError(e error) {
	log.Printf("Error: %s", e.Error())
	os.Exit(-1)
}

func parseCommaSep(raw string) []string {
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
	props[propTTL] = "TTL"
	props[propPriority] = "Priority"
	props[propFQDN] = "FQDN"
	props[propAdminMail] = "Admin Mail"
	props[propRetry] = "Retry"
	props[propRefresh] = "Refresh"
	props[propExpire] = "Expire"
	props[propMinTTL] = "MinTTL"
}

func printTable(records []*api.Record, props string) {
	parsedProps := parseCommaSep(props)
	header := make([]string, len(parsedProps))
	var err error
	for i, prop := range parsedProps {
		if header[i], err = getHeader(prop); err != nil {
			throwError(err)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, rec := range records {
		data := make([]string, len(parsedProps))
		for q, prop := range parsedProps {
			if data[q], err = getValue(prop, rec); err != nil {
				throwError(err)
			}
		}
		table.Append(data)
	}

	table.Render()
}

func printList(records []*api.Record, props string) {
	parsedProps := parseCommaSep(props)
	header := make([]string, len(parsedProps))
	var maxLen int
	var err error
	for i, prop := range parsedProps {
		if header[i], err = getHeader(prop); err != nil {
			throwError(err)
		}
		if len(header[i]) > maxLen {
			maxLen = len(header[i])
		}
	}

	for _, rec := range records {
		var value string
		for i, prop := range parsedProps {
			value, err = getValue(prop, rec)
			fmt.Printf("  %s  %s  %s\n", header[i], strings.Repeat(" ", maxLen-len(header[i])), value)
		}
		fmt.Println("")
	}
}

func printResponse(response api.Response, format, props string, types []string) {
	switch format {
	case formatJson:
		fmt.Print(response.Json)
	case formatList:
		printList(filterRecords(response.Records, types), props)
	case formatTable:
		printTable(filterRecords(response.Records, types), props)
	default:
		throwError(errors.New(fmt.Sprintf(`Unknown output format "%s".`, format)))
	}
}

func filterRecords(recs []api.Record, types []string) []*api.Record {
	var filteredRecs []*api.Record
	for i, rec := range recs {
		if isAllowedType(&rec, types) {
			filteredRecs = append(filteredRecs, &recs[i])
		}
	}
	return filteredRecs
}

func isAllowedType(rec *api.Record, types []string) bool {
	for _, t := range types {
		if typeAll == t || strings.ToUpper(rec.RecordType) == strings.ToUpper(t) {
			return true
		}
	}
	return false
}

func getHeader(prop string) (string, error) {
	if title, ok := props[strings.ToLower(prop)]; ok {
		return title, nil
	} else {
		return "", errors.New(fmt.Sprintf(`Unknown record property "%s".`, prop))
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
	case propTTL:
		val = strconv.Itoa(r.TTL)
	case propPriority:
		val = fmt.Sprintf("%v", r.Priority)
	case propFQDN:
		val = r.FQDN
	case propAdminMail:
		val = r.AdminMail
	case propRetry:
		val = strconv.Itoa(r.Retry)
	case propRefresh:
		val = strconv.Itoa(r.Refresh)
	case propExpire:
		val = strconv.Itoa(r.Expire)
	case propMinTTL:
		val = strconv.Itoa(r.MinTTL)
	default:
		return "", errors.New(fmt.Sprintf(`Unknown record property %s.`, prop))
	}
	return val, nil
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("format", "f", "", fmt.Sprintf("format output (%s|%s|%s)", formatList, formatTable, formatJson))
	viper.BindPFlag("format", listCmd.Flags().Lookup("format"))
	viper.SetDefault("format", formatList)

	listCmd.Flags().StringP("props", "p", "", fmt.Sprintf("comma separated record properties for display (available: %s) (does not work for json format)", strings.Join([]string{propAll, propId, propType, propContent, propSubdomain, propPriority, propTTL, propFQDN, propAdminMail, propRetry, propRefresh, propExpire, propMinTTL}, ", ")))
	viper.BindPFlag("props", listCmd.Flags().Lookup("props"))
	viper.SetDefault("props", propsDefault)

	listCmd.Flags().StringP("types", "t", "", fmt.Sprintf("comma separated record types for display (available: %s) (does not work for json format)", strings.Join([]string{typeAll, typeA, typeAAAA, typeCNAME, typeSRV, typeTXT, typeSOA, typeMX, typeNS}, ", ")))
	viper.BindPFlag("types", listCmd.Flags().Lookup("types"))
	viper.SetDefault("types", "*")
}
