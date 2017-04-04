package middleware

// Note that I import the versions bundled with vulcand. That will make our lives easier, as we'll use exactly the same versions used
// by vulcand. We are escaping dependency management troubles thanks to Godep.
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/trusch/jwtd/jwt"
	"github.com/vulcand/vulcand/plugin"
)

const (
	// Type of middleware
	Type = "jwtd"
	// UserHeader is the header used to store user claims for
	// downstream services
	UserHeader = "X-USER"
)

type LabelSet map[string]string

func NewLabelSetFromString(str string) LabelSet {
	res := make(LabelSet)
	parts := strings.Split(str, ",")
	for _, part := range parts {
		leftAndRight := strings.Split(part, "=")
		if len(leftAndRight) == 2 {
			res[leftAndRight[0]] = leftAndRight[1]
		}
	}
	return res
}

func (set LabelSet) String() string {
	res := ""
	for key, val := range set {
		res += key + "=" + val + ","
	}
	res = res[:len(res)-1]
	return res
}

// GetSpec is part of the Vulcan middleware interface
func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,       // A short name for the middleware
		FromOther: FromOther,  // Tells vulcand how to rcreate middleware from another one (this is for deserialization)
		FromCli:   FromCli,    // Tells vulcand how to create middleware from command line tool
		CliFlags:  CliFlags(), // Vulcand will add this flags to middleware specific command line tool
	}
}

// JwtMiddleware struct holds configuration parameters and is used to
// serialize/deserialize the configuration from storage engines.
type JwtMiddleware struct {
	PublicKey interface{}
	Service   string
	Required  LabelSet
}

// JwtHandler is the HTTP handler for the JWT middleware
type JwtHandler struct {
	cfg  JwtMiddleware
	next http.Handler
}

// This function will be called each time the request hits the location with this middleware activated
func (a *JwtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Let OPTIONS go on by
	if r.Method == "OPTIONS" {
		a.next.ServeHTTP(w, r)
		return
	}

	log.Print("got request in middleware")

	// Pass the request to the next middleware in chain
	a.next.ServeHTTP(w, r)
}

// New is optional but handy, used to check input parameters when creating new middlewares
func New(publicKey interface{}, service string, required LabelSet) (*JwtMiddleware, error) {
	return &JwtMiddleware{PublicKey: publicKey, Service: service, Required: required}, nil
}

// NewHandler is important, it's called by vulcand to create a new handler from the middleware config and put it into the
// middleware chain. Note that we need to remember 'next' handler to call
func (c *JwtMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	return &JwtHandler{next: next, cfg: *c}, nil
}

// String() will be called by loggers inside Vulcand and command line tool.
func (c *JwtMiddleware) String() string {
	return fmt.Sprintf("service=%v required=%v", c.Service, c.Required)
}

// FromOther Will be called by Vulcand when engine or API will read the middleware from the serialized format.
// It's important that the signature of the function will be exactly the same, otherwise Vulcand will
// fail to register this middleware.
// The first and the only parameter should be the struct itself, no pointers and other variables.
// Function should return middleware interface and error in case if the parameters are wrong.
func FromOther(c JwtMiddleware) (plugin.Middleware, error) {
	return New(c.PublicKey, c.Service, c.Required)
}

// FromCli constructs the middleware from the command line
func FromCli(c *cli.Context) (plugin.Middleware, error) {
	var (
		pubKey   interface{}
		service  string
		required LabelSet
	)
	keyFile := c.String("key")
	if keyFile == "" {
		return nil, errors.New("supply a public key file")
	}
	k, err := jwt.LoadPublicKey(keyFile)
	if err != nil {
		return nil, err
	}
	pubKey = k
	service = c.String("service")
	if service == "" {
		return nil, errors.New("supply a service identifier")
	}
	reqStr := c.String("require")
	required = NewLabelSetFromString(reqStr)
	if len(required) == 0 {
		return nil, errors.New("supply required labels")
	}
	return New(pubKey, service, required)
}

// CliFlags will be used by Vulcand construct help and CLI command for the vctl command
func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "key, k", Usage: "Path to file with Public Key"},
		cli.StringFlag{Name: "service, s", Usage: "service id"},
		cli.StringFlag{Name: "require, r", Usage: "required label set"},
	}
}
