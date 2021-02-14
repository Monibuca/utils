package utils

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sync/errgroup"
)

// ListenAddrs Listen http and https
func ListenAddrs(addr, addTLS, cert, key string, handler http.Handler) {
	var g errgroup.Group
	if addTLS != "" {
		g.Go(func() error {
			return http.ListenAndServeTLS(addTLS, cert, key, handler)
		})
	}
	if addr != "" {
		g.Go(func() error { return http.ListenAndServe(addr, handler) })
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func CORS(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin[0])
	}
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ReadFileLines(filename string) (lines []string, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	bio := bufio.NewReader(file)
	for {
		var line []byte

		line, _, err = bio.ReadLine()
		if err != nil {
			if err == io.EOF {
				file.Close()
				return lines, nil
			}
			return
		}

		lines = append(lines, string(line))
	}
}

func CurrentDir(path ...string) string {
	_, currentFilePath, _, _ := runtime.Caller(1)
	if len(path) == 0 {
		return filepath.Dir(currentFilePath)
	}
	return filepath.Join(filepath.Dir(currentFilePath), filepath.Join(path...))
}
