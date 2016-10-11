package hls

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/grafov/m3u8"
)

// GetPlaylist fetch content from remote url and return a list of segments
func GetPlaylist(url string) (*m3u8.MediaPlaylist, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, nil
	}

	playlist, listType, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		return nil, err
	}

	switch listType {
	case m3u8.MEDIA:
		p := playlist.(*m3u8.MediaPlaylist)
		return p, nil
	default:
		return nil, nil
	}
}

func BuildSegments(u string) ([]string, error) {
	playlistURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	p, err := GetPlaylist(u)
	if err != nil {
		return nil, err
	}

	var urls []string

	for _, v := range p.Segments {
		if v == nil {
			continue
		}

		var segmentURI string
		if strings.HasPrefix(v.URI, "http") {
			segmentURI, err = url.QueryUnescape(v.URI)
			if err != nil {
				return nil, err
			}
		} else {
			msURL, err := playlistURL.Parse(v.URI)
			if err != nil {
				continue
			}

			segmentURI, err = url.QueryUnescape(msURL.String())
			if err != nil {
				return nil, err
			}
		}
		urls = append(urls, segmentURI)
	}

	return urls, nil
}

func DownloadSegments(u, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	urls, err := BuildSegments(u)
	if err != nil {
		return err
	}
	segmentCount := len(urls)
	ch := make(chan *http.Response, segmentCount)
	errc := make(chan error, 1)

	for _, u := range urls {
		go func(u string) {
			res, e := http.Get(u)
			if e != nil {
				errc <- e
				return
			}
			//defer res.Body.Close()
			ch <- res
		}(u)
	}

	segmentProcessed := 0
	for {
		select {
		case r := <-ch:
			segmentProcessed++
			// need to append to the outout file
			_, err = io.Copy(out, r.Body)
			if err != nil {
				errc <- err
			}
			if segmentProcessed == segmentCount {
				return nil
			}
		case err := <-errc:
			if err != nil {
				return err
			}
		}
	}
}

// Download hls segments into a single output file based on the remote index
func Download(u, output string) error {
	err := DownloadSegments(u, output)
	if err != nil {
		return err
	}

	return nil
}
