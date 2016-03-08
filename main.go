package main

import "fmt"
import "net/http"
import "os"
import "bufio"

func httpworker(id int, client *http.Client, requests <-chan *http.Request, responses chan<- *http.Response) {
	for req := range requests {
		fmt.Println("worker", id, "downloading url", req.URL)
		resp, err := client.Do(req)
		if err != nil {
			// handle error
			fmt.Println("worker", id, "- Error downloading url", req.URL)
			continue
		}
		responses <- resp
	}
}

func main() {

	httpreqs := make(chan *http.Request, 100)
	httpresps := make(chan *http.Response, 100)

	client := &http.Client{}

	for w := 1; w <= 20; w++ {
		go httpworker(w, client, httpreqs, httpresps)
	}

	inFile, _ := os.Open("alexa.txt")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		req, _ := http.NewRequest("GET", scanner.Text(), nil)
		httpreqs <- req
	}
	close(httpreqs)
	<-httpresps
}
