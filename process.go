package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

var globalstate map[string]string

func RunProcess(Type_Path string, NumOf_ProversNodes int, repairOrNot bool, storagelimit_pernode float64, replication_factor float64, request_filename string) {
	// --------get the results in each timestamps--------
	// assign blk.dat file to the first timestamp.
	assignblkToFirst(Type_Path, NumOf_ProversNodes, storagelimit_pernode, replication_factor)
	// // record the lost and repair in each timestamp
	calculateLostAndRepair(Type_Path, NumOf_ProversNodes, repairOrNot, storagelimit_pernode, replication_factor)

	// --------Get the request results--------
	numofprovers := strconv.Itoa(NumOf_ProversNodes)
	// storagelimit := fmt.Sprintf("%f", storagelimit_pernode)
	storagelimit := strconv.Itoa(int(storagelimit_pernode))
	output_path := "../../nodes/results/" + numofprovers + "provers/" + storagelimit + "G/" + Type_Path + "nodes_withBlk_state/"
	// ReadRequest(output_path, "10Baseline.csv")
	ReadRequestsAnalysis(Type_Path, output_path, request_filename)
}

func calculateLostAndRepair(Type_Path string, NumOf_ProversNodes int, repairOrNot bool, storagelimit_pernode float64, replication_factor float64) {
	numofprovers := strconv.Itoa(NumOf_ProversNodes)
	// storagelimit := fmt.Sprintf("%f", storagelimit_pernode)
	storagelimit := strconv.Itoa(int(storagelimit_pernode))
	output_path := "../../nodes/results/" + numofprovers + "provers/" + storagelimit + "G/" + Type_Path + "nodes_withBlk_state/"

	names_Blk := GetFilesName(output_path)
	var path_Status string = "../../nodes/nodes_states/"
	names_States := GetFilesName(path_Status)

	first_state := make(map[string][]string)
	byteValue, err := ioutil.ReadFile(output_path + names_Blk[0] + "_states.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &first_state)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
	}
	// Initialization of variables
	repair := make(map[string][]string)
	states := make(map[string]string)
	var previous_timestamp time.Time
	init_timestamp, _ := strconv.ParseInt(names_Blk[0], 10, 64)
	previous_timestamp = time.Unix(init_timestamp, 0)

	// record each iteration of states
	for i := 1; i < len(names_States); i++ {
		// convert timestamp and assigned by current_timestamp
		current, err := strconv.ParseInt(names_States[i], 10, 64)
		current_timestamp := time.Unix(current, 0)

		byteValue, err := ioutil.ReadFile(path_Status + names_States[i] + "_states.json")
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(byteValue, &states)
		if err != nil {
			fmt.Println("something wrong while parsing json!")
		}
		// repair time = 10mins, first examines the timestamp has already 10 mins (two snapshots), then add the repair element to the nodes that isn't preversed blk
		diff := current_timestamp.Sub(previous_timestamp)
		// fmt.Println(diff.Minutes())
		if repairOrNot {
			if diff.Minutes() > 10 {
				for k, v := range repair {
					for kf, vf := range first_state {
						if len(vf) == 0 {
							first_state[kf] = v
							delete(repair, k)
							break
						}
					}
				}
				// fmt.Println("repair time: ", names_States[i])
				// fmt.Println(repair)

				//remember the timestamp as preivous_timestamp for the next iteration
				forprevious_timestamp, _ := strconv.ParseInt(names_States[i], 10, 64)
				previous_timestamp = time.Unix(forprevious_timestamp, 0)
			}

		}

		// Iterate elements in states, if the value is off or left, remove from the first_state, if the value is join, add to first_state with empty string slice.
		for k, v := range states {
			if v == "off" || v == "left" {
				repair[k] = first_state[k]
				delete(first_state, k)
			} else if v == "join" {
				first_state[k] = []string{}
			}

		}
		//output this state
		outputStates(first_state, names_States[i], output_path)
	}

}
func outputStates(state_map map[string][]string, name string, Directory_Path string) {
	//out the assign result to the folder nodes_withBlk_state
	file, err := json.Marshal(state_map)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	// Write to file
	_ = ioutil.WriteFile(Directory_Path+name+"_states.json", file, 0644)
}

