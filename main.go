package main

import (
	"log"
	"net/http"
	"strconv"
)

const Port = 80

var Version = "0.0"
var VersionTime = "0000-00-00 00:00"

func download(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.ParseInt(r.URL.RawQuery, 10, 64)
	if err != nil || length <= 0 {
		length = 204800
	}

	size := strconv.FormatInt(length, 10)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=dntest"+size)
	w.Header().Set("Content-Length", size)
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method == http.MethodHead {
		return
	}

	var buf [1024]byte
	for i := range buf {
		buf[i] = '0'
	}

	for remaining := length; remaining > 0; {
		chunk := min(remaining, int64(len(buf)))
		n, err := w.Write(buf[:chunk])
		if err != nil {
			log.Printf("writing download: %v", err)
			return
		}
		if int64(n) != chunk {
			log.Printf("writing download: short write (%d of %d bytes)", n, chunk)
			return
		}
		remaining -= int64(n)
	}
}

func main() {
	http.HandleFunc("/dntest.php", download)
	log.Printf("Download tests service v%s (%s)\n", Version, VersionTime)
	log.Printf("Listening on :%d", Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Port), nil))
}
