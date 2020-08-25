package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	//%-----------collect data from bitnodes' Apis--------------------%
	// download from Bitnodes  ex: https://bitnodes.io/api/v1/snapshots/?page=1000 
	GetNodeSnapshots()
	//plot the pattern
	plotsnapshots()

	// Get https://bitnodes.io/api/v1/snapshots/<timestamp>
	// if the request limit is reached, then print "reuquest got throttled and stop the program."
	timestamplist := GetTimestamps()
	for i := 0; i < len(timestamplist); i++ {
		fmt.Println("downlaoding timestamp: ", timestamplist[i])
		if !GetSnapshotsWithTimestamps(timestamplist[i]) {
			fmt.Println("request got throttled")
			break
		}
	}
	
	// %------------------------------Run------------------------------%
	// Find who is always up in the network in this duration.
	WhoIsAlwaysUp()	
	// Calculate node states in each timestamp.
	CalculateEachSession()

	// 500 prover nodes
	// --------get 500/5G results--------
	RunProcess("Baseline_norepair_",500, false, 5, 1,"5Baseline.csv")
	// RunProcess("R2_",500, false, 5, 2,"5R2.csv")
	// RunProcess("R4_",500, false, 5, 4,"5R4.csv")
	// RunProcess("R8_", 500, false, 5, 8, "5R8.csv")
	// RunProcess("R16_", 500, false, 5, 16, "5R16.csv")
	// RunProcess("R32_", 500, false, 5, 32, "5R32.csv")

	// --------get 500/10G results--------
	// RunProcess("Baseline_norepair_",500, false, 10, 1,"10Baseline.csv")
	// RunProcess("R2_",500, false, 10, 2,"10R2.csv")
	// RunProcess("R4_",500, false, 10, 4,"10R4.csv")
	// RunProcess("R8_", 500, false, 10, 8, "10R8.csv")
	// RunProcess("R16_", 500, false, 10, 16, "5R16.csv")
	// RunProcess("R32_", 500, false, 10, 32, "5R32.csv")

	// --------get 500/20G results--------
	// RunProcess("Baseline_norepair_",500, false, 20, 1,"20Baseline.csv")
	// RunProcess("R2_",500, false, 20, 2,"20R2.csv")
	// RunProcess("R4_",500, false, 20, 4,"20R4.csv")
	// RunProcess("R8_", 500, false, 20, 8, "20R8.csv")
	// RunProcess("R16_", 500, false, 20, 16, "20R16.csv")
	// RunProcess("R32_", 500, false, 20, 32, "20R32.csv")

	// --------get 500/40G results--------
	// RunProcess("Baseline_norepair_",500, false, 40, 1,"40Baseline.csv")
	// RunProcess("R2_",500, false, 40, 2,"40R2.csv")
	// RunProcess("R4_",500, false, 40, 4,"40R4.csv")
	// RunProcess("R8_", 500, false, 40, 8, "40R8.csv")
	// RunProcess("R16_", 500, false, 40, 16, "40R16.csv")
	// RunProcess("R32_", 500, false, 40, 32, "40R32.csv")

	// 1000 prover nodes
	// --------get 1000/5G results--------
	// RunProcess("Baseline_norepair_",1000, false, 5, 1,"1000-5Baseline.csv")
	// RunProcess("R2_",1000, false, 5, 2,"1000-5R2.csv")
	// RunProcess("R4_",1000, false, 5, 4,"1000-5R4.csv")
	// RunProcess("R8_", 1000, false, 5, 8, "1000-5R8.csv")
	// RunProcess("R16_", 1000, false, 5, 16, "1000-5R16.csv")
	// RunProcess("R32_", 1000, false, 5, 32, "1000-5R32.csv")

	// --------get 1000/10G results--------
	// RunProcess("Baseline_norepair_",1000, false, 10, 1,"1000-10Baseline.csv")
	// RunProcess("R2_",1000, false, 10, 2,"1000-10R2.csv")
	// RunProcess("R4_",1000, false, 10, 4,"1000-10R4.csv")
	// RunProcess("R8_", 1000, false, 10, 8, "1000-10R8.csv")
	// RunProcess("R16_", 1000, false, 10, 16, "1000-10R16.csv")
	// RunProcess("R32_", 1000, false, 10, 32, "1000-10R32.csv")

	// --------get 1000/20G results--------
	// RunProcess("Baseline_norepair_",1000, false, 20, 1,"1000-20Baseline.csv")
	// RunProcess("R2_",1000, false, 20, 2,"1000-20R2.csv")
	// RunProcess("R4_",1000, false, 20, 4,"1000-20R4.csv")
	// RunProcess("R8_", 1000, false, 20, 8, "1000-20R8.csv")
	// RunProcess("R16_", 1000, false, 20, 16, "1000-20R16.csv")
	// RunProcess("R32_", 1000, false, 20, 32, "1000-20R32.csv")

	// --------get 1000/40G results--------
	// RunProcess("Baseline_norepair_",1000, false, 40, 1,"1000-40Baseline.csv")
	// RunProcess("R2_",1000, false, 40, 2,"1000-40R2.csv")
	// RunProcess("R4_",1000, false, 40, 4,"1000-40R4.csv")
	// RunProcess("R8_", 1000, false, 40, 8, "1000-40R8.csv")
	// RunProcess("R16_", 1000, false, 40, 16, "1000-5R16.csv")
	// RunProcess("R32_", 1000, false, 40, 32, "1000-5R32.csv")

	// 2000 prover nodes
	// --------get 2000/5G results--------
	// RunProcess("Baseline_norepair_",2000, false, 5, 1,"2000-5Baseline.csv")
	// RunProcess("R2_",2000, false, 5, 2,"2000-5R2.csv")
	// RunProcess("R4_",2000, false, 5, 4,"2000-5R4.csv")
	// RunProcess("R8_", 2000, false, 5, 8, "2000-5R8.csv")
	// RunProcess("R16_", 2000, false, 5, 16, "2000-5R16.csv")
	// RunProcess("R32_", 2000, false, 5, 32, "2000-5R32.csv")

	// --------get 2000/10G results--------
	// RunProcess("Baseline_norepair_",2000, false, 10, 1,"2000-10Baseline.csv")
	// RunProcess("R2_",2000, false, 10, 2,"2000-10R2.csv")
	// RunProcess("R4_",2000, false, 10, 4,"2000-10R4.csv")
	// RunProcess("R8_", 2000, false, 10, 8, "2000-10R8.csv")
	// RunProcess("R16_", 2000, false, 10, 16, "2000-10R16.csv")
	// RunProcess("R32_", 2000, false, 10, 32, "2000-10R32.csv")

	// --------get 2000/20G results--------
	// RunProcess("Baseline_norepair_",2000, false, 20, 1,"2000-20Baseline.csv")
	// RunProcess("R2_",2000, false, 20, 2,"2000-20R2.csv")
	// RunProcess("R4_",2000, false, 20, 4,"2000-20R4.csv")
	// RunProcess("R8_", 2000, false, 20, 8, "2000-20R8.csv")
	// RunProcess("R16_", 2000, false, 20, 16, "2000-20R16.csv")
	// RunProcess("R32_", 2000, false, 20, 32, "2000-20R32.csv")

	// --------get 2000/40G results--------
	// RunProcess("Baseline_norepair_",2000, false, 40, 1,"2000-40Baseline.csv")
	// RunProcess("R2_",2000, false, 40, 2,"2000-40R2.csv")
	// RunProcess("R4_",2000, false, 40, 4,"2000-40R4.csv")
	// RunProcess("R8_", 2000, false, 40, 8, "2000-40R8.csv")
	// RunProcess("R16_", 2000, false, 40, 16, "2000-40R16.csv")
	// RunProcess("R32_", 2000, false, 40, 32, "2000-40R32.csv")
	
	// %------------------------------Others------------------------------%
	// calculateChurn_Paper()
	// CalculateEachSession()

	// ReadIP_States()
	// Calculate_totalChurn()
	// Calculate_averageChurn_Simple()
}

