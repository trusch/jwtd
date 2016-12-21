package main

import (
	"crypto/tls"
	"net"
	"net/http"
)

func ListenAndServeTLSSNI(addr string, certs []*TLSConfig, handler http.Handler) error {
	httpsServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	config := &tls.Config{}
	var err error
	config.Certificates = make([]tls.Certificate, len(certs))
	for i, v := range certs {
		config.Certificates[i], err = tls.LoadX509KeyPair(v.Cert, v.Key)
		if err != nil {
			return err
		}
	}

	config.BuildNameToCertificate()

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(conn, config)
	return httpsServer.Serve(tlsListener)
}
