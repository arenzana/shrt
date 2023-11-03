// Copyright © 2023 Ismael Arenzana <isma@arenzana.org>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ShortUrls struct {
	SU struct {
		Data []struct {
			ShortCode      string `json:"shortCode"`
			ShortUrl       string `json:"shortUrl"`
			LongUrl        string `json:"longUrl"`
			DeviceLongUrls struct {
				Android interface{} `json:"android"`
				Ios     interface{} `json:"ios"`
				Desktop interface{} `json:"desktop"`
			} `json:"deviceLongUrls"`
			DateCreated   string `json:"dateCreated"`
			VisitsSummary struct {
				Total   int `json:"total"`
				NonBots int `json:"nonBots"`
				Bots    int `json:"bots"`
			} `json:"visitsSummary"`
			Tags []string `json:"tags"`
			Meta struct {
				ValidSince string      `json:"validSince"`
				ValidUntil interface{} `json:"validUntil"`
				MaxVisits  interface{} `json:"maxVisits"`
			} `json:"meta"`
			Domain    interface{} `json:"domain"`
			Title     interface{} `json:"title"`
			Crawlable bool        `json:"crawlable"`
		} `json:"data"`
	} `json:"shortUrls"`
	Pagination struct {
		CurrentPage        int `json:"currentPage"`
		PagesCount         int `json:"pagesCount"`
		ItemsPerPage       int `json:"itemsPerPage"`
		ItemsInCurrentPage int `json:"itemsInCurrentPage"`
		TotalItems         int `json:"totalItems"`
	} `json:"pagination"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Version: ApplicationVersion,
	Short:   "Lists all the URLs shortened in the configured Shlink service",
	Long:    `Lists all the URLs stored in the configured Shlink service for reuse`,
	Run: func(cmd *cobra.Command, args []string) {

		if !quiet {
			fmt.Printf("Listing URLs from %s\n", viper.GetString("shlink_url"))
		}

		URLResponse, err := getShortURLList(viper.GetString("shlink_url"), viper.GetString("api_key"), viper.GetInt("timeout"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		outputURLs := [][]string{}

		for _, shortURL := range URLResponse.SU.Data {
			t, err := time.Parse(time.RFC3339, shortURL.DateCreated)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			longUrl := strings.Split(shortURL.LongUrl, "?")
			record := []string{t.Format("2006-01-02"), shortURL.ShortUrl, longUrl[0], fmt.Sprintf("%d", shortURL.VisitsSummary.Total)}
			outputURLs = append(outputURLs, record)
		}

		//Display table options and render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Created", "Short URL", "Long URL", "Visits"})

		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoFormatHeaders(true)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		table.SetBorder(false)

		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.Bold})

		table.SetColumnColor(tablewriter.Colors{},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
			tablewriter.Colors{})

		table.SetRowSeparator("-")
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)
		table.AppendBulk(outputURLs) // Add Bulk Data
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// getShortURLList returns a list of short URLs from the Shlink API
func getShortURLList(host string, apiKey string, timeout int) (ShortUrls, error) {
	// Get the list of short URLs
	// GET /rest/v2/short-urls
	// https://app.shlink.io/#/docs#operation/getShortUrls

	shortURLResponse, _, err := RestRequest("GET", host+"/rest/v2/short-urls", apiKey, time.Duration(timeout)*time.Second, nil)
	if err != nil {
		return ShortUrls{}, err
	}

	var shortUrls ShortUrls

	err1 := json.Unmarshal(shortURLResponse, &shortUrls)
	if err1 != nil {
		return ShortUrls{}, err1
	}

	return shortUrls, nil
}

// RestRequest returns data from a REST endpoint
func RestRequest(method string, uri string, apiKey string, timeout time.Duration, body io.Reader) ([]byte, time.Duration, error) {
	start := time.Now()

	restClient := http.Client{
		Timeout: timeout,
	}

	// Create an HTTP request with custom headers
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := restClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode == 401 {
		return nil, 0, fmt.Errorf("Unauthorized")
	}

	if resp.StatusCode == 500 {
		return nil, 0, fmt.Errorf("Internal Server Error")
	}

	if resp.StatusCode == 400 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}

		return nil, 0, fmt.Errorf("%+v", string(data))
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, time.Since(start), nil
}
