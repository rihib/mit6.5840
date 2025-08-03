// https://go.dev/tour/concurrency/10
package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	var wg sync.WaitGroup
	visited := make(map[string]struct{})
	var mu sync.Mutex

	var crawl func(string, int)
	crawl = func(url string, depth int) {
		defer wg.Done()
		if depth <= 0 {
			return
		}
		mu.Lock()
		if _, ok := visited[url]; ok {
			mu.Unlock()
			return
		}
		visited[url] = struct{}{}
		mu.Unlock()

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			wg.Add(1)
			go crawl(u, depth-1)
		}
	}

	wg.Add(1)
	go crawl(url, depth)
	wg.Wait()
}

func main() {
	fetcher := newFakeFetcher()
	Crawl("https://golang.org/", 4, fetcher)
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type fakeFetcher struct {
	data    fakeData
	visited *fakeVisitedURLs
}

type fakeData map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

type fakeVisitedURLs struct {
	mu   sync.Mutex
	urls map[string]struct{}
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if _, ok := f.visited.urls[url]; ok {
		return "", nil, fmt.Errorf("already visited: %s", url)
	}
	f.visited.mu.Lock()
	f.visited.urls[url] = struct{}{}
	f.visited.mu.Unlock()

	if res, ok := f.data[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

func newFakeFetcher() *fakeFetcher {
	f := new(fakeFetcher)
	f.data = fakeData{
		"https://golang.org/": &fakeResult{
			"The Go Programming Language",
			[]string{
				"https://golang.org/pkg/",
				"https://golang.org/cmd/",
			},
		},
		"https://golang.org/pkg/": &fakeResult{
			"Packages",
			[]string{
				"https://golang.org/",
				"https://golang.org/cmd/",
				"https://golang.org/pkg/fmt/",
				"https://golang.org/pkg/os/",
			},
		},
		"https://golang.org/pkg/fmt/": &fakeResult{
			"Package fmt",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			},
		},
		"https://golang.org/pkg/os/": &fakeResult{
			"Package os",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			},
		},
	}
	f.visited = newFakeVisitedURLs()
	return f
}

func newFakeVisitedURLs() *fakeVisitedURLs {
	visited := new(fakeVisitedURLs)
	visited.urls = make(map[string]struct{})
	return visited
}
