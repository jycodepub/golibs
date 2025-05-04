package net

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type endpoint struct {
	Method        string   `json:"method"`
	Url           string   `json:"url"`
	ContentType   string   `json:"contentType"`
	ResponseFiles []string `json:"responseFiles"`
	Responses     []string `json:"responses"`
}

func (e endpoint) String() string {
	return fmt.Sprintf("{ method: %v, url: %v, path: %v }\n", e.Method, e.Url, e.ResponseFiles)
}

func (e endpoint) isValid() bool {
	if e.Method == "" {
		return false
	}
	if e.Url == "" {
		return false
	}
	if e.Responses == nil && e.ResponseFiles == nil {
		return false
	}
	return true
}

var config = make(map[string]endpoint)
var responseDir string

type Server struct {
	addr    string
	handler http.HandlerFunc
}

func NewMockServer(addr string, configPath string, respdir string) *Server {
	responseDir = respdir
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var endpoints []endpoint
	err = json.Unmarshal(b, &endpoints)
	if err != nil {
		log.Fatal(err)
	}

	for _, ep := range endpoints {
		if !ep.isValid() {
			log.Fatal("Invalid config")
		}
		config[ep.Method+ep.Url] = ep
		log.Println("Endpoint:", ep)
	}

	log.Println("Starting server at", addr)
	return &Server{
		addr:    addr,
		handler: serve,
	}
}

func (s *Server) Start() {
	ss := &http.Server{
		Addr:           s.addr,
		Handler:        s.handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(ss.ListenAndServe())
}

func serve(w http.ResponseWriter, req *http.Request) {
	key := req.Method + req.URL.Path
	if req.URL.RawQuery != "" {
		key += ("?" + req.URL.RawQuery)
	}
	endpoint, ok := config[key]
	if ok {
		var err error
		var response string
		if endpoint.Responses == nil {
			response, err = loadPayload(endpoint.ResponseFiles[randomIndex(len(endpoint.ResponseFiles))])
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			response = endpoint.Responses[randomIndex(len(endpoint.Responses))]
		}
		w.Header().Add("Content-Type", endpoint.ContentType)
		fmt.Fprint(w, response)
		log.Println("Served request:", endpoint.Method, endpoint.Url, endpoint.ContentType)
	} else {
		http.NotFound(w, req)
	}
}

func loadPayload(path string) (string, error) {
	var apath string
	if !filepath.IsAbs(path) {
		apath = filepath.Join(responseDir, path)
	} else {
		apath = path
	}
	f, err := os.Open(apath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func randomIndex(maxIndex int) int {
	if maxIndex == 1 {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	return r.Intn(maxIndex)
}
