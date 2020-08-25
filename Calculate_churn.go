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
	"gonum.org/v1/plot/plotter"
)

type mydata struct {
	timestamp int64
	nodes     int
	churn_r   float64
	churn_n   int
	add_n     int
	add_r     float64
}

func calculateChurn_Paper() {
	names := GetFilesName("../../nodes/nodes_states")

	filelist := make(map[int][]int64)
	var day int = 1
	var previous, current int64
	for index := range names {
		if index == 0 {
			previous, _ = strconv.ParseInt(names[index], 10, 64)
			filelist[day] = append(filelist[day], previous)
		} else {
			current, _ = strconv.ParseInt(names[index], 10, 64)

			t1 := time.Unix(previous, 0)
			t2 := time.Unix(current, 0)
			diff := t2.Sub(t1)
			// fmt.Println(diff)
			if diff.Hours()/24 < 1 {
				filelist[day] = append(filelist[day], current)

			} else {
				previous = current
				day++
			}
		}
	}
	file, err := json.Marshal(filelist)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	var jsonpath string = "./node_jsons_days.json"
	// Write to file
	_ = ioutil.WriteFile(jsonpath, file, 0644)

	// fmt.Println(filelist)
	// allmap := make(map[int]map[string]string)
	// var path string = "../../nodes/nodes_states/"
	// // var pre_map map[string]string
	// for i := 1; i <= 32; i++ {
	// 	for j := 0; j < len(filelist[i]); j++ {
	// 		filename := filelist[i][j]
	// 		str_filename := strconv.FormatInt(filename, 10) + "_states.json"
	// 		allmap[i] = ReadStates(path + str_filename)

	// 		//read each day each csv file
	// 		// if j==0{
	// 		// 	pre_map = ReadStates(str_filename)
	// 		// }else{
	// 		// 	cur_map := ReadStates(str_filename)

	// 		// 	fmt.Println(cur_map)
	// 		// }

	// 	}
	// 	fmt.Println("day: ",len(allmap))

	// }

	// first_day, _ := strconv.ParseInt(names[0], 10, 64)
	// last_day, _ := strconv.ParseInt(names[len(names)-1], 10, 64)
	// first := time.Unix(first_day, 0)
	// last := time.Unix(last_day, 0)
	// fmt.Println("first day:", first)
	// fmt.Println("last day:", last)
	// d := last.Sub(first)
	// fmt.Println(d)
	// fmt.Println(d.Hours() / 24)

	// //check size
	// var tmp int
	// for _, v := range filelist {
	// 	tmp += len(v)
	// 	// fmt.Println(len(v))
	// }
	// fmt.Println(tmp)
	// fmt.Println(len(filelist))
	// fmt.Println("1 ", filelist[1])
	// fmt.Println("last ", filelist[len(filelist)])
}
func CalculateEachSession() {
	names := GetFilesName("../../nodes/nodes_states")
	IP_States := make(map[string]map[string][]string)
	for i := range names {
		fmt.Println("Processing...", names[i])
		byteValue, err := ioutil.ReadFile("../../nodes/nodes_states/" + names[i] + "_states.json")
		if err != nil {
			fmt.Println("something wrong!")
			// return err
		}

		var datamap map[string]string
		err = json.Unmarshal(byteValue, &datamap)
		if err != nil {
			fmt.Println("something wrong while parsing json!")
			// return err
		}
		if i == 0 {
			for k, _ := range datamap {
				tmp := make(map[string][]string)
				tmp["on"] = append(tmp["on"], names[i])
				IP_States[k] = tmp
			}
		} else {
			for k, v := range datamap {
				if val, ok := IP_States[k]; ok {
					//do something here
					// fmt.Println(val)

					if v == "on" && len(IP_States[k]["on"]) == len(val["on"]) {
						IP_States[k]["on"] = append(IP_States[k]["on"], names[i])
					} else if v == "off" && len(IP_States[k]["off"]) < len(IP_States[k]["on"]) {
						IP_States[k]["off"] = append(IP_States[k]["off"], names[i])
					}

				} else if v == "on" {
					tmp := make(map[string][]string)
					tmp["on"] = append(tmp["on"], names[i])
					IP_States[k] = tmp
					// IP_States[k]["on"] = append(IP_States[k]["on"], names[i])
				}

				//maybe delete the element after adding to IP_States, probably may save some memory...?!
			}

		}

	}
	// fmt.Println(IP_States)
	file, err := json.Marshal(IP_States)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	var jsonpath string = "./nodes_jsons_sessions.json"
	// Write to file
	_ = ioutil.WriteFile(jsonpath, file, 0644)

}
func Calculate_totalChurn() {
	csvfile, err := os.Open("Churn_Rate.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	var count float64
	var lineinCSV int
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
		// fmt.Println(record)
		if record[1] == "+Inf" {
			continue
		} else {
			tmp, _ := strconv.ParseFloat(record[1], 64)

			if tmp > 10 {
				fmt.Println(tmp)
			}
			count += tmp
			lineinCSV++
		}

	}
	fmt.Println("average churn rate: ", count/float64(lineinCSV))

}