func assignblkToFirst(Type_Path string, simulate_n int, storagelimit_pernode float64, replication_factor float64) {
	names := GetFilesName("../../nodes/nodes_states")
	var path string = "../../nodes/nodes_states"
	first_state := make(map[string]string)
	byteValue, err := ioutil.ReadFile(path + "/" + names[0] + "_states.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &first_state)
	if err != nil {
		fmt.Println(err)
	}

	//get assignments
	// var simulate_n int = 500
	assign_nums := simulate(simulate_n, storagelimit_pernode, replication_factor)
	// fmt.Println("assign_nums length: ", len(assign_nums))

	//get keys from state.json to obtain the appeared IP addresses.
	var keys []string
	for k, _ := range first_state {
		keys = append(keys, k)
	}

	var assign_keys []int
	for k, _ := range assign_nums {
		assign_keys = append(assign_keys, k)
	}

	// assign output
	output := make(map[string][]string)
	survivors := map[string]int{"101.100.174.240:8333": 1}
	// survivors := map[string]int{"101.100.174.240:8333": 1, "103.85.190.218:20000": 1, "109.236.81.138:58333": 1, "111.221.46.63:8333": 1, "16.203.185.207:8446": 1, "136.243.142.20:8333": 1, "130.185.77.105:8333": 1, "136.143.48.128:8333": 1, "144.76.116.8:8446": 1, "142.91.11.100:8333": 1, "13.228.130.195:8333": 1, "139.59.247.235:8333": 1, "142.4.211.170:8333": 1, "142.93.109.28:8333": 1, "142.93.137.232:8333": 1, "144.76.108.6:8444": 1, "151.80.33.185:8333": 1, "144.76.117.234:8433": 1, "160.20.145.62:8333": 1, "165.22.112.46:8333": 1}
	// survivors := map[string]int{"101.100.174.240:8333": 1, "103.85.190.218:20000": 1, "109.236.81.138:58333": 1, "111.221.46.63:8333": 1, "16.203.185.207:8446": 1}
	// for baseline
	// survivors := ReadSurvivors(100)
	// for R2
	// survivors := ReadSurvivors(800)
	// for R4 R8
	// survivors := ReadSurvivors(500)

	var survivors_keys []string
	for k, _ := range survivors {
		survivors_keys = append(survivors_keys, k)
	}
	// fmt.Println("keys: ", keys)
	num := 0
	for _, v := range assign_nums {
		//old
		// output[keys[num]] = v

		// new
		// index := num * 5
		// // fmt.Println("index", index)
		// // fmt.Println(keys[index])
		// output[keys[index]] = v

		// new with survivors
		// if num < len(survivors_keys) {
		if num < len(survivors_keys) {
			output[survivors_keys[num]] = v
			fmt.Println("eternal storage: ", output)
		} else {
			//index := num * 5 before 2000
			index := num
			output[keys[index]] = v
		}

		num++
	}

	//output the assign result to the folder nodes_withBlk_state
	file, err := json.Marshal(output)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	numofprovers := strconv.Itoa(simulate_n)
	// storagelimit := fmt.Sprintf("%f", storagelimit_pernode)
	storagelimit := strconv.Itoa(int(storagelimit_pernode))
	output_path := "../../nodes/results/" + numofprovers + "provers/" + storagelimit + "G/" + Type_Path + "nodes_withBlk_state/"

	mkdir(output_path)

	// Write to file
	_ = ioutil.WriteFile(output_path+names[0]+"_states.json", file, 0644)

	assignblkToFirst_withEmpty(Type_Path, numofprovers, storagelimit, replication_factor)

}

func mkdir(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
}

func assignblkToFirst_withEmpty(Type_Path string, numofprovers string, storagelimit string, replication_factor float64) {
	namesofstates := GetFilesName("../../nodes/nodes_states")
	first_state := make(map[string]string)
	byteValue, err := ioutil.ReadFile("../../nodes/nodes_states/" + namesofstates[0] + "_states.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &first_state)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
	}
	path := "../../nodes/results/" + numofprovers + "provers/" + storagelimit + "G/" + Type_Path + "nodes_withBlk_state/"
	namesofBlk := GetFilesName(path)
	first_stateWithBlk := make(map[string][]string)
	byteValue, err = ioutil.ReadFile(path + namesofBlk[0] + "_states.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &first_stateWithBlk)
	if err != nil {
		fmt.Println("something wrong while parsing json!")
	}

	for k, _ := range first_state {
		for k1, _ := range first_stateWithBlk {
			if k1 == k {
				delete(first_state, k)
			}
		}
	}

	for k, _ := range first_state {
		first_stateWithBlk[k] = []string{}
	}

	//output the assign result to the folder nodes_withBlk_state
	file, err := json.Marshal(first_stateWithBlk)
	if err != nil {
		fmt.Println("something wrong while writing json!")
	}
	// Write to file
	_ = ioutil.WriteFile(path+namesofBlk[0]+"_states.json", file, 0644)
}

func ReadSurvivors(num int) map[string]string {
	suvivors := make(map[string]string)
	csvfile, err := os.Open("survivors.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		if num == 0{
			break
		}
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		suvivors[record[0]] = record[1]

		num--
	}
	return suvivors

}

func WhoIsAlwaysUp() {
	names := GetFilesName("../../nodes/nodes_states")
	var names_int []int
	for k, _ := range names {
		tmp, _ := strconv.Atoi(names[k])
		names_int = append(names_int, tmp)
	}
	sort.Ints(names_int)

	first := make(map[string]string)
	current := make(map[string]string)

	// read first state
	str := strconv.Itoa(names_int[0])
	byteValue, err := ioutil.ReadFile("../../nodes/nodes_states/" + str + "_states.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(byteValue, &first)
	if err != nil {
		fmt.Println(err)
	}

	for index := 1; index < len(names_int); index++ {
		str := strconv.Itoa(names_int[index])
		fmt.Println("reading....", str)
		byteValue, err := ioutil.ReadFile("../../nodes/nodes_states/" + str + "_states.json")
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(byteValue, &current)
		if err != nil {
			fmt.Println(err)
		}

		for k, v := range first {
			for k1, v1 := range current {
				if k == k1 {
					if v == v1 {
						break
					} else {
						delete(first, k)
					}
				}
			}

		}
		fmt.Println("current survivors: ", len(first))

	}

	fmt.Println("Number of Survivors: ", len(first))
	fmt.Println("survivor IP Addresses: ", first)
	for k, v := range first {
		appendToCSV_pure(k, v, "./survivors.csv")
	}

	// for _, v := range first {
	// 	fmt.Print(v + " ")
	// 	if v == "on" {
	// 		fmt.Println("yes")
	// 	} else if v == "left" {
	// 		fmt.Println("no")
	// 	}
	// }

}

func WhoIsAlwaysUp_test() {
	names := GetFilesName("../../nodes/test")
	var previous []string

	for i := 0; i < len(names); i++ {
		var current []string
		var path1 string = "../../nodes/test/" + names[i] + ".json"
		//read path1
		byteValue, err := ioutil.ReadFile(path1)
		if err != nil {
			fmt.Println(err)
			break
		}
		var result map[string]json.RawMessage
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			fmt.Println(err)
			break
		}
		if i == 0 {
			for k, _ := range result {
				previous = append(previous, k)
			}
		} else {
			for k, _ := range result {
				current = append(current, k)
			}
			whosUp := WhoisSurvive(previous, current)
			fmt.Println("num of nodes: ", len(whosUp))
			fmt.Println("survivors: ", whosUp)
			previous = current
		}

	}

}

func WhoIsAlwaysUp_() {
	names := GetFilesName("../../nodes/test")
	// var path1 string = "../../nodes/node_jsons_reverse/" + names[0] + ".json"
	// var path2 string = "../../nodes/node_jsons_reverse/" + "1588410937" + ".json"
	var result map[string]json.RawMessage
	var result2 map[string]json.RawMessage
	var keys1, keys2 []string
	for i := 0; i < len(names); i++ {
		var path1 string = "../../nodes/test/" + names[i] + ".json"
		//read path1
		byteValue, err := ioutil.ReadFile(path1)
		if err != nil {
			fmt.Println(err)
			fmt.Println("something wrong!")
			// return err
		}
		if i == 0 {
			err = json.Unmarshal(byteValue, &result)
			if err != nil {
				// fmt.Println("something wrong while parsing json!")
				fmt.Println(err)
				// return err
			}

		} else {
			err = json.Unmarshal(byteValue, &result2)
			if err != nil {
				// fmt.Println("something wrong while parsing json!")
				fmt.Println(err)
				// return err
			}

		}

		for k, _ := range result {
			keys1 = append(keys1, k)
		}
		for k, _ := range result2 {
			keys2 = append(keys2, k)
		}

	}
	whosUp := WhoisSurvive(keys1, keys2)
	fmt.Println(whosUp)
	fmt.Println(len(whosUp))

}

func WhoisSurvive(first []string, second []string) []string {
	var survivors []string

	for i := range first {
		for j := range second {
			if first[i] == second[j] {
				survivors = append(survivors, first[i])
				// fmt.Println("survivor: ", first[i])
				break
			}
		}
	}

	return survivors

}

func simulate(simulate_nodes int, storagelimit_pernode float64, replication_factor float64) map[int][]string {
	// rand.Seed(time.Now().UnixNano())
	var numofBlk float64 = 2116
	// var replication_factor float64 = 1
	// var storagelimit_pernode float64 = 40 //GB
	min := 0
	// simulate_nodes := 1000 //actually depends on how many prover nodes we want to simulate

	maximum_numofblk_pernode := math.Ceil(storagelimit_pernode * 1024 / 128)
	// fmt.Println(maximum_numofblk_pernode)

	datanode_m := make(map[int][]string)

	var tmpnumofblk int = int(numofBlk)
	var previous_num []int

	for i := 0; i < int(numofBlk); i++ {
		for j := 0; j < int(replication_factor); j++ {
			var num int
			// generate random num that is in th range of simulate_nodes and min
			num = rand.Intn(simulate_nodes-min) + min
			//To make sure there's no any nums are the same.
			if len(previous_num) != 0 {
				for {
					var check bool = false
					for index := range previous_num {
						if num == previous_num[index] {
							continue
						} else {
							check = true
						}
					}
					if !check {
						num = rand.Intn(simulate_nodes-min) + min
						continue
					} else {
						previous_num = append(previous_num, num)
						break
					}
				}

			} else {
				previous_num = append(previous_num, num)
			}
			// fmt.Println(num)
			blknum := strconv.Itoa(tmpnumofblk)
			// assign blk
			if datanode_m[num] != nil && len(datanode_m[num]) <= int(maximum_numofblk_pernode) {
				datanode_m[num] = append(datanode_m[num], "blk"+blknum)
			} else {
				// the case that the node's storage is full, then we assign to the other node who has enough space
				if len(datanode_m[num]) > int(maximum_numofblk_pernode) {
					for s := 0; s < len(datanode_m); s++ {
						if len(datanode_m[s]) < int(maximum_numofblk_pernode) {
							datanode_m[s] = append(datanode_m[s], "blk"+blknum)
						}
					}
				} else {
					// initialization
					datanode_m[num] = []string{"blk" + blknum}
				}

			}
			// datanode_m[num] = append
		}
		tmpnumofblk--
	}
	fmt.Println("datanode map: ", datanode_m)
	fmt.Println("length of datanode: ", len(datanode_m))
	return datanode_m
}

func RecordStateInEachSnapshots() {
	files, err := ioutil.ReadDir("../../nodes/node_jsons_reverse")
	if err != nil {
		log.Fatal(err)
	}
	// a string slice to hold the keys
	// k := make([]string, len(files))
	var k []string
	for _, file := range files {
		k = append(k, file.Name()[:10])
		// fmt.Println(file.Name()[:10])
		// fmt.Println(file.Name()[:10])
	}
	// fmt.Println(k)
	//k[i][:10]

	for i := 0; i < len(k); i++ {
		fmt.Println("processing: ", k[i])
		var path string = "../../nodes/node_jsons_reverse/" + k[i] + ".json"
		byteValue, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("something wrong!")
			// return err
		}

		var result map[string]json.RawMessage
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			fmt.Println("something wrong while parsing json!")
			// return err
		}
		states := make(map[string]string)
		for s, _ := range result {
			states[s] = "on"
		}
		//first time just store the state to disk
		if i == 0 {
			globalstate = states
			outputjson(k[i], states)
		} else {
			RecordState(globalstate, states)
			outputjson(k[i], globalstate)
		}
		// k := make(map[string]string)
		// for s, _ := range result {
		// 	k[s] = "on"
		// }

	}

}

