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
	ChampionsLink      []string `json:"champions"`
	SummonersLink      []string `json:"summoners"`
	RunesLink          []string `json:"runes"`
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

func checkLoadings() any {
	_checkLinksLoading()
	_checkVersionLoading()
	return nil
}

var lolDataLinks _LoLDataLinks = _LoLDataLinks{}
var loadDataLinksError error = loadFileInStruct("./loldatalinks.json", &lolDataLinks)

var lolDataVersion string
var loadDataVersionError error = loadFileInStruct(generateLink(lolDataLinks.VersionLink, []string{}), &lolDataVersion)

// kinda cursed ^^
var _ any = checkLoadings()

// handling data structs

type LoLLanguages []string

type LoLItems struct {
	Data map[string]LoLItem `json:"data"`
}

/*
   "gold": {
       "base": 300,
       "purchasable": true,
       "total": 300,
       "sell": 210
   },
   "tags": [
       "Boots"
   ],
   "maps": {
       "11": true,
       "12": true,
       "21": true,
       "22": false,
       "30": false
   },
*/

type LoLItem struct {
	ItemKey string
	Name    string   `json:"name"`
	From    []string `json:"from"`
	Into    []string `json:"into"`
	Gold    struct {
		Base        int  `json:"base"`
		Purchasable bool `json:"purchasable"`
		Total       int  `json:"total"`
	} `json:"gold"`
	Tags     []string        `json:"tags"`
	Maps     map[string]bool `json:"maps"`
	Metadata []string
}

type LoLChampions struct {
	Data map[string]LoLChampion `json:"data"`
}

type LoLChampion struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	Metadata []string
}

type LoLSummoners struct {
	Data map[string]LoLSummoner `json:"data"`
}

type LoLSummoner struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LoLRunes []LoLRuneTree

type LoLRuneTree struct {
	Name  string       `json:"name"`
	Slots []LoLRuneRow `json:"slots"`
}

type LoLRuneRow struct {
	Runes []LoLRune `json:"runes"`
}

type LoLRune struct {
	Name string `json:"name"`
}

// LoLLanguages methods

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

// LoLItems methods

func getLoLItems(language string) LoLItems {
	var items LoLItems = LoLItems{}
	err := loadFileInStruct(generateLink(lolDataLinks.ItemsLink, []string{lolDataVersion, language}), &items)

	for key, item := range items.Data {
		item.ItemKey = key

		if !item.hasChildren() {
			item.Metadata = append(item.Metadata, "final")
		}

		if !item.hasParents() {
			item.Metadata = append(item.Metadata, "base")
		}

		items.Data[key] = item
	}

	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Items")
	}
	return items
}

func (lolItems LoLItems) getLoLItem(num string) LoLItem {
	var lolItem LoLItem = lolItems.Data[num]
	return lolItem
}

func (lolItem LoLItem) getLoLItemParents(lolItems LoLItems) []LoLItem {
	var parents []LoLItem = []LoLItem{}
	for _, e := range lolItem.From {
		parents = append(parents, lolItems.getLoLItem(e))
	}
	return parents
}

func (lolItem LoLItem) getLoLItemChildren(lolItems LoLItems) []LoLItem {
	var children []LoLItem = []LoLItem{}
	for _, e := range lolItem.Into {
		children = append(children, lolItems.getLoLItem(e))
	}
	return children
}

func (lolItem LoLItem) hasParents() bool {
	return len(lolItem.From) > 0
}

func (lolItem LoLItem) hasChildren() bool {
	return len(lolItem.Into) > 0
}

// LoLLanguages methods

func getLoLChampions(language string) LoLChampions {
	var champs LoLChampions = LoLChampions{}
	err := loadFileInStruct(generateLink(lolDataLinks.ChampionsLink, []string{lolDataVersion, language}), &champs)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Champions")
	}
	return champs
}

// LoLLanguages methods

func getLoLSummoners(language string) LoLSummoners {
	var summs LoLSummoners = LoLSummoners{}
	err := loadFileInStruct(generateLink(lolDataLinks.SummonersLink, []string{lolDataVersion, language}), &summs)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Summoners")
	}
	return summs
}

// LoLLanguages methods

func getLoLRunes(language string) LoLRunes {
	var runes LoLRunes = LoLRunes{}
	err := loadFileInStruct(generateLink(lolDataLinks.RunesLink, []string{lolDataVersion, language}), &runes)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Runes")
	}
	return runes
}

// data manager

type LoLDataManager struct {
	Languages      LoLLanguages
	ChosenLanguage string
	Items          LoLItems
	Champions      LoLChampions
	Summoners      LoLSummoners
	Runes          LoLRunes
}

func NewLoLDataManager() LoLDataManager {
	var lolDataManager LoLDataManager = LoLDataManager{}
	lolDataManager.Languages = getLoLLanguages()
	lolDataManager.ChosenLanguage = lolDataManager.Languages.getBaseLoLLanguage()
	lolDataManager.updateData()
	return lolDataManager
}

func (lolDataManager *LoLDataManager) updateData() {
	lolDataManager.Items = getLoLItems(lolDataManager.ChosenLanguage)
	lolDataManager.Champions = getLoLChampions(lolDataManager.ChosenLanguage)
	lolDataManager.Summoners = getLoLSummoners(lolDataManager.ChosenLanguage)
	lolDataManager.Runes = getLoLRunes(lolDataManager.ChosenLanguage)
}

// main

func Test() {
	test := NewLoLDataManager()
	testitems := test.Items
	item := testitems.getLoLItem("3067")
	fmt.Println(item)
	fmt.Println(item.getLoLItemParents(testitems))
	fmt.Println(item.getLoLItemChildren(testitems))
}
