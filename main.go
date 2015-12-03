package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
)

type cfg struct {
	interval time.Duration
	perms    uint

	csrfile string
	cfgfile string

	certpath string
	certname string
}

var c = &cfg{}
var s = &cli.Config{}

func init() {
	flag.DurationVar(&c.interval, "interval", 0, "repeat interval")
	flag.UintVar(&c.perms, "perms", 0600, "permissions for resulting dirs/files")

	flag.StringVar(&c.csrfile, "csrfile", "", "path to csr json file")
	flag.StringVar(&c.cfgfile, "config", "", "path to file with auth info")

	flag.StringVar(&c.certpath, "certpath", "", "path of resulting cert and key")
	flag.StringVar(&c.certname, "certname", "", "name of resulting cert and key (${certname}.pem, ${certname}-key.pem)")

	flag.StringVar(&s.Hostname, "hostname", "", "hostname for the cert, comma separated")
	flag.StringVar(&s.Profile, "profile", "", "signing config profile")
	flag.StringVar(&s.Label, "label", "", "signing config label")
	flag.StringVar(&s.Remote, "remote", "", "remote cfssl server")
}

func main() {
	flag.Parse()

	if c.interval == 0 {
		log.Fatal("interval is required and must be > 0")
	}

	if len(c.csrfile) == 0 {
		log.Fatal("csrfile is required")
	}

	if len(c.certpath) == 0 {
		log.Fatal("path is required")
	}

	if len(c.certname) == 0 {
		log.Fatal("certname is required")
	}

	if len(s.Remote) == 0 {
		log.Fatal("remote is required")
	}

	log.Println("cfssl-sidecar")
	log.Printf("operating with interval %s", c.interval)

	if err := os.MkdirAll(c.certpath, os.FileMode(c.perms)); err != nil {
		log.Fatal(err)
	}

	if len(c.cfgfile) > 0 {
		loadedCfg, err := config.LoadFile(c.cfgfile)
		if err != nil {
			log.Fatal(err)
		}
		s.CFG = loadedCfg
	} else {
		log.Println("WARNING: no config file specified.")
	}

	log.Println("updating key and csr.")
	keypath, csrpath, err := createKey(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("key(%s) and csr(%s) updated.", *keypath, *csrpath)

	s.CSRFile = *csrpath
	for {
		log.Println("updating certificate.")
		certpath, err := createCert(c, s)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("cert(%s) updated.", *certpath)

		log.Printf("sleeping for interval %s.", c.interval)
		time.Sleep(c.interval)
	}
}

func createKey(c *cfg) (*string, *string, error) {
	keypath := path.Join(c.certpath, fmt.Sprintf("%s-key.pem", c.certname))
	csrpath := path.Join(c.certpath, fmt.Sprintf("%s.csr", c.certname))

	_, keyerr := os.Stat(keypath)
	_, csrerr := os.Stat(csrpath)

	if keyerr != nil && csrerr != nil {
		log.Println("key and csr already exist")
		return &keypath, &csrpath, nil
	}

	bytes, err := ioutil.ReadFile(c.csrfile)
	if err != nil {
		return nil, nil, err
	}

	req := csr.CertificateRequest{KeyRequest: csr.NewBasicKeyRequest()}
	if err := json.Unmarshal(bytes, &req); err != nil {
		return nil, nil, err
	}

	var keybytes, csrbytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrbytes, keybytes, err = g.ProcessRequest(&req)
	if err != nil {
		return nil, nil, err
	}

	if err := ioutil.WriteFile(keypath, keybytes, 0600); err != nil {
		return nil, nil, err
	}

	if err := ioutil.WriteFile(csrpath, csrbytes, 0600); err != nil {
		return nil, nil, err
	}

	return &keypath, &csrpath, nil
}

func createCert(c *cfg, s *cli.Config) (*string, error) {
	signr, err := sign.SignerFromConfig(*s)
	if err != nil {
		return nil, err
	}

	csrbytes, err := ioutil.ReadFile(s.CSRFile)
	if err != nil {
		return nil, err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrbytes),
		Hosts:   signer.SplitHosts(s.Hostname),
		Profile: s.Profile,
		Label:   s.Label,
	}

	cert, err = signr.Sign(signReq)
	if err != nil {
		return nil, err
	}

	certpath := path.Join(c.certpath, fmt.Sprintf("%s.pem", c.certname))
	if err := ioutil.WriteFile(certpath, cert, 0600); err != nil {
		return nil, err
	}

	return &certpath, nil
}
