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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
)

type shortenURLAPIResponse struct {
	ShortCode      string `json:"shortCode"`
	ShortUrl       string `json:"shortUrl"`
	LongUrl        string `json:"longUrl"`
	DeviceLongUrls struct {
		Android interface{} `json:"android"`
		Ios     interface{} `json:"ios"`
		Desktop interface{} `json:"desktop"`
	} `json:"deviceLongUrls"`
	DateCreated   time.Time `json:"dateCreated"`
	VisitsSummary struct {
		Total   int `json:"total"`
		NonBots int `json:"nonBots"`
		Bots    int `json:"bots"`
	} `json:"visitsSummary"`
	Tags      []string    `json:"tags"`
	Meta      meta        `json:"meta"`
	Domain    interface{} `json:"domain"`
	Title     interface{} `json:"title"`
	Crawlable bool        `json:"crawlable"`
}

type meta struct {
	ValidSince time.Time   `json:"validSince"`
	ValidUntil interface{} `json:"validUntil"`
	MaxVisits  int         `json:"maxVisits"`
}

type shortenURLAPIRequest struct {
	LongUrl      string `json:"longUrl"`
	CustomSlug   string `json:"customSlug,omitempty"`
	FindIfExists bool   `json:"findIfExists"`
}

type APIResponse struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

var customSlug string

// shortCmd represents the short command
var shortCmd = &cobra.Command{
	Use:     "short",
	Version: ApplicationVersion,
	Short:   "Shortens a long URL using the Shlink API",
	Long:    `Shortens a given URL using the Shlink service through its API and copies it to the clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		URLsToShorten := []string{}
		URLsToShorten = append(URLsToShorten, args...)
		enableCopyToClipboard := false

		if len(URLsToShorten) == 0 {
			fmt.Println("No URLs to shorten were provided")
			os.Exit(0)
		}

		if len(URLsToShorten) == 1 {
			enableCopyToClipboard = true
		}

		//Shorten all URLs
		for _, URLToShorten := range URLsToShorten {
			fmt.Printf("Shortening %s\n", URLToShorten)
			shortenedURL, err := shortenURL(viper.GetString("shlink_url"), viper.GetString("api_key"), viper.GetInt("timeout"), URLToShorten)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if enableCopyToClipboard {
				clipboard.Write(clipboard.FmtText, []byte(shortenedURL))
			}

			color.Blue(shortenedURL)

		}
	},
}

func init() {
	rootCmd.AddCommand(shortCmd)

	shortCmd.PersistentFlags().StringVarP(&customSlug, "slug", "g", "", "Set a specific slug for the shortened URL")
}

func shortenURL(host string, apiKey string, timeout int, URLToShorten string) (string, error) {
	var shortResponse shortenURLAPIResponse
	buildBody := shortenURLAPIRequest{
		LongUrl:      URLToShorten,
		FindIfExists: true,
	}

	if customSlug != "" {
		buildBody.CustomSlug = customSlug
	}
	payloadBytes, err := json.Marshal(buildBody)
	if err != nil {
		return "", err
	}

	suResponse, _, err1 := RestRequest("POST", host+"/rest/v2/short-urls", apiKey, time.Duration(timeout)*time.Second, bytes.NewBuffer(payloadBytes))
	if err1 != nil {
		return "", err1
	}

	err = json.Unmarshal(suResponse, &shortResponse)
	if err != nil {
		return "", err
	}

	return shortResponse.ShortUrl, nil
}
