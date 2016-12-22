package main

func (s *JWTDSuite) TestRequestEchoAdminPageCorrect() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-1", "role", "admin")
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.EchoOneRequest("GET", "/admin", token, "")
	s.NoError(err)
	s.Equal("/admin", resp.URL)
}

func (s *JWTDSuite) TestRequestEchoAdminPageWrongRole() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-1", "role", "user")
	s.NoError(err)
	s.NotEmpty(token)
	_, err = s.EchoOneRequest("GET", "/admin", token, "")
	s.Error(err)
	s.Equal("http error: 401", err.Error())
}

func (s *JWTDSuite) TestRequestEchoAdminPageWrongService() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-2", "role", "admin")
	s.NoError(err)
	s.NotEmpty(token)
	_, err = s.EchoOneRequest("GET", "/admin", token, "")
	s.Error(err)
	s.Equal("http error: 401", err.Error())
}

func (s *JWTDSuite) TestRequestEchoAdminPageWrongProject() {
	token, err := s.GetToken("wrong", "admin", "admin", "http-echo-1", "role", "admin")
	s.Error(err)
	s.Empty(token)
	_, err = s.EchoOneRequest("GET", "/admin", token, "")
	s.Error(err)
	s.Equal("http error: 401", err.Error())
}

func (s *JWTDSuite) TestRequestEchoUserPageCorrect() {
	token, err := s.GetToken("default", "user", "user", "http-echo-1", "role", "user")
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.EchoOneRequest("GET", "/user", token, "")
	s.NoError(err)
	s.Equal("/user", resp.URL)
}

func (s *JWTDSuite) TestRequestEchoUserPageAsAdmin() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-1", "role", "user")
	s.NoError(err)
	s.NotEmpty(token)
	resp, err := s.EchoOneRequest("GET", "/user", token, "")
	s.NoError(err)
	s.Equal("/user", resp.URL)
}