// the necessary function for using random APIs
func init() {
	rand.Seed(time.Now().UnixNano())
}

func test() {
	// var total_nodes int = 500
	min := 10
	max := 30
	fmt.Println(rand.Intn(max-min) + min)

	var numofBlk float64 = 3000
	replication_factor := []float64{1, 2, 4, 8}
	var GBslice []float64
	for i := 0; i < len(replication_factor); i++ {
		GBslice = append(GBslice, replication_factor[i]*numofBlk*0.128)
	}
	fmt.Println(GBslice)
	// d1 := [4]float64{265,529,1052,2116}
	d2 := [4]string{"Baseline", "R2", "R4", "R8"}
	s1 := [7]float64{2, 4, 8, 10, 16, 32, 64}

	for i := 0; i < len(GBslice); i++ {
		fmt.Println("Replicaiton factor: ", d2[i])
		for j := 0; j < len(s1); j++ {
			ans := GBslice[i] / s1[j]
			// fmt.Printf("storage limit: %f, nodes: %f \n",s1[j], ans)
			fmt.Printf("storage limit: %f, nodes: ", s1[j])
			fmt.Println(math.Ceil(ans))
		}
	}

}

func readcsv_reverse() {
	csvfile, err := os.Open("nodes_snapshots.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))
	var slice1, slice2 []string

	//my code

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
		slice1 = append(slice1, record[0])
		slice2 = append(slice2, record[1])
	}

	for i := len(slice1) - 1; i <= len(slice1); i-- {
		appendToCSV_pure(slice1[i], slice2[i], "nodes_snapshots_reverse_forchurn.csv")
	}

}

// func appendToJsonl(nodesjson map[string]json.RawMessage, totalNodes string, filename string) {
// 	outputmap := make(map[string]map[string]string)
// 	// output.total_nodes = totalNodes

