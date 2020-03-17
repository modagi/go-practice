package gatewayserver

import (
	"io/ioutil"
	"log"
	"net/http"

	"../endpoint"

	"gopkg.in/yaml.v2"
)

var manager endpoint.Manager

// gateway server
type Server struct {
	Config struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

// http request handler
type reqHandler struct {
	http.Handler
}

func (h *reqHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res, err := manager.ProcessRequest(req)
	if err != nil {
		log.Print(err)
		return
	}
	w.Write(res)
}

func (s *Server) Init(configPath string) {
	s.parseConfigFile(configPath)
	manager.Init(configPath)

	http.Handle("/", new(reqHandler))
}

func (s *Server) parseConfigFile(configPath string) {
	yamlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlData, s)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) RunServer() {
	log.Print("Run gateway server")
	http.ListenAndServe(":"+s.Config.Port, nil)
}
