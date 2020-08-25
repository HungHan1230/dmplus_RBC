package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	timestamp_d   int64
	total_nodes_d int
}

// type timestampjson struct {
// 	total_nodes string
// 	IP_address  []string
// }

var base string = "https://bitnodes.io/api/v1/snapshots/"
var has_next bool
var next string

func GetFilesName(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var k []string
	for _, file := range files {
		// fmt.Println(file.Name())
		k = append(k, file.Name()[:10])
		// fmt.Println(file.Name()[:10])
		// fmt.Println(file.Name()[:10])
	}
	return k

}

func GetFromNodesJson() {
	// Open the file
	csvfile, err := os.Open("nodes_churn.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// m := make(map[string]int)
	var data []mydata

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])

		//My experiments only need a month, I don't need more..
		if record[0] == "1590941023" {
			break
		}

		var element mydata
		element.timestamp, _ = strconv.ParseInt(record[0], 10, 64)
		element.nodes, _ = strconv.Atoi(record[1])
		element.churn_r, _ = strconv.ParseFloat(record[3], 64)
		element.churn_n, _ = strconv.Atoi(record[4])
		element.add_n, _ = strconv.Atoi(record[5])
		element.add_r, _ = strconv.ParseFloat(record[6], 64)
		// fmt.Println(element)
		data = append(data, element)

	}	

}

func GetNodeSnapshots() {
	// next = "https://bitnodes.io/api/v1/snapshots/"
	next = base
	has_next = true
	i := 1
	for has_next {
		fmt.Println("Page", i)
		GetSnapshots(next)
		i++
	}
	// timestamplist := GetTimestamps()

	// for i = 0; i < len(timestamplist); i++ {
	// 	fmt.Println("downlaoding timestamp: ", timestamplist[i])
	// 	GetSnapshotsWithTimestamps(timestamplist[i])
	// }

}
func appendToJson(nodesjson map[string]json.RawMessage, totalNodes string, filename string) {
	outputmap := make(map[string]map[string]string)

	for k, v := range nodesjson {
		objectmap := make(map[string]string)
		//The cool thing is that here I set int slice, it only read the int data into this slice
		var temp []int
		var tmp []string
		json.Unmarshal(v, &temp)
		json.Unmarshal(v, &tmp)

		objectmap["start_time"] = strconv.Itoa(temp[2]) // result: s = "-18"temp[2]
		objectmap["Timezone"] = tmp[10]
		objectmap["ASN"] = tmp[11]
		objectmap["Organization"] = tmp[12]
		// tmpstr = tmpstr[:len(tmpstr)-1]
		// objectmap["info"] = tmpstr
		outputmap[k] = objectmap

	}

	file, err := json.Marshal(outputmap)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	var path string = "../../node_jsons_reverse/" + filename + ".json"
	// Write to file
	_ = ioutil.WriteFile(path, file, 0644)

}
func GetSnapshotsWithTimestamps(timestamp string) (judge bool) {
	url := base + "/" + timestamp
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//request body, which is byte[]
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// a map container to decode the JSON structure into
	result := make(map[string]json.RawMessage)

	// unmarschal JSON
	err = json.Unmarshal(sitemap, &result)
	if err != nil {
		return
	}

	// if request got throttle, return false to terminate the program.
	var jud bool
	for k, _ := range result {
		if k != "detail" {
			jud = true
			break
		} else {
			jud = false
		}
	}
	// fmt.Println(result)
	if jud {
		total_nodes := result["total_nodes"]
		nodes := result["nodes"]

		nodesjson := make(map[string]json.RawMessage)
		err = json.Unmarshal(nodes, &nodesjson)
		appendToJson(nodesjson, string(total_nodes), timestamp)

	}
	return jud

}

func GetTimestamps() (list []string) {
	// Open the file
	csvfile, err := os.Open("nodes_snapshots_reverse.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// m := make(map[string]int)
	var timestamplist []string
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
		timestamplist = append(timestamplist, record[0])

	}
	return timestamplist
}

func GetChurn() []mydata {
	// Open the file
	csvfile, err := os.Open("nodes_churn.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// m := make(map[string]int)
	var data []mydata

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])

		//My experiments only need a month, I don't need more..
		if record[0] == "1590941023" {
			break
		}

		var element mydata
		element.timestamp, _ = strconv.ParseInt(record[0], 10, 64)
		element.nodes, _ = strconv.Atoi(record[1])
		element.churn_r, _ = strconv.ParseFloat(record[3], 64)
		element.churn_n, _ = strconv.Atoi(record[4])
		element.add_n, _ = strconv.Atoi(record[5])
		element.add_r, _ = strconv.ParseFloat(record[6], 64)
		// fmt.Println(element)
		data = append(data, element)

	}
	return data

}

func appendToCSV(data Data, csvfile string) {
	// check if nodes.csv exists
	_, err := os.Open(csvfile)
	if err != nil {
		// fmt.Println(os.IsNotExist(err)) //true  證明檔案已經存在
		// fmt.Println(err)                //open widuu.go: no such file or directory
		os.Create(csvfile)
	}

	var path = csvfile
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	var tocsv [][]string
	tocsv = append(tocsv, []string{strconv.FormatInt(data.timestamp_d, 10), strconv.Itoa(data.total_nodes_d)})

	w.WriteAll(tocsv)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
func appendToCSV_pure(data1 string, data2 string, csvfile string) {
	// check if nodes.csv exists
	of, err := os.Open(csvfile)
	if err != nil {
		// fmt.Println(os.IsNotExist(err)) //true  證明檔案已經存在
		// fmt.Println(err)                //open widuu.go: no such file or directory
		os.Create(csvfile)
	}

	var path = csvfile
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	// defer f.Close()

	w := csv.NewWriter(f)

	var tocsv [][]string
	tocsv = append(tocsv, []string{data1, data2})

	w.WriteAll(tocsv)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	f.Close()
	of.Close()

}
func appendToCSV_pure_multiple(data string, csvfile string) {
	str := strings.Split(data,",")
	// check if nodes.csv exists
	of, err := os.Open(csvfile)
	if err != nil {
		// fmt.Println(os.IsNotExist(err)) //true  證明檔案已經存在
		// fmt.Println(err)                //open widuu.go: no such file or directory
		os.Create(csvfile)
	}

	var path = csvfile
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	// defer f.Close()

	w := csv.NewWriter(f)

	var tocsv [][]string	
	
	tocsv = append(tocsv, str)

	w.WriteAll(tocsv)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	f.Close()
	of.Close()

}

func GetSnapshots(nextURL string) {
	fmt.Printf("Downloading from: %s \n", nextURL)

	res, err := http.Get(nextURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//request body, which is byte[]
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// parse json
	var result map[string]interface{}
	err = json.Unmarshal(sitemap, &result)
	if err != nil {
		return
	}
	next_url := result["next"]
	// fmt.Println(next_url)

	if next_url != nil {
		results := result["results"].([]interface{})
		// fmt.Println("results: ", results)
		for _, obj := range results {
			object := obj.(map[string]interface{})
			timestamp := int64(object["timestamp"].(float64))
			total_nodes := int(object["total_nodes"].(float64))

			data := Data{timestamp_d: timestamp, total_nodes_d: total_nodes}
			appendToCSV(data, "nodes_snapshots.csv")
		}

		str := fmt.Sprintf("%v", next_url)
		next = str

	} else {
		has_next = false
	}

}
