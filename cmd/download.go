package cmd

import (
	"fmt"
	"github.com/ej-agas/himawari-dl/internal"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"sync"

	"github.com/spf13/cobra"
)

var downloadDir *string
var consecutiveDownloads *int

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Himawari satellite images from data.jma.go.jp",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		download()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadDir = downloadCmd.Flags().String("dir", "img", "Directory where to save downloaded images")
	consecutiveDownloads = downloadCmd.Flags().Int("parallel", 5, "Number of consecutive downloads")
}

func download() {
	response, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("error: %s", response.Status)
	}

	root, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	links := make(map[internal.Band][]internal.ImageLink)
	internal.ExtractImgLinks(root, links)

	var wg sync.WaitGroup
	errChan := make(chan error, len(links))
	limiter := make(chan int, *consecutiveDownloads)

	for _, v := range links {
		for _, link := range v {
			wg.Add(1)
			limiter <- 1
			go func(link internal.ImageLink) {
				if err := internal.DownloadImage(link, *downloadDir, &wg); err != nil {
					errChan <- err
				}
				<-limiter
			}(link)

		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		fmt.Println(err)
	}

	fmt.Println("done!")
}
