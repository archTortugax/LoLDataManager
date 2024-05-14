package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoLVersions []string

type LoLLanguages []string

func getLoLVersions() LoLVersions {
	var versionslink string = "https://ddragon.leagueoflegends.com/api/versions.json"
	var versions LoLVersions = LoLVersions{}
	err := loadResponseInStruct(versionslink, &versions) 
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Versions")
	}
	return versions
}

func (lolVersions LoLVersions) getLatestLoLVersion() string {
	return lolVersions[0]
}

func getLoLLanguages() LoLLanguages {	
	var languageslink string = "https://ddragon.leagueoflegends.com/cdn/languages.json"
	var languages LoLLanguages = LoLLanguages{}
	err := loadResponseInStruct(languageslink, &languages) 
	if err != nil {
		fmt.Println(err.Error())
		panic("Error in loading Languages")
	}
	return languages
}

func (lolLanguages LoLLanguages) getBaseLoLLanguage() string {
	return lolLanguages[0]
}

type LoLDataManager struct {
	Versions LoLVersions
	ChosenVersion string
	
	Languages LoLLanguages
	ChosenLanguage string
}

func NewLoLDataManager() LoLDataManager {
	var lolDataManager LoLDataManager = LoLDataManager{}
	lolDataManager.Versions = getLoLVersions()
	lolDataManager.ChosenVersion = lolDataManager.Versions.getLatestLoLVersion()
	lolDataManager.Languages = getLoLLanguages()
	lolDataManager.ChosenLanguage = lolDataManager.Languages.getBaseLoLLanguage()
	return lolDataManager
}

func loadResponseInStruct(url string, responseStruct any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var lolDataManager LoLDataManager = NewLoLDataManager()
	fmt.Println(lolDataManager.Versions)
	fmt.Println(lolDataManager.ChosenVersion)
	fmt.Println(lolDataManager.Languages)
	fmt.Println(lolDataManager.ChosenLanguage)
}
