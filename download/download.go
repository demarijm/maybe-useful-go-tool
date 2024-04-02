package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func DownloadData(pages int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	combinedData := make([]json.RawMessage, 0)

	for i := 0; i < pages; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			url := fmt.Sprintf("turn this into an env variable %d", i)
			// request the dynamic url
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal("There was an error trying to get this")
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("Receieved non 200 code %d", resp.StatusCode)
				resp.Body.Close()
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failed to read body", err)
			}

			mu.Lock()
			combinedData = append(combinedData, body)
			mu.Unlock()

		}(i)
	}
	wg.Wait()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, data := range combinedData {
		buffer.Write(data)
		if i < len(combinedData)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")

	fileName := "combined_data.json"
	if err := os.WriteFile(fileName, buffer.Bytes(), 0666); err != nil {
		log.Fatalf("Faileds to write combined data %s", err)
	}

	log.Println("All data has been combined", fileName)
}