func Calculate_averageChurn_Simple() {
	byteValue, err := ioutil.ReadFile("./nodes_jsons_sessions.json")
	if err != nil {
		fmt.Println("something wrong!")
		// return err
	}
	var datamap map[string]map[string][]string
	err = json.Unmarshal(byteValue, &datamap)
	if err != nil {
		fmt.Println(err)
	}
	var churn_count int
	var MaxPreMaxInner_count int
	var MaxPreIP string
	for k, v := range datamap {
		var MaxInner_count int
		for k1, v1 := range v {
			if k1 == "off" {
				churn_count += len(v1)
				MaxInner_count += len(v1)
				break
			}
		}
		if MaxInner_count > MaxPreMaxInner_count{
			MaxPreMaxInner_count = MaxInner_count
			MaxPreIP = k		
		}
	}
	fmt.Println("Max Churns: ", MaxPreMaxInner_count)
	fmt.Println("Max ip: ", MaxPreIP)
	fmt.Println("length of churn ip: ", len(datamap))
	fmt.Println("total churns: ", churn_count)
	fmt.Println("average churns: ", float64(churn_count/len(datamap)))

// 	total ip:  25913
//  total churns:  15932155
//  average churns:  614
}

func ReadIP_States() {
	byteValue, err := ioutil.ReadFile("./nodes_jsons_sessions.json")
	if err != nil {
		fmt.Println("something wrong!")
		// return err
	}

	var datamap map[string]map[string][]string
	err = json.Unmarshal(byteValue, &datamap)
	if err != nil {
		fmt.Println(err)
	}

	// var ip_list []string
	// var churn_list []float64
	var Churn_Rate float64

	for k, v := range datamap {
		fmt.Println("processing...", k)
		// ip_list = append(ip_list, k)
		var duration_list []float64
		for i := range v["on"] {
			if len(v["on"]) != len(v["off"]) {
				continue
			} else {
				on, _ := strconv.ParseInt(v["on"][i], 10, 64)
				off, _ := strconv.ParseInt(v["off"][i], 10, 64)
				t_on := time.Unix(on, 0)
				t_off := time.Unix(off, 0)
				diff := t_off.Sub(t_on).Hours() / 24
				duration_list = append(duration_list, diff)

			}

		}
		var count float64
		for i := range duration_list {
			count += duration_list[i]
		}

		Churn_Rate = 1 / count
		// output key, churn_rate
		var ipaddress string = k
		var data_Churn = fmt.Sprintf("%f", Churn_Rate)
		appendToCSV_pure(ipaddress, data_Churn, "Churn_Rate.csv")
	}

	//output ip, churn rate csv

	// fmt.Println(len(datamap))

}

func CalculateEachSession_Old() {
	byteValue, err := ioutil.ReadFile("./node_jsons_days.json")
	if err != nil {
		fmt.Println("something wrong!")
		// return err
	}

	var datamap map[string][]int
	err = json.Unmarshal(byteValue, &datamap)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
		// return err
	}

	for _, v := range datamap {
		first_state := make(map[string]string)
		for i := range v {
			str_filename := strconv.Itoa(v[i])
			byteValue, err = ioutil.ReadFile(str_filename)
			if err != nil {
				fmt.Println("something wrong!")
				// return err
			}
			state_map := make(map[string]string)
			err = json.Unmarshal(byteValue, &state_map)
			if err != nil {
				fmt.Println("something wrong while parsing json!")
				// return err
			}
			if i == 0 {
				first_state = state_map
			} else {
				// for k, _ := range first_state {
				// 	if val, ok := state_map[k]; ok {
				// 		//do something here
				// 	}
				// }

			}

		}
		fmt.Println(first_state)

	}

}

