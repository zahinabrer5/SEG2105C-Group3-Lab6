package main

import (
    "fmt"
    "io"
    "net/http"
    "sync"
)

// Structure to store results
type FetchResult struct {
    URL        string
    StatusCode int
    Size       int
    Error      error
}

// Worker function
func worker(wg *sync.WaitGroup, id int, jobs <-chan string, results chan<- FetchResult) {
    defer wg.Done() // wg.Add(-1)

    for url := range jobs {
        // fetch the URL
        resp, err := http.Get(url)
        if err != nil {
            results <- FetchResult{URL: url, StatusCode: 0, Size: 0, Error: err}
            continue
        }
        body, err := io.ReadAll(resp.Body)
        resp.Body.Close()

        // send Result struct to results channel
        size := 0
        statusCode := 0
        if err == nil {
            size = len(body)
            statusCode = resp.StatusCode
        }
        results <- FetchResult{URL: url, StatusCode: statusCode, Size: size, Error: err}
    }
}

func main() {
    var wg sync.WaitGroup

    urls := []string{
        "https://example.com",
        "https://golang.org",
        "https://uottawa.ca",
        "https://github.com",
        "https://httpbin.org/get",
    }

    numWorkers := 5

    jobs := make(chan string, len(urls))
    results := make(chan FetchResult, len(urls))

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(&wg, w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= len(urls); j++ {
        jobs <- urls[j-1]
    }
    close(jobs)

    wg.Wait()
    // Collect results
    fmt.Println("Fetching URLs concurrently using worker pool...")
    fmt.Println()
    for i := 1; i <= len(urls); i++ {
        result := <- results
        fmt.Printf("%s | Status: %d | Size: %d bytes\n", result.URL, result.StatusCode, result.Size)
    }

    fmt.Println("\nScraping complete!")
}
