package main 
	
import ("fmt" 
		"net/http"
		"time"
		"sync"
		"io"
		ij "steinar/influxjson")

func main() {

	http.HandleFunc("/", influxTest)
	if err := http.ListenAndServe(":8921", nil); err != nil {
		fmt.Println("failed to start server", err)
	}

}

func influxTest(w http.ResponseWriter, r *http.Request) {
	page := `<h1>InfluxDB testing interface</h1></br>
			<form action="/" method="POST">
				<div>
					<label for="connection_settings"><b>Connection settings</b></label><br/>

					<label for="ip">IP:</label><br/>
					<input type="text" name="ip" value="10.161.28.157" size="30"><br />

					<label for="port">Port:</label><br/>
					<input type="text" name="port" value="8086" size="30"><br />
				</div><br/>

				<div>
					<label for="query_settings"><b>Query settings:</b></label><br/>

					<label for="dbname">Database name:</label><br/>
					<input type="text" name="dbname" value="testdb" size="30"><br />

					<label for="query">Query:</label><br/>
					<input type="text" name="query" value="query" size="30"><br />
				</div><br/>

				<div>
					<label for="test_settings"><b>Test settings:</b></label><br/>

					<label for="queries_number">Number of queries:</label><br/>
					<input type="text" name="queries_number" value="1" size="30"><br />

					<label for="threads_number">Number of threads:</label><br/>
					<input type="text" name="threads_number" value="1" size="30"><br />
				</div><br/>

				<div>
					<p><input type="submit" name="do_read" value="Read"></p><p><input type="submit" name="do_write" value="Write"></p>
				</form>
				<p>Duration:</p>`

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "+++STEINAR: Web-error", err)
		} else {
			if message, ok := processRequest(r); ok {
			page += "<p>" + message + "</p>"
		} else if message != "" {
			fmt.Fprintf(w, "<p style='color: red'>+++STEINAR: Web-error!</p>")
		}
		}
		io.WriteString(w, page)
}

func processRequest(r *http.Request) (string, bool) {

	ip, ip_b := r.Form["ip"]
	port, port_b := r.Form["port"]
	dbname, dbname_b := r.Form["dbname"]
	query, query_b := r.Form["query"]
	queries_number, queries_number_b := r.Form["queries_number"]
	threads_number, threads_number_b := r.Form["threads_number"]
	do_read, do_read_b := r.Form["do_read"]
	do_write, do_write_b := r.Form["do_write"]

	fmt.Println(ip, port, dbname, query, queries_number, threads_number, do_read, do_write)
	fmt.Println(ip_b, port_b, dbname_b, query_b, queries_number_b, threads_number_b, do_read_b, do_write_b)


	var returnMessage string
	var returnBool bool 
	switch {
	case (!ip_b):
		returnMessage = "Enter IP!"
		returnBool = true
	case (!port_b):
		returnMessage = "Enter port!"
		returnBool = true
	case (!dbname_b):
		returnMessage = "Enter database name!"
		returnBool = true
	case (!query_b):
		returnMessage = "Enter query!"
		returnBool = true
	case (!queries_number_b):
		returnMessage = "Enter number of queries!"
		returnBool = true
	case (!threads_number_b):
		returnMessage = "Enter number of threads!"
		returnBool = true
	case (do_read_b):
		duration := influxMultiThreadReading(ip[0], port[0], dbname[0], query[0], 1, 1)
		returnMessage = duration.String()
		returnBool = true
	case (do_write_b):
		duration := influxMultiThreadWriting(ip[0], port[0], dbname[0], query[0], 1, 1)
		returnMessage = duration.String()
		returnBool = true
	default:
		returnMessage = ""
		returnBool = true
	}
	
	return returnMessage, returnBool
}

//InfluxDB multithread reading
func influxMultiThreadReading(ip string, port string, dbname string, request string, reqCount int, threadCount int) time.Duration {

	var wg sync.WaitGroup

	start := time.Now()

	for t := 0; t < threadCount; t++ {
		wg.Add(1)
		go influxOneThreadReading(ip, port, dbname, request, reqCount, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	return elapsed

}

func influxOneThreadReading(ip string, port string, dbname string, request string, reqCount int, wg *sync.WaitGroup) {

	for r := 0; r < reqCount; r++ {
		ij.ReadRequest(ip, port, dbname, request)
	}
	defer wg.Done()
}


//InfluxDB multithread writing
func influxMultiThreadWriting(ip string, port string, dbname string, request string, reqCount int, threadCount int) time.Duration {

	var wg sync.WaitGroup

	start := time.Now()

	for t := 0; t < threadCount; t++ {
		wg.Add(1)
		go influxOneThreadWriting(ip, port, dbname, request, reqCount, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	return elapsed

}

func influxOneThreadWriting(ip string, port string, dbname string, request string, reqCount int, wg *sync.WaitGroup) {

	for r := 0; r < reqCount; r++ {
		ij.WriteRequest(ip, port, dbname, request)
	}
	defer wg.Done()
}
