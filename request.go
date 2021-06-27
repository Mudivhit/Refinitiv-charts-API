package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const tokenURL = "https://api.rkd.refinitiv.com/api/TokenManagement/TokenManagement.svc/REST/Anonymous/TokenManagement_1/CreateServiceToken_1"
const chartURL = "http://api.rkd.refinitiv.com/api/Charts/Charts.svc/REST/Charts_1/GetChart_2"

type data struct{
	Server string `json:"Server"`
	Tag string `json:"Tag"`
	Url string `json:"Url"`
	SecureUrl string `json:"SecureUrl"`
	MPLSURL string `json:"MPLSURL"`
	SecureMPLSURL string `json:"SecureMPLSURL"`
}

type response struct {
	GetChart_Response_2 chart `json:"GetChart_Response_2"`
}

type chart struct {
	ChartImageResult data `json:"ChartImageResult"`

}

func main(){
	xmlFile, err := os.Open("sample_request.json")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var token string = string(GetToken())
	resp := ExecuteRequest(chartURL,byteValue , token )
	//wanted := string(resp)
	//fmt.Println(wanted)
	var info response
	err = json.Unmarshal(resp,&info)
	if err != nil{
		log.Println(err)
	}

	link := info.GetChart_Response_2.ChartImageResult.Url
	err = DownloadFile("Chart_001.png", link)
	if err != nil{
		panic(err)
	}
	fmt.Println("Downloaded chart : ", link)
}

//Creating service token to access API
func GetToken() []byte {
	login, err := ioutil.ReadFile("./logins.json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(login))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var jsonData map[string]interface{}
	json.Unmarshal(body, &jsonData)

	serv, _ := json.Marshal(jsonData["CreateServiceToken_Response_1"])
	json.Unmarshal(serv, &jsonData)

	token, _ := json.Marshal(jsonData["Token"])
	json.Unmarshal(token, &jsonData)

	return token[1 : len(token)-1]
}

//Sending request to charts API
func ExecuteRequest(link string, body []byte, token string) []byte {
	jsonReader := io.Reader(bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", link , jsonReader)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Trkd-Auth-Token",token)
	req.Header.Add("X-Trkd-Auth-ApplicationID", "AdriAndilesolutionsCom")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	rates := make([]byte, resp.ContentLength)
	resp.Body.Read(rates)
	return rates
}

//Downloading to save locally
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


