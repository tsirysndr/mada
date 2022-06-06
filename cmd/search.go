/*
Copyright © 2022 Tsiry Sandratraina
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tsirysndr/mada/mada"
	"github.com/tsirysndr/mada/types"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a region, district, commune, or fokontany",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		outputInJSON, _ := cmd.Flags().GetBool("json")
		searchForFokontany, _ := cmd.Flags().GetBool("fokontany")
		searchForCommune, _ := cmd.Flags().GetBool("commune")
		searchForDistrict, _ := cmd.Flags().GetBool("district")
		searchForRegion, _ := cmd.Flags().GetBool("region")
		openInBrowser, _ := cmd.Flags().GetBool("open")

		options := types.SearchOptions{
			OutputInJSON:       outputInJSON,
			SearchForFokontany: searchForFokontany,
			SearchForCommune:   searchForCommune,
			SearchForDistrict:  searchForDistrict,
			SearchForRegion:    searchForRegion,
			OpenInBrowser:      openInBrowser,
		}

		db, err := mada.OpenDatabaseConnection()
		if err != nil {
			panic(err)
		}

		index, err := mada.InitializeBleve(db)
		if err != nil {
			panic(err)
		}

		s := mada.NewSearchService(db, index)

		result, err := s.Search(args[0], options)
		if err != nil {
			panic(err)
		}
		mada.FormatSearchResultOrOpenInBrowser(db, index, result, options)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	searchCmd.Flags().BoolP("fokontany", "f", false, "Search for a fokontany")
	searchCmd.Flags().BoolP("commune", "c", false, "Search for a commune")
	searchCmd.Flags().BoolP("district", "d", false, "Search for a district")
	searchCmd.Flags().BoolP("region", "r", false, "Search for a region")
	searchCmd.Flags().BoolP("open", "o", false, "Open the result in the browser")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
