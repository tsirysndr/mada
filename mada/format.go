package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/pkg/browser"
	"github.com/tsirysndr/mada/types"
)

func FormatResultOrOpenInBrowser(db *sql.DB, index bleve.Index, searchResults *bleve.SearchResult, openInBrowser, outputInJSON bool) {
	if openInBrowser {
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
		if err != nil {
			fmt.Printf("Open http://localhost:%d in your browser\n", PORT)
		}
		StartHttpServer(db, index)
	}

	if !outputInJSON {
		fmt.Println(searchResults)
		return
	}

	b, _ := json.MarshalIndent(searchResults.Hits, "", "  ")

	fmt.Println(string(b))
}

func FormatSearchResultOrOpenInBrowser(db *sql.DB, index bleve.Index, searchResults *types.SearchResult, opt types.SearchOptions) {
	if opt.OpenInBrowser {
		url := fmt.Sprintf("http://localhost:%d", PORT)
		if searchResults.Fokontany != nil {
			url = fmt.Sprintf("%s/#/fokontany/%s", url, searchResults.Fokontany.ID)
		}
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Printf("Open %s in your browser\n", url)
		}
		StartHttpServer(db, index)
	}

	if !opt.OutputInJSON {
		FormatSearchResults(searchResults)
		return
	}

	if searchResults.Fokontany != nil {
		b, _ := json.MarshalIndent(searchResults.Fokontany, "", "  ")
		fmt.Println(string(b))
		return
	}
	if searchResults.Commune != nil {
		b, _ := json.MarshalIndent(searchResults.Commune, "", "  ")
		fmt.Println(string(b))
		return
	}
	if searchResults.District != nil {
		b, _ := json.MarshalIndent(searchResults.District, "", "  ")
		fmt.Println(string(b))
		return
	}
	if searchResults.Region != nil {
		b, _ := json.MarshalIndent(searchResults.Region, "", "  ")
		fmt.Println(string(b))
		return
	}
	b, _ := json.MarshalIndent(searchResults.Result.Hits, "", "  ")

	fmt.Println(string(b))
}

func FormatSearchResults(searchResults *types.SearchResult) {
	if searchResults.Fokontany != nil {
		fmt.Printf(`
				point
								%v
				id
								%s
				name
								%s
				commune
								%s
				district
								%s
				region
								%s
				country
								%s
				type
								fokontany
				geometry
								%v
	`,
			searchResults.Fokontany.Point,
			searchResults.Fokontany.ID,
			searchResults.Fokontany.Name,
			searchResults.Fokontany.Commune,
			searchResults.Fokontany.District,
			searchResults.Fokontany.Region,
			searchResults.Fokontany.Country,
			searchResults.Fokontany.Coordinates,
		)
		return
	}

	if searchResults.Commune != nil {
		fmt.Printf(`
				point
								%v
				id
								%s
				name
								%s
				district
								%s
				region
								%s
				country
								%s
				type
								commune
				geometry
								%v
	`,
			searchResults.Commune.Point,
			searchResults.Commune.ID,
			searchResults.Commune.Name,
			searchResults.Commune.District,
			searchResults.Commune.Region,
			searchResults.Commune.Country,
			searchResults.Commune.Coordinates,
		)
		return
	}

	if searchResults.District != nil {
		fmt.Printf(`
				point
								%v
				id
								%s
				name
								%s
				region
								%s
				country
								%s
				type
								district
				geometry
								%v
	`,
			searchResults.District.Point,
			searchResults.District.ID,
			searchResults.District.Name,
			searchResults.District.Region,
			searchResults.District.Country,
			searchResults.District.Coordinates,
		)
		return
	}

	if searchResults.Region != nil {
		fmt.Printf(`
				point
								%v
				id
								%s
				name
								%s
				country
								%s
				type
								region
				geometry
								%v
	`,
			searchResults.Region.Point,
			searchResults.Region.ID,
			searchResults.Region.Name,
			searchResults.Region.Country,
			searchResults.Region.Coordinates,
		)
		return
	}

	fmt.Println(searchResults.Result)

}

func FormatOrOpenFokontanyInBrowser(db *sql.DB, index bleve.Index, fokontany *types.Fokontany, openInBrowser, outputInJSON bool) {
	if openInBrowser {
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d/#/fokontany/%s", PORT, fokontany.ID))
		if err != nil {
			fmt.Printf("Open http://localhost:%d/#/fokontany/%s", PORT, fokontany.ID)
		}
		StartHttpServer(db, index)
	}
	if outputInJSON {
		b, _ := json.MarshalIndent(fokontany, "", "  ")
		fmt.Println(string(b))
		return
	}
	fmt.Printf(`
					id
									%s
					name
									%s
					commune
									%s
					district
									%s
					region
									%s
          country
									%s
					type
									fokontany
					geometry
									%v
		`,
		fokontany.ID,
		fokontany.Name,
		fokontany.Commune,
		fokontany.District,
		fokontany.Region,
		fokontany.Country,
		fokontany.Coordinates)

}

func FormatOrOpenCommuneInBrowser(db *sql.DB, index bleve.Index, commune *types.Commune, openInBrowser, outputInJSON bool) {
	if openInBrowser {
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d/#/communes/%s", PORT, commune.ID))
		if err != nil {
			fmt.Printf("Open http://localhost:%d/#/communes/%s", PORT, commune.ID)
		}
		StartHttpServer(db, index)
	}
	if outputInJSON {
		b, _ := json.MarshalIndent(commune, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf(`
        id
                %s
        name
                %s
        district
                %s
        
        region
                %s
        type
                commune
        country
                Madagascar
        geometry
                %v
	`, commune.ID, commune.Name, commune.District, commune.Region, commune.Coordinates)

}

func FormatOrOpenDistrictInBrowser(db *sql.DB, index bleve.Index, district *types.District, openInBrowser, outputInJSON bool) {
	if openInBrowser {
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d/#/districts/%s", PORT, district.ID))
		if err != nil {
			fmt.Printf("Open http://localhost:%d/#/districts/%s", PORT, district.ID)
		}
		StartHttpServer(db, index)
	}
	if outputInJSON {
		b, _ := json.MarshalIndent(district, "", "  ")
		fmt.Println(string(b))
		return
	}
	fmt.Printf(`
        id
                %s
        name
                %s
        
        region
                %s
        type
                district
        country
                Madagascar
        geometry
                %v
	`, district.ID, district.Name, district.Region, district.Coordinates)
}

func FormatOrOpenRegionInBrowser(db *sql.DB, index bleve.Index, region *types.Region, openInBrowser, outputInJSON bool) {
	if openInBrowser {
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d/#/regions/%s", PORT, region.ID))
		if err != nil {
			fmt.Printf("Open http://localhost:%d/#/regions/%s", PORT, region.ID)
		}
		StartHttpServer(db, index)
	}
	if outputInJSON {
		b, _ := json.MarshalIndent(region, "", "  ")
		fmt.Println(string(b))
		return
	}
	fmt.Printf(`
        id
                %s
        name
                %s
        type
                region
        country
                Madagascar
        geometry
               %v
	`, region.ID, region.Name, region.Coordinates)
}
