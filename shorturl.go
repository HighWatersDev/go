package main


import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Short struct {
	ShortUrl string
}
type UrlInfo struct {
	Url string `json:"Url"`
	Title string `json:"Title"`
	ShortUrl string `json:"ShortUrl"`
//	Clicks int
//	IsArchived bool
//	PartitionKey string
//	RowKey string
//	Timestamp string
//	ETag string
}
type UList struct {
	UrlList []UrlInfo `json:"UrlList"`
}

func listUrls() {
	listUrl := "https://<Azure API url>"
	method := "GET"
	payload := strings.NewReader("")
	client := &http.Client {
	}
	req, err := http.NewRequest(method, listUrl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "<Host url>")
	req.Header.Add("Ocp-Apim-Subscription-Key", "<Secret Key>")
	req.Header.Add("Ocp-Apim-Trace", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.Status)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	text := UList{}
	jsonErr := json.Unmarshal(body, &text)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for _, value := range text.UrlList {
		//fmt.Println(value)
		out, err := json.Marshal(value)
		if err != nil {
			panic (err)
		}
		var Sr UrlInfo
		json.Unmarshal(out, &Sr)
		short := strings.ReplaceAll(Sr.ShortUrl, "<original url>", "<target url>")
		fmt.Println("Url: ", Sr.Url)
		fmt.Println("Title: ", Sr.Title)
		fmt.Println("Short Url: ", short, "\n")
	}
}
func findUrl(find string) string {
	listUrl := "https://<Azure API url>"
	method := "GET"
	payload := strings.NewReader("")
	client := &http.Client {
	}
	req, err := http.NewRequest(method, listUrl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "<Host url>")
	req.Header.Add("Ocp-Apim-Subscription-Key", "<Secret Key>")
	req.Header.Add("Ocp-Apim-Trace", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.Status)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	text := UList{}
	jsonErr := json.Unmarshal(body, &text)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var Sr UrlInfo
	var result string
	for _, value := range text.UrlList {
		//fmt.Println(value)
		out, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		//var Sr UrlInfo
		json.Unmarshal(out, &Sr)
		short := strings.ReplaceAll(Sr.ShortUrl, "<original url>", "<target url>")
		if strings.Contains(short, find) {
			result := short
			fmt.Println("Long URL:", Sr.Url)
			fmt.Println("Short URL:", result)
			fmt.Println("Title:", Sr.Title)
			return result
		} else if strings.Contains(Sr.Url, find) {
			result := Sr.Url
			fmt.Println("Long URL:", Sr.Url)
			fmt.Println("Short URL:", short)
			fmt.Println("Title:", Sr.Title)
			return result
		}
	}

	fmt.Println("No result found")
	return result
}
func createUrl(url, vanity, title string) string {

	createUrl := "https://<Azure API url>"
	method := "POST"
	payload, err := json.Marshal(map[string]string{"url":url, "title":title, "vanity":vanity})
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client {
	}
	req, err := http.NewRequest(method, createUrl, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "<Host url>")
	req.Header.Add("Ocp-Apim-Subscription-Key", "<Secret Key>")
	req.Header.Add("Ocp-Apim-Trace", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.Status)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var shortUrl Short
	json.Unmarshal(body, &shortUrl)
	vanityPart := strings.Split(shortUrl.ShortUrl, "/")
	newShortUrl := "https://<new url>/" + vanityPart[3]
	return newShortUrl
}
func main() {

	//var url string
	list := flag.Bool("list", false, "Provide -list to list existing shortened URLs")
	url := flag.String("lu", "", "long, original URL including https:// or http://")
	find := flag.String("find", "", "Look up short URL")
	var vanity string
	flag.StringVar(&vanity, "vanity", "", "string at the end of URL")
	var title string
	flag.StringVar(&title, "title", "default description", "default description")
	flag.Parse()

	if *url == "" && *list == false && *find == "" {
		fmt.Println("Provide correct input values")
		flag.PrintDefaults()
		os.Exit(1)
	} else if (*list && *url != "") || (*list && *find != "") || (*url != "" && *find != ""){
		fmt.Println("Choose either -list or -lu or -find")
		flag.PrintDefaults()
		os.Exit(1)
	} else if !strings.Contains(*url, "http:") && !strings.Contains(*url, "https://") && *list == false && *find == ""{
		fmt.Println("Link should contain http:// or https://")
		flag.PrintDefaults()
		os.Exit(1)
	}  else if *url == "" && *list && *find == ""{
		fmt.Println("Only listing URLs")
		listUrls()
	} else if *url != "" && *list == false && *find == ""{
		shortenedUrl := createUrl(*url, vanity, title)
		fmt.Println("Your short URL: ", shortenedUrl)
	} else if *find != "" && *url == "" && *list == false {
		findUrl(*find)
	}
}
