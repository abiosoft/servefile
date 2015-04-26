package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {

	// flags
	var port int
	var filePath string
	var help bool

	flag.IntVar(&port, "p", 8080, "port to listen on")
	flag.StringVar(&filePath, "f", ".", "path to file or directory to serve, defaults to current directory")
	flag.BoolVar(&help, "h", false, "show this help")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if filePath == "." {
		log.Println("-file not specified, serving current directory.")
	}

	// Server
	http.Handle("/", fileHandler(filePath))
	log.Printf("servefile serving on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func fileHandler(filePath string) http.HandlerFunc {
	// path validation
	fd, err := os.Stat(filePath)
	if err != nil {
		log.Println("Error:", err)
		return http.NotFound
	}

	filePath, err = filepath.Abs(filePath)
	if err != nil {
		log.Println("Error:", err)
		return http.NotFound
	}

	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		// scope variables to this function
		fd := fd
		filePath := filePath

		// If servefile is serving a dir and a file is requested
		// from the directory
		if fd.IsDir() && r.URL.Path != "/" {
			filePath = path.Join(filePath, r.URL.Path)
			fd1, err := os.Stat(filePath)
			if err != nil {
				log.Println("Error:", err)
				http.NotFound(w, r)
				return
			}
			fd = fd1
		}

		info := requestInfo{r.RemoteAddr, r.URL.Path, filePath}

		// If a file is about to served
		if !fd.IsDir() {
			// Set file name
			w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%v"`, fd.Name()))
		}

		log.Println(info)
		http.ServeFile(w, r, filePath)
	}
}

type requestInfo struct {
	ip, path, file string
}

func (r requestInfo) String() string {
	return fmt.Sprintf("Request from %v request: %v file: %v", r.ip, r.path, r.file)
}
