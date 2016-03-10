package main

import "fmt"
import "net/http"
import "os"
import "bufio"
import "time"
import "sync/atomic"
import "encoding/json"
import "flag"

var usage = `brau
brau brau
`

type Configuration struct {
	Id      string
	Workers int
	URLs    string
}

func loadConfig(fpath string) Configuration {
	config := Configuration{}
	if &fpath != nil {
		configFile, err := os.Open(fpath)
		if err != nil {
			fmt.Println("Error opening config file:", fpath)
			return config
		}
		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("Error parsing config file:", fpath)
		}
	}
	return config
}

const (
	defaultconfig = "config.json"
)

func main() {

	config := loadConfig(defaultconfig)
	//configFile := flag.String("config", "config.json", "Path to JSON Configuration File")
	flag.StringVar(&config.Id, "id", config.Id, "String Identifier")
	flag.IntVar(&config.Workers, "workers", config.Workers, "Number of Parallel Goroutines")
	flag.StringVar(&config.URLs, "urls", config.URLs, "File with list of URLs")

	flag.Usage = func() {
		os.Stderr.WriteString(usage)
		os.Stderr.WriteString("\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Println("Starting Node", config.Id)

	var reqs int64 = 0

	httpreqs := make(chan *http.Request, 100)
	//httpresps := make(chan *http.Response, 100)
	done := make(chan bool)
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	for w := 1; w <= config.Workers; w++ {
		go func(id int) {
			for {
				req, more := <-httpreqs
				if more {
					atomic.AddInt64(&reqs, 1)
					fmt.Println("worker", id, "downloading url", reqs, req.URL)
					_, err := client.Do(req)
					if err != nil {
						// handle error
						fmt.Println("worker", id, "- Error downloading url", reqs, req.URL)
						continue
					}
				} else {
					done <- true
					return
				}
				//httpresps <- resp
			}
		}(w)
	}

	inFile, _ := os.Open(config.URLs)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		req, _ := http.NewRequest("GET", scanner.Text(), nil)
		//req, _ := http.NewRequest("GET", "http://mediavida.com", nil)
		httpreqs <- req
	}
	close(httpreqs)
	<-done
}