func ReadStates(str_filename string) map[string]string {
	// Open the file
	byteValue, err := ioutil.ReadFile(str_filename)
	if err != nil {
		fmt.Println("something wrong!")
		// return err
	}

	var datamap map[string]string
	err = json.Unmarshal(byteValue, &datamap)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
		// return err
	}

	return datamap
}

func caluculateCount() {
	data := GetChurn()
	m := make(map[int]int)
	for i := 0; i < len(data); i++ {
		if data[i].churn_n == 0 {
			continue
		}
		if m[data[i].churn_n] != 0 {
			m[data[i].churn_n] = m[data[i].churn_n] + 1
		} else {
			m[data[i].churn_n] = 1
		}

	}
	fmt.Println("record the frequencies of churn nodes: ", m)

	var total int = 0
	m2 := make(map[int]int)
	var keys []int

	for k, v := range m {
		keys = append(keys, k)
		total += v
	}
	// sort the keys slice
	sort.Ints(keys)
	// fmt.Println(keys)

	// prepare points for plotting
	var points plotter.XYs
	for i := 0; i < len(keys); i++ {
		// fmt.Println(keys[i])
		if keys[i] > 300 {
			continue
		}
		m2[keys[i]] = total

		points = append(points, struct{ X, Y float64 }{float64(keys[i]), float64(total)})
		total -= m[keys[i]]
	}
	fmt.Println("The cumulative calculation of churn nodes: ", m2)
	// plot cumulative graph of churn nodes
	plottest(points)
}

func calculateDailyChurnCumulative() {
	data := GetChurn()
	day_map := make(map[string][]float64)

	var previous, unixTimeUTC time.Time
	layout := "2006-01-02"

	for i := 0; i < len(data); i++ {
		if data[i].churn_n == 0 {
			// fmt.Println("churn nodes = 0")
			continue
		}
		now := data[i].timestamp
		unixTimeUTC = time.Unix(now, 0)
		judgement := unixTimeUTC.Sub(previous)
		key := unixTimeUTC.Format(layout)
		// fmt.Println(unixTimeUTC.Date())
		// fmt.Println(previous.Date())

		if i == 0 {
			previous = unixTimeUTC
			day_map[key] = []float64{float64(data[i].churn_n)}
		} else if judgement.Hours()/24 >= 1 {
			// } else if unixTimeUTC.Date() != previous.Date() {
			previous = unixTimeUTC
			// key := unixTimeUTC.Format(layout)
			if day_map[key] == nil {
				day_map[key] = []float64{float64(data[i].churn_n)}
			} else {
				day_map[key] = append(day_map[key], float64(data[i].churn_n))
			}

		} else {
			p_key := previous.Format(layout)
			day_map[p_key] = append(day_map[p_key], float64(data[i].churn_n))
		}

	}
	fmt.Println("dat_map: ", day_map)

}

func calculateChurnCumulative() {

	data := GetChurn()
	// fmt.Println(data)
	percentage_m := make(map[string]int)
	var judge float64 = 0.0

	for i := 0; i < 100; i++ {
		for j := 0; j < len(data); j++ {
			if data[j].churn_r > judge {
				key := strconv.FormatFloat(judge, 'g', 5, 64)
				percentage_m[key] = percentage_m[key] + 1
			}
		}
		// fmt.Println(judge)
		judge += 0.01
	}
	fmt.Println(percentage_m)
}

