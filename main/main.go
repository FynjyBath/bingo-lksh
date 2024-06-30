package main

import (
	"config"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Contest struct {
	Problems [][]string       `json:"problems"`
	Teams    map[string][]int `json:"solved"`
	Link     string           `json:"link"`
}

func QueryToAPI(cfg config.Config) (Contest, error) {
	now := strconv.Itoa(int(time.Now().Unix()))

	randomNumber := rand.Intn(900000) + 100000
	text := fmt.Sprintf("%d/contest.standings?apiKey=%s&contestId=%d&time=%s#%s",
		randomNumber, cfg.ApiKey, cfg.ContestID, now, cfg.ApiSecret)
	hash := sha512.Sum512([]byte(text))
	hashString := strconv.Itoa(randomNumber) + hex.EncodeToString(hash[:])

	apiURL := "https://codeforces.com/api/contest.standings"

	params := url.Values{}
	params.Add("apiKey", cfg.ApiKey)
	params.Add("contestId", strconv.Itoa(cfg.ContestID))
	params.Add("time", now)
	params.Add("apiSig", hashString)

	apiURL = apiURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return Contest{}, errors.New("Ошибка при создании запроса: " + fmt.Sprint(err))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Contest{}, errors.New("Ошибка при отправке запроса: " + fmt.Sprint(err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Contest{}, errors.New("Ошибка при чтении ответа: " + fmt.Sprint(err))
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return Contest{}, err
	}
	result := jsonData["result"].(map[string]interface{})

	var contest Contest
	contest.Link = cfg.Link

	var lst = make([]string, 0)
	contest.Teams = make(map[string][]int)

	problems := result["problems"].([]interface{})
	for _, problem := range problems {
		lst = append(lst, (problem.(map[string]interface{}))["index"].(string))
	}

	sz := int(math.Sqrt(float64(len(lst) + 1)))
	contest.Problems = make([][]string, sz)
	for i := 0; i < len(lst); i += sz {
		for j := i; j < i+sz; j++ {
			contest.Problems[i/sz] = append(contest.Problems[i/sz], lst[j])
		}
	}

	rows := result["rows"].([]interface{})
	for _, rowInterface := range rows {
		row := rowInterface.(map[string]interface{})

		listMembers := (row["party"].(map[string]interface{}))["members"].([]interface{})
		team := ""
		for _, member := range listMembers {
			if team != "" {
				team += ", "
			}
			team += (member.(map[string]interface{}))["handle"].(string)
		}

		listProblems := row["problemResults"].([]interface{})
		for _, problem := range listProblems {
			contest.Teams[team] = append(contest.Teams[team], int((problem.(map[string]interface{})["points"]).(float64)))
		}
	}

	return contest, nil
}

func GetContest() Contest {
	cfg := config.LoadConfig("config.json")
	contest, err := QueryToAPI(cfg)
	if err != nil {
		log.Println("!ERROR!", err)
		return Contest{}
	}
	log.Println(contest)
	return contest
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	contest := GetContest()

	jsonData, err := json.Marshal(contest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(jsonData))
}

func main() {
	http.HandleFunc("/get_table", GetTable)

	log.Println("Listening on :3000...")
	err := http.ListenAndServeTLS(":3000", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
