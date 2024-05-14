package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// useful funcs

func loadFileInStruct(url string, responseStruct any) error {
	data, err := os.ReadFile(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func generateLink(linkparts []string, params []string) string {
	if len(params) != len(linkparts)-1 {
		panic("Error generating link")
	}

	var link string = ""
	link += linkparts[0]
	for i := 1; i < len(linkparts); i++ {
		link += params[i-1]
		link += linkparts[i]
	}

	return link
}

// handling links to data

type _LoLDataLinks struct {
	LanguagesLink      []string `json:"languages"`
	VersionLink        []string `json:"version"`
	ItemsLink          []string `json:"items"`
	ItemsFullImageLink []string `json:"itemfullimage"`
}

func _checkLinksLoading() {
	if loadDataLinksError != nil {
		fmt.Println(loadDataLinksError)
		panic("Error in links loading")
	}
}

func _checkVersionLoading() {
	if loadDataVersionError != nil {
		fmt.Println(loadDataLinksError)
		panic("Error in version loading")
	}
}

func checkLoadings() {
	_checkLinksLoading()
	_checkVersionLoading()
}

var lolDataLinks _LoLDataLinks = _LoLDataLinks{}
var loadDataLinksError error = loadFileInStruct("./loldatalinks.json", &lolDataLinks)

var lolDataVersion string
var loadDataVersionError error = loadFileInStruct(generateLink(lolDataLinks.VersionLink, []string{}), &lolDataVersion)

// handling data

type LoLLanguages []string

type LoLItems struct {
	Data map[string]LoLItem
}

type LoLItem struct {
	Name  string `json:"name"`
	Image struct {
		Full string `json:"full"`
	} `json:"image"`
}

func getLoLLanguages() LoLLanguages {
	var languages LoLLanguages = LoLLanguages{}
	err := loadFileInStruct(generateLink(lolDataLinks.LanguagesLink, []string{}), &languages)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Languages")
	}
	return languages
}

func (lolLanguages LoLLanguages) getBaseLoLLanguage() string {
	return lolLanguages[0]
}

func getLoLItems(language string) LoLItems {
	var items LoLItems = LoLItems{}
	err := loadFileInStruct(generateLink(lolDataLinks.ItemsLink, []string{lolDataVersion, language}), &items)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Items")
	}
	return items
}

// data manager

type LoLDataManager struct {
	Languages      LoLLanguages
	ChosenLanguage string
	Items          LoLItems
}

func NewLoLDataManager() LoLDataManager {
	var lolDataManager LoLDataManager = LoLDataManager{}
	lolDataManager.Languages = getLoLLanguages()
	lolDataManager.ChosenLanguage = lolDataManager.Languages.getBaseLoLLanguage()
	lolDataManager.Items = getLoLItems(lolDataManager.ChosenLanguage)
	return lolDataManager
}

// main

func main() {
	checkLoadings()

	test := NewLoLDataManager()
	full := test.Items.Data["1001"].Image.Full
	fmt.Println(generateLink(lolDataLinks.ItemsFullImageLink, []string{lolDataVersion, full}))
}
