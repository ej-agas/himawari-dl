package internal

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const dateTimeFormat = "Mon, 02 Jan 2006 15:04:05 MST"

func ExtractImgLinks(n *html.Node, links map[Band][]ImageLink) map[Band][]ImageLink {
	baseUrl := "https://www.data.jma.go.jp/mscweb/data/himawari/"

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			re := regexp.MustCompile(`r2w_(\w+)_\d+\.jpg`)
			matches := re.FindStringSubmatch(attr.Val)
			if len(matches) < 2 {
				continue
			}
			band := Band(matches[1])
			links[band] = append(links[band], *NewImageLink(band, baseUrl+attr.Val))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = ExtractImgLinks(c, links)
	}

	return links
}

func DownloadImage(link ImageLink, dir string, wg *sync.WaitGroup) error {
	defer wg.Done()

	response, err := http.Get(link.Link)
	if err != nil {
		return fmt.Errorf("error downloading image: %w", err)
	}

	defer response.Body.Close()

	lastModified := response.Header.Get("Last-Modified")

	if lastModified == "" {
		lastModified = time.Now().UTC().Format(dateTimeFormat)
	}

	dateTime, err := time.Parse(dateTimeFormat, lastModified)
	dateTime.Format("Jan")

	if err != nil {
		return fmt.Errorf("error parsing Last-Modified header %s: %w", lastModified, err)
	}

	imgDir := filepath.Join(dir, string(link.Band))

	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating image dir: %w", err)
	}

	seqNum, err := link.SequenceNum()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%d_%s_%d_%s.jpg", dateTime.Day(), strings.ToLower(dateTime.Format("Jan")), dateTime.Year(), seqNum)

	fullFileName := filepath.Join(imgDir, fileName)

	// If file exists, return early
	if _, err := os.Stat(fullFileName); err == nil {
		return nil
	}

	file, err := os.Create(fullFileName)
	if err != nil {
		return fmt.Errorf("error creating file %w", err)
	}

	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("error saving image to file: %w", err)
	}

	fmt.Printf("Saved: %s\n", fileName)

	return nil
}
