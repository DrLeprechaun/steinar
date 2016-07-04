
package influxjson 
	
import ("fmt"; "net/http"; "io/ioutil"; "strings"; "bytes")

//JSON STRUCTURES
type InfluxRequestResult struct {
	series []Serie
}

type Serie struct {
	name string
	columns []string
	valuesBlocks []ValuesBlock
}

type ValuesBlock struct {
	values []string
}

//InfluxDB response parser
func InfluxJsonParser(jsonString string) InfluxRequestResult {

	var influxRequestResult InfluxRequestResult

	results := strings.Split(jsonString, "results")[1:]

	for _, r := range results {
		seriesRaw := strings.Split(r, "series")

			series := strings.Split(seriesRaw[1], "name")[1:]

			var sers []Serie
			for _, s := range series {

				var serie Serie

				//Serie name extraction
				nameRaw := strings.Split(s, `:"`)
				nameRaw = strings.Split(nameRaw[1], `[`)
				nameRaw = strings.Split(nameRaw[0], `,`)
				name := nameRaw[0]
				name = strings.Trim(name, `"`)
				serie.name = name
				
				//Column names extraction
				columnsRaw := strings.Split(s, "columns")
				columnsRaw = strings.Split(columnsRaw[1], "[")
				columnsRaw = strings.Split(columnsRaw[1], "]")
				columns := strings.Split(columnsRaw[0], ",")

				var cols []string
				for c := 0; c < len(columns); c++ {
					cols = append(cols, strings.Trim(columns[c], `"`))
				}

				//Values extraction
				valuesBlocksRaw := strings.Split(s, "values")[1:2]
				valuesBlocksRaw = strings.Split(valuesBlocksRaw[0], ":[")[1:2]
				valuesBlocksRaw = strings.Split(valuesBlocksRaw[0], "]}")[0:1]
				valuesBlocks := strings.Split(valuesBlocksRaw[0], "],")

				var valsBlocks []ValuesBlock

				for _, vb := range valuesBlocks {

					var valBlock ValuesBlock

					valuesRaw := strings.Split(vb, "[")[1:2]
					values := strings.Split(valuesRaw[0], ",")

					var vals []string
					for _, v := range values {

						v = strings.Trim(v, "]")
						v = strings.Trim(v, `"`)

						vals = append(vals, v)
					}

					valBlock.values = vals
					valsBlocks = append(valsBlocks, valBlock)
				}
				
				serie.columns = cols
				serie.valuesBlocks = valsBlocks
				sers = append(sers, serie)
			}

		influxRequestResult.series = sers
	}

	return influxRequestResult
}


//Influx write request
func WriteRequest(ip string, port string, dbname string, requestBody string) {
	bodyBytes := bytes.NewBufferString(requestBody)
	resp, err := http.Post("http://" + ip + ":" + port + "/write?db=" + dbname, "application/x-www-form-urlencoded", bodyBytes)
	if err != nil {
		fmt.Println("+++STEINAR: Error while excecuting request! Check your connection settings.")
	}
	defer resp.Body.Close()

	bodyArray, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyArray[:len(bodyArray)])

	if (strings.Contains(bodyString, "error")) {
		fmt.Println("+++STEINAR: Bad request! Check request synthax.")
	}
	if (bodyString == "") {
		fmt.Println("+++STEINAR: Data inserted successfuly.")
	}
}

//Influx read request
func ReadRequest(ip string, port string, dbname string, request string) string {
	url := "http://" + ip + ":" + port + "/query?db=" + dbname + "&q=" + request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("+++STEINAR: Error while excecuting request! Check your connection settings.")
	}
	defer resp.Body.Close()
	bodyArray, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyArray[:len(bodyArray)])

	if (strings.Contains(bodyString, "error")) {
		fmt.Println("+++STEINAR: Bad request! Check request synthax.")
	}

	return bodyString
}