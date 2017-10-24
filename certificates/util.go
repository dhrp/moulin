package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	demoKeyPair  *tls.Certificate
	demoCertPool *x509.CertPool
)

// How to generate your own self-signed certificate:
// openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt

// GetCert returns a certicicate pair and pool
func GetCert() (*tls.Certificate, *x509.CertPool) {
	crtPath, _ := filepath.Abs("./certificates/server.crt")
	serverCrt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		log.Fatal(err)
	}
	keyPath, _ := filepath.Abs("./certificates/server.key")
	serverKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	pair, err := tls.X509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Fatal(err)
	}
	demoKeyPair = &pair
	demoCertPool = x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM(serverCrt)
	if !ok {
		log.Fatal("bad certs")
	}

	return demoKeyPair, demoCertPool
}
