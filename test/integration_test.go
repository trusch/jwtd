package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JWTDSuite struct {
	suite.Suite
}

type Label struct {
	Key   string
	Value string
}

func (s *JWTDSuite) SetupSuite() {
	script := `
  cp jwtd.yaml.tmpl default.yaml

  docker stop jwtd
  docker rm jwtd
  docker run --name jwtd -d \
    -v $(pwd)/pki/jwtd.key:/etc/jwtd/jwtd.key \
    -v $(pwd)/default.yaml:/etc/jwtd/default.yaml \
    trusch/jwtd

  docker stop http-echo-1
  docker rm http-echo-1
  docker run --name http-echo-1 -d trusch/http-echo

  docker stop http-echo-2
  docker rm http-echo-2
  docker run --name http-echo-2 -d trusch/http-echo

  docker stop jwtd-proxy
  docker rm jwtd-proxy
  docker run --name jwtd-proxy -d \
    -v $(pwd)/pki/jwtd.crt:/etc/jwtd-proxy/jwtd.crt \
    -v $(pwd)/pki/jwtd.key:/etc/jwtd-proxy/jwtd.key \
    -v $(pwd)/pki/http-echo-1.crt:/etc/jwtd-proxy/http-echo-1.crt \
    -v $(pwd)/pki/http-echo-1.key:/etc/jwtd-proxy/http-echo-1.key \
    -v $(pwd)/pki/http-echo-2.crt:/etc/jwtd-proxy/http-echo-2.crt \
    -v $(pwd)/pki/http-echo-2.key:/etc/jwtd-proxy/http-echo-2.key \
    -v $(pwd)/jwtd-proxy.yaml:/etc/jwtd-proxy/config.yaml \
    --link http-echo-1 \
    --link http-echo-2 \
    --link jwtd \
    -p 443:443 \
    trusch/jwtd-proxy
  `
	cmd := exec.Command("/bin/bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Print("prepare test... ")
	err := cmd.Run()
	fmt.Println("done.")
	s.NoError(err)

	out := &bytes.Buffer{}
	cmd = exec.Command("docker", "inspect", "--format", "{{ .State.Running }}", "jwtd")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	s.NoError(err)
	s.Equal("true\n", string(out.Bytes()), "jwtd not running after setup!")

	out = &bytes.Buffer{}
	cmd = exec.Command("docker", "inspect", "--format", "{{ .State.Running }}", "jwtd-proxy")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	s.NoError(err)
	s.Equal("true\n", string(out.Bytes()), "jwtd-proxy not running after setup!")

	out = &bytes.Buffer{}
	cmd = exec.Command("docker", "inspect", "--format", "{{ .State.Running }}", "http-echo-1")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	s.NoError(err)
	s.Equal("true\n", string(out.Bytes()), "http-echo-1 not running after setup!")

	out = &bytes.Buffer{}
	cmd = exec.Command("docker", "inspect", "--format", "{{ .State.Running }}", "http-echo-2")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	s.NoError(err)
	s.Equal("true\n", string(out.Bytes()), "http-echo-2 not running after setup!")
}

func (s *JWTDSuite) reset() {
	script := `
	docker stop jwtd
  docker rm jwtd
	cp jwtd.yaml.tmpl default.yaml
  docker run --name jwtd -d \
    -v $(pwd)/pki/jwtd.key:/etc/jwtd/jwtd.key \
    -v $(pwd)/default.yaml:/etc/jwtd/default.yaml \
    trusch/jwtd
  `
	cmd := exec.Command("/bin/bash", "-c", script)
	cmd.Run()
}

func (s *JWTDSuite) GetToken(project, username, password, service string, labels []*Label) (string, error) {
	labelStr := ""
	for idx, label := range labels {
		if idx != 0 {
			labelStr += ","
		}
		labelStr += fmt.Sprintf(`"%v":"%v"`, label.Key, label.Value)
	}
	payload := fmt.Sprintf(`{"project":"%v","username":"%v","password":"%v","service":"%v","labels":{%v}}`, project, username, password, service, labelStr)
	token, err := s.DoRequest("jwtd", "POST", "/token", "", payload)
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("no token returned")
	}
	return token, nil
}

type HttpEchoResponse struct {
	URL    string
	Header http.Header
	Body   string
}

func (s *JWTDSuite) EchoOneRequest(method, url, token, body string) (*HttpEchoResponse, error) {
	data, err := s.DoRequest("http-echo-1", method, url, token, body)
	if err != nil {
		return nil, err
	}
	resp := &HttpEchoResponse{}
	err = json.Unmarshal([]byte(data), resp)
	return resp, err
}

func (s *JWTDSuite) EchoTwoRequest(method, url, token, body string) (*HttpEchoResponse, error) {
	data, err := s.DoRequest("http-echo-2", method, url, token, body)
	if err != nil {
		return nil, err
	}
	resp := &HttpEchoResponse{}
	err = json.Unmarshal([]byte(data), resp)
	return resp, err
}

func (s *JWTDSuite) DoRequest(host, method, url, token, body string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, "https://localhost"+url, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Host = host
	if token != "" {
		req.Header.Set("Authorization", "bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return string(bs), fmt.Errorf("http error: %v", resp.StatusCode)
	}
	return string(bs), nil
}

func TestJWTDSuite(t *testing.T) {
	suite.Run(t, new(JWTDSuite))
}
