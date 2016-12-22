package main

import (
	"encoding/json"

	"github.com/trusch/jwtd/db"
)

func (s *JWTDSuite) TestListUsers() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "GET", "/project/default/user", token, "")
	s.NoError(err)
	users := []*db.User{}
	err = json.Unmarshal([]byte(resp), &users)
	s.NoError(err)
	s.Equal(2, len(users))
	s.Equal("admin", users[0].Name)
	s.Equal("user", users[1].Name)
}

func (s *JWTDSuite) TestGetUser() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "GET", "/project/default/user/admin", token, "")
	s.NoError(err)
	user := &db.User{}
	err = json.Unmarshal([]byte(resp), user)
	s.NoError(err)
	s.Equal("admin", user.Name)
}

func (s *JWTDSuite) TestGetNonExistingUser() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	_, err = s.DoRequest("jwtd", "GET", "/project/default/user/wrong", token, "")
	s.Error(err)
	s.Equal("http error: 404", err.Error())
}

func (s *JWTDSuite) TestCreateUser() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "POST", "/project/default/user", token, `{"username":"test","password":"test","groups":["http-echo-user"]}`)
	s.NoError(err)
	s.Equal("create ok", resp)
	resp, err = s.DoRequest("jwtd", "GET", "/project/default/user/test", token, "")
	s.NoError(err)
	user := &db.User{}
	err = json.Unmarshal([]byte(resp), user)
	s.NoError(err)
	s.Equal("test", user.Name)
	s.reset()
}

func (s *JWTDSuite) TestUpdateUser() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "PATCH", "/project/default/user/user", token, `{"groups":[]}`)
	s.NoError(err)
	s.Equal("update ok", resp)
	resp, err = s.DoRequest("jwtd", "GET", "/project/default/user/user", token, "")
	s.NoError(err)
	user := &db.User{}
	err = json.Unmarshal([]byte(resp), user)
	s.NoError(err)
	s.Equal([]string{}, user.Groups)
	s.reset()
}

func (s *JWTDSuite) TestDeleteUser() {
	token, err := s.GetToken("default", "admin", "admin", "jwtd", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.DoRequest("jwtd", "POST", "/project/default/user", token, `{"username":"test","password":"test","groups":["http-echo-user"]}`)
	s.NoError(err)
	s.Equal("create ok", resp)
	resp, err = s.DoRequest("jwtd", "DELETE", "/project/default/user/test", token, "")
	s.NoError(err)
	s.Equal("delete ok", resp)
	_, err = s.DoRequest("jwtd", "GET", "/project/default/user/test", token, "")
	s.Error(err)
	s.Equal("http error: 404", err.Error())
	s.reset()
}
