package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	collection := client.Database("vietnamworks").Collection("jobs")
	url := "https://ms.vietnamworks.com/job-search/v1.0/search"

	payload := []byte(`{"query":"","filter":[{"field":"workingLocations.cityId","value":"29"},{"field":"workingLocations.districtId","value":"[{\"cityId\":29,\"districtId\":[-1]}]"}],"ranges":[],"order":[],"hitsPerPage":1000,"page":0,"retrieveFields":["benefits","jobTitle","salaryMax","isSalaryVisible","jobLevelVI","isShowLogo","salaryMin","companyLogo","userId","jobLevel","jobId","companyId","approvedOn","isAnonymous","alias","expiredOn","industries","workingLocations","services","companyName","salary","onlineOn","simpleServices","visibilityDisplay","isShowLogoInSearch","priorityOrder","skills","profilePublishedSiteMask"]}`)
	page := 0
	end := false
	for end != true {
		payloadWithPage := bytes.Replace(payload, []byte("\"page\":0"), []byte(fmt.Sprintf("\"page\":%d", page)), -1)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadWithPage))

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		insertManyResult, err := collection.InsertMany(context.Background(), result["data"].([]interface{}))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted multiple documents: ", len(insertManyResult.InsertedIDs))
		if err != nil {
			log.Fatal(err)
		}
		nPages := int(result["meta"].(map[string]interface{})["nbPages"].(float64))
		page++
		if page == nPages {
			end = true
		}
		resp.Body.Close()
	}
}