func outputjson(filename string, result map[string]string) {
	// Convert golang object back to byte
	byteValue, err := json.Marshal(result)
	if err != nil {
		fmt.Println("something wrong!")
		// return err
	}

	// Write back to file
	err = ioutil.WriteFile("../../nodes/nodes_states/"+filename+"_states.json", byteValue, 0644)
	// fmt.Println(file.Name()[:10])
}

func RecordState(resultjson map[string]string, m2 map[string]string) {
	//copy new map
	m1 := make(map[string]string)
	for k, v := range resultjson {
		m1[k] = v
	}
	// fmt.Println(m1)
	// remove the same nodes and record to the reultjson as "on"
	for k2, _ := range m2 {
		for k1, _ := range m1 {
			result := k1 == k2
			if result {
				if m1[k1] == "off" || m1[k1] == "left" {
					delete(m1, k1)
					delete(m2, k2)
					resultjson[k1] = "join"
				} else {
					delete(m1, k1)
					delete(m2, k2)
					resultjson[k1] = "on"
				}

			}

		}
	}

	// fmt.Println("to on: ", resultjson)
	// the rest of nodes in k2 are the new nodes, recorded as "join"
	for k2, _ := range m2 {
		resultjson[k2] = "join"
	}
	// fmt.Println("to join: ", resultjson)

	//the nodes who has already off
	for k, v := range resultjson {
		if v == "left" {
			resultjson[k] = "off"
		}
	}
	// fmt.Println("to off: ", resultjson)

	// the rest of nodes in k1 are the left nodes, recorded as "left"
	// fmt.Println(m1)
	for k1, _ := range m1 {
		if m1[k1] != "left" && m1[k1] != "off" {
			resultjson[k1] = "left"
		}
	}
	// fmt.Println("to left: ", resultjson)
	fmt.Println()
	globalstate = resultjson

}
func RecordState_test(resultjson map[string]string, m2 map[string]string) {
	//copy new map
	m1 := make(map[string]string)
	for k, v := range resultjson {
		m1[k] = v
	}
	// fmt.Println(m1)
	// remove the same nodes and record to the reultjson as "on"
	for k2, _ := range m2 {
		for k1, _ := range m1 {
			result := k1 == k2
			if result {
				if m1[k1] == "off" || m1[k1] == "left" {
					delete(m1, k1)
					delete(m2, k2)
					resultjson[k1] = "join"
				} else {
					delete(m1, k1)
					delete(m2, k2)
					resultjson[k1] = "on"
				}

			}

		}
	}
	// fmt.Println("m1", m1)
	// fmt.Println("m2", m2)

	// fmt.Println("to on: ", resultjson)
	// the rest of nodes in k2 are the new nodes, recorded as "join"
	for k2, _ := range m2 {
		resultjson[k2] = "join"
	}
	// fmt.Println("to join: ", resultjson)

	//the nodes who has already off
	for k, v := range resultjson {
		if v == "left" {
			resultjson[k] = "off"
		}
	}
	// fmt.Println("to off: ", resultjson)

	// the rest of nodes in k1 are the left nodes, recorded as "left"
	// fmt.Println(m1)
	for k1, _ := range m1 {
		if m1[k1] != "left" && m1[k1] != "off" {
			resultjson[k1] = "left"
		}
	}

	// fmt.Println("to left: ", resultjson)
	fmt.Println()
	globalstate = resultjson

	// fmt.Println("m1: ", m1)
	// fmt.Println("m2: ", m2)
}
func mytest2() {

	m1 := make(map[string]string)
	m1["1.172.126.13:8333"] = "on"
	m1["1.32.249.37:8333"] = "on"
	m1["100.24.150.236:8333"] = "on"

	globalstate = m1

	m2 := make(map[string]string)
	m2["1.32.249.37:8333"] = "on"
	m2["100.24.150.236:8333"] = "on"
	RecordState_test(globalstate, m2)
	fmt.Println(globalstate)

	m3 := make(map[string]string)
	m3["1.32.249.37:8333"] = "on"
	m3["100.24.150.236:8333"] = "on"
	m3["Iamfuckingonian:8333"] = "on"
	RecordState_test(globalstate, m3)
	fmt.Println(globalstate)

	m4 := make(map[string]string)
	m4["1.172.126.13:8333"] = "on"
	m4["1.32.249.37:8333"] = "on"
	m4["100.24.150.236:8333"] = "on"
	m4["Iamfuckingonian:8333"] = "on"
	m4["Iamfuckingonian2:8333"] = "on"
	RecordState_test(globalstate, m4)
	fmt.Println(globalstate)

	m5 := make(map[string]string)
	m5["Iamfuckingonian2:8333"] = "on"
	RecordState_test(globalstate, m5)
	fmt.Println(globalstate)

	m6 := make(map[string]string)
	m6["1.172.126.13:8333"] = "on"
	m6["1.32.249.37:8333"] = "on"
	m6["100.24.150.236:8333"] = "on"
	m6["Iamfuckingonian:8333"] = "on"
	m6["Iamfuckingonian2:8333"] = "on"
	RecordState_test(globalstate, m6)
	fmt.Println(globalstate)
}

