package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'avgRotorSpeed' function below.
 *
 * URL for cut and paste
 * https://jsonmock.hackerrank.com/api/iot_devices/search?status={statusQuery}&page={number}
 *
 * The function is expected to return an INTEGER.
 * The function accepts following parameters:
 *  1. STRING statusQuery
 *  2. INTEGER parentId
 */

type T2 struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Data       []struct {
		Id              int    `json:"id"`
		Timestamp       int64  `json:"timestamp"`
		Status          string `json:"status"`
		OperatingParams struct {
			RotorSpeed    int     `json:"rotorSpeed"`
			Slack         float64 `json:"slack"`
			RootThreshold int     `json:"rootThreshold"`
		} `json:"operatingParams"`
		Asset struct {
			Id    int    `json:"id"`
			Alias string `json:"alias"`
		} `json:"asset"`
		Parent *struct {
			Id    int    `json:"id"`
			Alias string `json:"alias"`
		} `json:"parent,omitempty"`
	} `json:"data"`
}

func avgRotorSpeed(statusQuery string, parentId int32) int32 {
	base := "https://jsonmock.hackerrank.com/api/iot_devices/search?status=" + statusQuery + "&page=0"
	resp, err := http.Get(base)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data T2
	if err := json.Unmarshal([]byte(body), &data);	err != nil {
		panic(err)
	}
	count := 0
	totalSpeed := 0
	totalPages := data.TotalPages
	currentPage := 0
	for i := 0; i < totalPages; i++ {
		for _, devMap := range data.Data {
			if devMap.Parent != nil && devMap.Parent != nil && devMap.Parent.Id == int(parentId) {
				count++
				totalSpeed += devMap.OperatingParams.RotorSpeed
			}
		}
		currentPage++
		base := "https://jsonmock.hackerrank.com/api/iot_devices/search?status=" + statusQuery + "&page=" + strconv.Itoa(currentPage)
		resp, err := http.Get(base)
		if err != nil {
			panic(err)
		}
		//defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
	}
	return int32(totalSpeed/count)
}

func main() {
	//reader := bufio.NewReaderSize(os.Stdin, 16 * 1024 * 1024)
	pID := os.Getenv("PARENTID")
	//checkError(err)
	//defer stdout.Close()
	//writer := bufio.NewWriterSize(stdout, 16 * 1024 * 1024)
	statusQuery := os.Getenv("STATUSQUERY")
	parentIdTemp, err := strconv.ParseInt(strings.TrimSpace(pID), 10, 64)
	checkError(err)
	parentId := int32(parentIdTemp)
	result := avgRotorSpeed(statusQuery, parentId)
	fmt.Printf("%d\n", result)
	//writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
