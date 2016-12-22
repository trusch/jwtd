package main

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenCorrect() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-1", []*Label{&Label{"role", "admin"}})
	s.NoError(err)
	s.NotEmpty(token, "admin token is empty")

}

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenWrongPassword() {
	token, err := s.GetToken("default", "admin", "wrong-password", "http-echo-1", []*Label{&Label{"role", "admin"}})
	s.Error(err)
	s.Equal("http error: 401", err.Error(), "wrong error")
	s.Empty(token, "admin token is not empty")
}

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenWrongLabel() {
	token, err := s.GetToken("default", "user", "user", "http-echo-1", []*Label{&Label{"role", "wrong-label"}})
	s.Error(err)
	s.Equal("http error: 401", err.Error(), "wrong error")
	s.Empty(token, "admin token is not empty")
}

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenWrongProject() {
	token, err := s.GetToken("wrong", "admin", "admin", "http-echo-1", []*Label{&Label{"role", "admin"}})
	s.Error(err)
	s.Equal("http error: 401", err.Error(), "wrong error")
	s.Empty(token, "admin token is not empty")
}

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenWrongService() {
	token, err := s.GetToken("default", "admin", "admin", "wrong-service", []*Label{&Label{"role", "admin"}})
	s.Error(err)
	s.Equal("http error: 401", err.Error(), "wrong error")
	s.Empty(token, "admin token is not empty")
}

func (s *JWTDSuite) TestGetHttpEchoOneAdminTokenAsUser() {
	token, err := s.GetToken("default", "user", "user", "http-echo-1", []*Label{&Label{"role", "admin"}})
	s.Error(err)
	s.Equal("http error: 401", err.Error(), "wrong error")
	s.Empty(token, "admin token is not empty")
}

func (s *JWTDSuite) TestGetHttpEchoOneUserTokenAsUser() {
	token, err := s.GetToken("default", "user", "user", "http-echo-1", []*Label{&Label{"role", "user"}})
	s.NoError(err)
	s.NotEmpty(token)
}

func (s *JWTDSuite) TestGetHttpEchoOneUserTokenAsAdmin() {
	token, err := s.GetToken("default", "admin", "admin", "http-echo-1", []*Label{&Label{"role", "user"}})
	s.NoError(err)
	s.NotEmpty(token)
}