func mytest() {
	m1 := make(map[string]string)
	m1["a"] = "on"
	m1["b"] = "on"
	m1["c"] = "on"
	m1["d"] = "on"

	globalstate = m1

	m2 := make(map[string]string)
	m2["a"] = "on"
	m2["b"] = "on"
	m2["c"] = "on"

	RecordState_test(globalstate, m2)
	// fmt.Println(globalstate)
	fmt.Println()

	m3 := make(map[string]string)
	m3["a"] = "on"
	m3["b"] = "on"
	m3["e"] = "on"
	RecordState_test(globalstate, m3)
	fmt.Println()

	m4 := make(map[string]string)
	m4["a"] = "on"
	m4["b"] = "on"
	m4["e"] = "on"
	m4["f"] = "on"
	RecordState_test(globalstate, m4)

	m5 := make(map[string]string)
	m5["a"] = "on"
	m5["b"] = "on"
	m5["c"] = "on"
	m5["e"] = "on"
	m5["f"] = "on"
	RecordState_test(globalstate, m5)

	m6 := make(map[string]string)
	m6["b"] = "on"
	m6["c"] = "on"
	RecordState_test(globalstate, m6)

	m7 := make(map[string]string)
	m7["a"] = "on"
	m7["b"] = "on"
	m7["c"] = "on"
	RecordState_test(globalstate, m7)

	fmt.Println("g: ", globalstate)

	// m := make(map[string]string)
	// m["1.172.126.13:8333"] = "a"
	// m["1.32.249.37:8333"] = "c"
	// m["100.24.150.236:8333"] = "b"
	// fmt.Println(m)
}

// func FindAnotherNum(dontrepeat map[int]int, target int) int {
// 	var ans int
// 	// ans := rand.Intn(100-0) + 0
// 	// ans := rand.Intn(100)

// 	// for {
// 	// 	if dontrepeat[ans] == 0 && ans < target {
// 	// 		break
// 	// 	}
// 	// }

// 	// for i := 0; i < len(dontrepeat); i++ {
// 	// 	if dontrepeat[i] == 0 {
// 	// 		ans = i
// 	// 		break
// 	// 	}
// 	// }

// 	for {
// 		num := rand.Intn(target)
// 		if dontrepeat[num] == 0 {
// 			ans = num
// 			break
// 		}
// 	}

// 	return ans

// }