func calculateChurn() {
	// Open the file
	csvfile, err := os.Open("nodes_snapshots_reverse_forchurn.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// mycode
	var counter int
	counter = 1
	var tmptimestamp, tmpnodes, previoustimestamp, previousnodes, churn_rate, churn_nodes, add_nodes, add_rate float64
	var dataArr []mydata
	var points_churn, points_churn_nodes plotter.XYs
	var points_add plotter.XYs
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

		if record[0] == "1590940575" {
			break
		}
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
		if counter == 1 {
			// tmptimestamp, err = strconv.ParseInt(record[0], 10, 64)
			tmptimestamp, err = strconv.ParseFloat(record[0], 64)
			tmpnodes, err = strconv.ParseFloat(record[1], 64)
			//initialization
			previoustimestamp = tmptimestamp
			previousnodes = tmpnodes
			churn_rate = 0.0
			churn_nodes = 0
			add_nodes = 0
			add_rate = 0.0
			// first record
			var tmpdata mydata
			tmpdata.timestamp = int64(previoustimestamp)
			tmpdata.churn_r = churn_rate
			tmpdata.churn_n = 0
			tmpdata.add_n = int(add_nodes)
			tmpdata.nodes = int(previousnodes)
			tmpdata.add_r = add_rate
			dataArr = append(dataArr, tmpdata)
		} else {
			tmptimestamp, err = strconv.ParseFloat(record[0], 64)
			tmpnodes, err = strconv.ParseFloat(record[1], 64)

			var judgement float64
			judgement = previousnodes - tmpnodes

			//If preivousnodes greater than tmpnodes, it means that there were some nodes disconnected (churn).
			if judgement > 0 {
				churn_rate = judgement / previousnodes
				churn_nodes = judgement
				add_nodes = 0
				add_rate = 0
			} else {
				churn_rate = 0
				churn_nodes = 0
				add_nodes = tmpnodes - previousnodes
				add_rate = add_nodes / previousnodes
			}

			previoustimestamp = tmptimestamp
			previousnodes = tmpnodes

			var tmpdata mydata
			tmpdata.timestamp = int64(previoustimestamp)
			tmpdata.churn_r = churn_rate
			tmpdata.churn_n = int(churn_nodes)
			tmpdata.add_n = int(add_nodes)
			tmpdata.nodes = int(previousnodes)
			tmpdata.add_r = add_rate
			dataArr = append(dataArr, tmpdata)
			points_churn = append(points_churn, struct{ X, Y float64 }{previoustimestamp, churn_rate})
			points_churn_nodes = append(points_churn_nodes, struct{ X, Y float64 }{previoustimestamp, churn_nodes})
			points_add = append(points_add, struct{ X, Y float64 }{previoustimestamp, add_rate})

		}
		counter++
	}
	// fmt.Println(dataArr)
	writeChurnToCSV(dataArr)

	plotchurn(points_churn)
	plotchurn(points_churn_nodes)
	plotadd_r(points_add)
	// plotchurn_V2(points)

}

func writeChurnToCSV(data []mydata) {
	// check if nodes.csv exists
	_, err := os.Open("nodes_churn.csv")
	if err != nil {
		// fmt.Println(os.IsNotExist(err)) //true  證明檔案已經存在
		// fmt.Println(err)                //open widuu.go: no such file or directory
		os.Create("nodes_churn.csv")
	}

	var path = "nodes_churn.csv"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	var tocsv [][]string
	for i := 0; i < len(data); i++ {
		originalTimestamp := strconv.FormatInt(data[i].timestamp, 10)
		// originalNodes := string(data[i].nodes) // so interesting
		originalNodes := strconv.Itoa(data[i].nodes)
		var tmptime int64
		tmptime = data[i].timestamp
		unixTimeUTC := time.Unix(tmptime, 0) //gives unix time stamp in utc
		// tocsv = append(tocsv, []string{strconv.FormatInt(data[i].timestamp, 10), fmt.Sprintf("%f", data[i].churn_r), strconv.Itoa(data[i].add_n)})
		tocsv = append(tocsv, []string{originalTimestamp, originalNodes, unixTimeUTC.Format("2006-01-02 15:04:05"), fmt.Sprintf("%f", data[i].churn_r), strconv.Itoa(data[i].churn_n), strconv.Itoa(data[i].add_n), fmt.Sprintf("%f", data[i].add_r)})
	}

	w.WriteAll(tocsv)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