// 	for k, v := range nodesjson {
// 		objectmap := make(map[string]string)
// 		//The cool thing is that here I set int slice, it only read the int data into this slice
// 		var temp []int
// 		var tmp []string
// 		json.Unmarshal(v, &temp)
// 		json.Unmarshal(v, &tmp)
// 		// fmt.Println(temp)
// 		// fmt.Println(tmp)
// 		objectmap["start_time"] = strconv.Itoa(temp[2]) // result: s = "-18"temp[2]
// 		objectmap["Timezone"] = tmp[10]
// 		objectmap["ASN"] = tmp[11]
// 		objectmap["Organization"] = tmp[12]
// 		// tmpstr = tmpstr[:len(tmpstr)-1]
// 		// objectmap["info"] = tmpstr
// 		outputmap[k] = objectmap
// 		// fmt.Println(objectmap)
// 	}
// 	// fmt.Println(outputmap)
// 	// fmt.Println(output.total_nodes)
// 	file, err := json.Marshal(outputmap)

// 	if err != nil {
// 		fmt.Println("something wrong while writing json!")
// 	}
// 	var path string = "./node_jsons/" + filename + ".json"
// 	// Write to file
// 	_ = ioutil.WriteFile(path, file, 0644)
// }
// func getfromTimestamp() {
// 	url := "https://bitnodes.io/api/v1/snapshots/1593319209/"
// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	//request body, which is byte[]
// 	sitemap, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// a map container to decode the JSON structure into
// 	result := make(map[string]json.RawMessage)

// 	// unmarschal JSON
// 	err = json.Unmarshal(sitemap, &result)
// 	if err != nil {
// 		return
// 	}
// 	// fmt.Println(result)
// 	total_nodes := result["total_nodes"]
// 	nodes := result["nodes"]

// 	nodesjson := make(map[string]json.RawMessage)
// 	err = json.Unmarshal(nodes, &nodesjson)
// 	appendToJsonl(nodesjson, string(total_nodes), "1593319209")

// 	//pretty json
// 	// var prettyJSON bytes.Buffer
// 	// error := json.Indent(&prettyJSON, nodes, "", "\t")
// 	// if error != nil {
// 	// 	fmt.Println("JSON parse error: ", error)
// 	// 	return
// 	// }
// 	// fmt.Println("Pretty Json:", string(prettyJSON.Bytes()))

// }

// type UpLoadSomething struct {
// 	Type   string
// 	Object interface{}
// }

// type File struct {
// 	FileName string
// }

// type Png struct {
// 	Wide  int
// 	Hight int
// }
// func test() {

// 	input := `
//     {
//         "type": "File",
//         "object": {
//             "filename": "for test"
// 		}
//     }
//     `
// 	var object json.RawMessage
// 	ss := UpLoadSomething{
// 		Object: &object,
// 	}
// 	if err := json.Unmarshal([]byte(input), &ss); err != nil {
// 		panic(err)
// 	}
// 	switch ss.Type {
// 	case "File":
// 		var f File
// 		if err := json.Unmarshal(object, &f); err != nil {
// 			panic(err)
// 		}
// 		println(f.FileName)
// 	case "Png":
// 		var p Png
// 		if err := json.Unmarshal(object, &p); err != nil {
// 			panic(err)
// 		}
// 		println(p.Wide)
// 	}
// }

// func main() {
// 	// fmt.Println(base)
// 	GetNodeSnapshots()

// }

// The testing of process.go
// mytest()
// mytest2()
// GetFilesName()

// reverse the node_snapshots.csv, let it be the acsending order
// readcsv_reverse()

// calculate churn
// calculateChurn()
// calculateDailyChurnCumulative()

// calculate churn in cumulative form
// calculateChurnCumulative()
// caluculateCount()

//simulation test
// simulate()

// RecordStateInEachSnapshots()

// find who is always up in the snapshots
// WhoIsAlwaysUp()

// --------get Baseline result--------
// assignblkToFirst("Baseline_norepair_",500)
// assignblkToFirst_withEmpty("Baseline_norepair_")
// // record the lost and repair in each timestamp
// calculateLostAndRepair("Baseline_norepair_")

// --------get R2 result--------
// assignblkToFirst("R2_")
// assignblkToFirst_withEmpty("R2_")
// // record the lost and repair in each timestamp
// calculateLostAndRepair("R2_")
// RunProcess("R2_",500, true)

// --------get R4 result--------
// assignblkToFirst("R4_")
// assignblkToFirst_withEmpty("R4_")
// // record the lost and repair in each timestamp
// calculateLostAndRepair("R4_")
// RunProcess("R4_",500, true)

// --------get R8 result--------
// assignblkToFirst("R8_")
// assignblkToFirst_withEmpty("R8_")
// // record the lost and repair in each timestamp
// calculateLostAndRepair("R8_")
// RunProcess("R8_",500, true)

// %------------------------------Analysis------------------------------%
// RunAnalysis()
