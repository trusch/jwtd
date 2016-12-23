package main

import (
	"encoding/json"

	"github.com/trusch/jwtd/storage"
)

func (s *JWTDSuite) TestListGroups() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "GET", "/project/default/group", token, "")
	s.NoError(err)
	groups := []*storage.Group{}
	err = json.Unmarshal([]byte(resp), &groups)
	s.NoError(err)
	s.Equal(3, len(groups))
	s.Equal("jwtd-admin", groups[0].Name)
	s.Equal("http-echo-admin", groups[1].Name)
	s.Equal("http-echo-user", groups[2].Name)
}

func (s *JWTDSuite) TestGetGroup() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "GET", "/project/default/group/jwtd-admin", token, "")
	s.NoError(err)
	group := &storage.Group{}
	err = json.Unmarshal([]byte(resp), group)
	s.NoError(err)
	s.Equal("jwtd-admin", group.Name)
	s.Equal(map[string]map[string]string{"jwtd": map[string]string{"role": "admin"}}, group.Rights)
}

func (s *JWTDSuite) TestDeleteGroup() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	_, err = s.DoRequest("jwtd", "DELETE", "/project/default/group/http-echo-user", token, "")
	s.NoError(err)
	_, err = s.DoRequest("jwtd", "GET", "/project/default/group/http-echo-user", token, "")
	s.Error(err)
	s.Equal("http error: 404", err.Error())
	s.reset()
}

func (s *JWTDSuite) TestCreateGroup() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "POST", "/project/default/group", token, `{"name":"http-echo-tester","rights":{"http-echo-1":{"role":"tester"}}}`)
	s.NoError(err)
	s.Equal("create ok", resp)
	resp, err = s.DoRequest("jwtd", "GET", "/project/default/group/http-echo-tester", token, "")
	s.NoError(err)
	group := &storage.Group{}
	err = json.Unmarshal([]byte(resp), group)
	s.NoError(err)
	s.Equal("http-echo-tester", group.Name)
	s.Equal("tester", group.Rights["http-echo-1"]["role"])
	s.reset()
}

func (s *JWTDSuite) TestUpdateGroup() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "POST", "/project/default/group", token, `{"name":"http-echo-tester","rights":{"http-echo-1":{"role":"tester"}}}`)
	s.NoError(err)
	s.Equal("create ok", resp)
	resp, err = s.DoRequest("jwtd", "PATCH", "/project/default/group/http-echo-tester", token, `{"rights":{"http-echo-1":{"role":"admin"}}}`)
	s.NoError(err)
	s.Equal("update ok", resp)
	resp, err = s.DoRequest("jwtd", "GET", "/project/default/group/http-echo-tester", token, "")
	s.NoError(err)
	group := &storage.Group{}
	err = json.Unmarshal([]byte(resp), group)
	s.NoError(err)
	s.Equal("http-echo-tester", group.Name)
	s.Equal("admin", group.Rights["http-echo-1"]["role"])
	s.reset()
}
