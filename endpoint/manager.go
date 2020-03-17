package endpoint

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"gopkg.in/yaml.v2"
)

// JSON-RPC endpoint structure
type Web3Endpoint struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	ready chan bool
}

// JSON-RPC endpoint manager
type Manager struct {
	Endpoints []Web3Endpoint `yaml:"web3endpoint"`
}

func (m *Manager) Init(configPath string) {
	m.parseConfigFile(configPath)

	for i, _ := range m.Endpoints {
		m.Endpoints[i].ready = make(chan bool, 1024)
		m.Endpoints[i].ready <- true
	}
}

func (m *Manager) ProcessRequest(req *http.Request) ([]byte, error) {
	i := m.GetEndpointIndex(req.RemoteAddr)
	ready := m.Endpoints[i].ready

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	log.Print("######## ProcessRequest")
	log.Print(buf.String())

	<-ready

	reqx, err := http.NewRequest(req.Method, fmt.Sprintf("http://%s:%s", m.Endpoints[i].Host, m.Endpoints[i].Port), buf)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(reqx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ready <- true

	log.Print(string(bytes))
	return bytes, nil
}

func (m *Manager) GetEndpointIndex(remoteHost string) int {
	return rand.Intn(len(m.Endpoints))
}

func (m *Manager) parseConfigFile(configPath string) {
	yamlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlData, m)
	if err != nil {
		log.Fatal(err)
	}
}
