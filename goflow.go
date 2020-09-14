package main

import (
	"fmt"
	"os"
	"sort"
	"os/exec"
	"runtime"
	"net/http"
	"net/url"
	"log"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
)

var api = os.Getenv("FLOW_API_URL")

type authData struct {
	AccessToken string `json:"access_token"`
	Scope string `json:"scope"`
	Expires string `json:"expires_in"`
	TokenType string `json:"token_type"`
	Token string `json:"id_token"`
}

type results struct {
	Folders []folder_list
	Surveys []survey_list
}

type folder_list struct {
	Id string `json:"id"`
	Name string `json:"name"`
	ParentId string `json:"parentId"`
	FoldersUrl string `json:"foldersUrl"`
	SurveysUrl string `json:"surveysUrl"`
	CreatedAt string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

type survey_list struct {
	Id string `json:id`
	Name string `json:"name"`
	FolderId string `json:"folderId"`
	SurveyUrl string `json:"surveyUrl"`
	CreatedAt string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

type survey struct {
	Id string `json:id`
	Name string `json:name`
	Registration string `json:registrationFormId`
	Forms []form_list
	CreatedAt string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	DataPointsUrl string `json:"dataPointsUrl"`
}

type form_list struct {
	Id string `json:id`
	Name string `json:"name"`
	CreatedAt string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

func getData(url string, token string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.akvo.flow.v2+json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func ask(instance string, foldersUrl string, surveysUrl string, token string) {
	var folders_data results
	var surveys_data results
	var dataslice []string
	json.Unmarshal(getData(foldersUrl, token), &folders_data)
	json.Unmarshal(getData(surveysUrl, token), &surveys_data)
	fdata := folders_data.Folders
	sort.SliceStable(fdata, func(i, j int) bool {
		return fdata[i].Id < fdata[j].Id
	})
	for _, f := range fdata {
		dataslice = append(dataslice, "FOLDER|" + f.Id + "|" + f.Name)
	}
	sdata := surveys_data.Surveys
	sort.SliceStable(sdata, func(i, j int) bool {
		return sdata[i].Id < sdata[j].Id
	})
	for _, f := range surveys_data.Surveys {
		dataslice = append(dataslice, "SURVEY|" + f.Id + "|" + f.Name)
	}

	prompt := promptui.Select{
		Label: "Select folder or surveys",
		Size: 20,
		Items: dataslice,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Aborting ", err)
		return
	}
	selected := strings.Split(result, "|")
	if selected[0] == "SURVEY" {
		for _, s := range surveys_data.Surveys {
			if s.Id == selected[1] {
				var survey_data survey
				var formslice []string
				json.Unmarshal(getData(s.SurveyUrl, token), &survey_data)
				for _, f := range survey_data.Forms {
					formslice = append(formslice, "FORM|" + f.Id + "|" + f.Name)
				}
				endprompt := promptui.Select{
					Label: "Select Form to Open",
					Items: formslice,
				}
				_, forminfo , err := endprompt.Run()
				if err != nil {
					fmt.Println("Aborting ", err)
					return
				}
				form := strings.Split(forminfo, "|")
				formUrl := fmt.Sprintf("https://tech-consultancy.akvo.org/akvo-flow-web/%s/%s", instance, form[1])
				fmt.Printf("\nOpening: %s\n", formUrl)
				apiUrl := strings.Replace(formUrl, "akvo-flow-web", "akvo-flow-web-api", -1)
				fmt.Printf("API: %s/update\n", apiUrl)
				openbrowser(formUrl)
			}
		}
	}
	if selected[0] == "FOLDER" {
		for _, f := range folders_data.Folders {
			if f.Id == selected[1] {
				ask(instance, f.FoldersUrl, f.SurveysUrl, token)
			}
		}
	}
	return
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getheaders() string {
	formData := url.Values{
		"client_id":{os.Getenv("AUTH0_CLIENT_FLOW")},
		"username":{os.Getenv("AUTH0_USER")},
		"password":{os.Getenv("AUTH0_PWD")},
		"grant_type":{"password"},
		"scope":{"openid email"},
	}
	resp, err := http.PostForm(os.Getenv("AUTH0_URL") + "/token", formData)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var token authData
	json.Unmarshal(body, &token)
	return token.Token
}

func main() {
	fmt.Print("Flow Instance: ")
	var instance string
	_, err := fmt.Scanln(&instance)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Print("Logging in to ")
	color.Info.Println("https://" + instance + ".akvoflow.org")
	token := fmt.Sprintf("Bearer %s", getheaders())
	foldersUrl := fmt.Sprintf("%s/%s/folders", api, instance)
	surveysUrl := fmt.Sprintf("%s/%s/surveys", api, instance)
	ask(instance, foldersUrl, surveysUrl, token)
	return
}
