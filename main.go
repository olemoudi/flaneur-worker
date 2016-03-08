package main

import "fmt"
import "net/http"
import "os"
import "bufio"
import "time"
import "sync/atomic"

func main() {

	var reqs int64 = 0

	httpreqs := make(chan *http.Request, 100)
	//httpresps := make(chan *http.Response, 100)
	done := make(chan bool)
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	for w := 1; w <= 20; w++ {
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

	inFile, _ := os.Open(os.Args[1])
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// req, _ := http.NewRequest("GET", scanner.Text(), nil)
		req, _ := http.NewRequest("GET", "http://mediavida.com", nil)
		httpreqs <- req
	}
	close(httpreqs)
	<-done
}
