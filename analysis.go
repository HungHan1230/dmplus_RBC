package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"go4.org/sort"
)

func AverageNodes() {
	csvfile, err := os.Open("./nodes_snapshots_reverse_forchurn.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	var count, total int

	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] == "1590940575" {
			break
		}
		tmp, _ := strconv.Atoi(record[1])
		total += tmp
		count++
	}

	fmt.Println("Average nodes in this period: ", total/count)

}

func ReadRequestsAnalysis(Type_path string, path string, outputName string) {
	names := GetFilesName(path)

	// Open the file
	csvfile, err := os.Open("/home/hank/go/read.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	var tmp_timestamp time.Time
	i, j := 0, 0
	access_status := make(map[string]int)
	record_map := make(map[string]map[string]int)
	// Iterate through the records
	var requests []string
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		layout := "2006-01-02 15:04:05"
		current_timestamp, err := time.Parse(layout, record[1])
		if i == 0 {
			tmp_timestamp = current_timestamp
		} else {
			diff := current_timestamp.Sub(tmp_timestamp)
			// fmt.Println(diff)
			if diff.Minutes() >= 5 {
				tmp := make(map[string]int)
				fmt.Println(current_timestamp)
				fmt.Println("reading...", path+names[j]+"_states.json")
				counter := calculateAccessRate(requests, path+names[j]+"_states.json")
				access_status["success"] += counter["success"]
				tmp["success"] += access_status["success"]
				failure := len(requests) - counter["success"]
				if failure > 0 {
					access_status["failure"] += len(requests) - counter["success"]
					tmp["failure"] += access_status["failure"]
				} else {
					access_status["failure"] += 0
					tmp["failure"] += access_status["failure"]
				}

				record_map[names[j]] = tmp

				fmt.Println(access_status)
				j++

				tmp_timestamp = current_timestamp
				requests = []string{}
			} else {
				requests = append(requests, record[3])
				// fmt.Println(requests)
			}

		}

		i++
	}
	fmt.Println("total access result: ", access_status)
	// toCSV(record_map, outputName)
	s1 := float64(access_status["success"])
	f1 := float64(access_status["failure"])
	fmt.Println("failure / total: ", f1/(f1+s1))
	Fault_rate := fmt.Sprintf("%f", f1/(f1+s1))

	i1 := strconv.Itoa(access_status["success"])
	i2 := strconv.Itoa(access_status["failure"])
	output := Type_path + "," + Fault_rate + "," + i1 + "," + i2

	// output fault rate to analysis.csv
	appendToCSV_pure_multiple(output, "analysis.csv")
	writeAnalysisToCSV(record_map, "./CSVs/"+outputName)
}

func writeAnalysisToCSV(access_record map[string]map[string]int, filename string) {
	var keys []int
	for key, _ := range access_record {
		timestamp, _ := strconv.Atoi(key)
		keys = append(keys, timestamp)
	}
	// fmt.Println("keys:", keys)
	sort.Ints(keys)

	for index := range keys {
		key := strconv.Itoa(keys[index])
		success := strconv.Itoa(access_record[key]["success"])
		failure := strconv.Itoa(access_record[key]["failure"])
		// check if the file exists
		_, err := os.Open(filename)
		if err != nil {
			os.Create(filename)
		}

		var path = filename
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		// defer f.Close()

		w := csv.NewWriter(f)

		var tocsv [][]string
		tocsv = append(tocsv, []string{key, success, failure})

		w.WriteAll(tocsv)

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
		// Replace defer f.Close() with f.Close()
		f.Close()
	}

}

// func toCSV(access_record map[string]map[string]int, filename string) {
// 	// layout := "2006-01-02 15:04:05"
// 	for k1, v1 := range access_record {
// 		success := strconv.Itoa(v1["success"])
// 		failure := strconv.Itoa(v1["failure"])
// 		str := success + "," + failure
// 		appendToCSV_pure(k1, str, "/home/hank/go/"+filename)
// 	}

// }

func calculateAccessRate(requests []string, file string) map[string]int {
	counter := make(map[string]int)
	state := make(map[string][]string)
	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &state)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
		// return err
	}
	// fmt.Println(state)
	// fmt.Println(requests)
	// for k, v := range state {
	// 	// fmt.Println("finding match....",k)
	// 	if v == nil || len(v) == 0 {
	// 		delete(state, k)
	// 	} else {
	// 		fmt.Println(len(requests))
	// 		for i := range requests {
	// 			for j := range v {
	// 				if "blk"+requests[i] == v[j] {
	// 					counter["success"] += 1
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	var tonext bool
	for i := range requests {
		tonext = false
		for k, v := range state {
			if v == nil || len(v) == 0 {
				delete(state, k)
			} else {
				for j := range v {
					if "blk"+requests[i] == v[j] {
						counter["success"] += 1
						tonext = true
						break
					}
				}
			}
			if tonext {
				break
			}
		}
	}

	return counter

}

// func ReadRequest_withoutTimeLimit(path string) {
// 	fmt.Println("Get file names")
// 	names := GetFilesName(path)
// 	fmt.Println("Get Requests")
// 	requests := GetRequestsFromCSV()

// 	// Open the file
// 	csvfile, err := os.Open("/home/hank/go/read_old.csv")
// 	if err != nil {
// 		log.Fatalln("Couldn't open the csv file", err)
// 	}
// 	// Parse the file
// 	r := csv.NewReader(csvfile)
// 	var requests_slice []string
// 	access_status := make(map[string]int)

// 	// Iterate through the records
// 	for {
// 		// Read each record from csv
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// fmt.Println(record)
// 		requests_slice = append(requests_slice, record[3])

// 		for j := range names {
// 			state := make(map[string][]string)
// 			byteValue, err := ioutil.ReadFile(path + names[j] + "_states.json")
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			err = json.Unmarshal(byteValue, &state)
// 			if err != nil {
// 				fmt.Println("something wrong while parsing json!")
// 				// return err
// 			}
// 			for k, v := range state {
// 				// fmt.Println("finding match....",k)
// 				if v == nil || len(v) == 0 {
// 					delete(state, k)
// 				} else {
// 					for in := range v {
// 						if v[in] == requests[j] {
// 							access_status["success"] += 1
// 							break
// 						}
// 					}

// 				}
// 			}

// 		}
// 	}

// 	fmt.Println("totoal requests: ", len(requests))
// 	fmt.Println("total success: ", access_status["success"])
// 	fmt.Println("total failure: ", len(requests)-access_status["success"])

// }

// func GetRequestsFromCSV() []string {
// 	// Open the file
// 	csvfile, err := os.Open("/home/hank/go/read_old.csv")
// 	if err != nil {
// 		log.Fatalln("Couldn't open the csv file", err)
// 	}
// 	// Parse the file
// 	r := csv.NewReader(csvfile)
// 	var requests_slice []string

// 	// Iterate through the records
// 	for {
// 		// Read each record from csv
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// fmt.Println(record)
// 		requests_slice = append(requests_slice, record[3])
// 	}
// 	return requests_slice
// }
