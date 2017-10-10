package gen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"os"
	// "crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 1024
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub
//
// 	Generate ECDSA signing files via shell (preferred over RSA)
//
//	$ openssl ecparam -genkey -name secp384r1 -noout -out private.pem
//	$ openssl ec -in private.pem -pubout -out public.pem

// see https://golang.org/src/crypto/tls/generate_cert.go for some interesting
// (but not really related) certificate stuff.

// KeyConfig holds the basic certificate configuration
type KeyConfig struct {
	RSABits    int    // 0
	ECDSACurve string // "P256" || "P384" || "P521"
	TargetDir  string // target directory to write keys to
}

// GenerateJWTKeys creates some ECDSA or RSA keys
func (keyc *KeyConfig) GenerateJWTKeys() error {

	var privatePEM []byte
	var publicPEM []byte

	// create ECDSA-based-keys
	if keyc.ECDSACurve != "" {
		priv, privPEM, err := privateECDSAKey(keyc.ECDSACurve)
		if err != nil {
			return err
		}

		_, pubPEM, err := publicECDSAKey(priv)
		if err != nil {
			return err
		}
		privatePEM = privPEM
		publicPEM = pubPEM
	} else {
		// RSA: TODO
		return fmt.Errorf("GenerateJWTKeys does not support RSA key generation at the moment")
	}

	// if TargetDir does not exist, create it
	_, err := os.Stat(keyc.TargetDir)
	if err != nil {
		os.Mkdir(keyc.TargetDir, 0755)
	}

	// write key files to the local filesystem
	err = ioutil.WriteFile(keyc.TargetDir+"/private.pem", privatePEM, 0755)
	if err != nil {
		return err
	}
	log.Println("generated JWT private key:", keyc.TargetDir+"/private.pem")
	err = ioutil.WriteFile(keyc.TargetDir+"/public.pem", publicPEM, 0755)
	if err != nil {
		return err
	}
	log.Println("generated JWT public key:", keyc.TargetDir+"/public.pem")
	return nil
}

// privateECDSAKey generates a private ecdsa key based on the selected curve
func privateECDSAKey(curve string) (priv *ecdsa.PrivateKey, pemPriv []byte, err error) {

	// generate private key
	switch curve {
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	if err != nil {
		return nil, nil, err
	}

	// marshal the private key into DER format
	derKey, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	// create a pem-block using the private key in DER format
	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derKey,
	}

	// encode the private pem-block to memory (you could encode directly to an io.Writer here...)
	pemPriv = pem.EncodeToMemory(keyBlock)
	// fmt.Println("PEM Encoded Private EC KEY:")
	// fmt.Println(string(pemPriv))

	return priv, pemPriv, nil
}

// publicECDSAKey returns a public ecdsa key based on the private-key
func publicECDSAKey(priv *ecdsa.PrivateKey) (public *ecdsa.PublicKey, pemPublic []byte, err error) {

	// get the public-key from the private-key
	public = &priv.PublicKey

	// marshal the public key into DER format
	derBytes, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, nil, err
	}

	// create a pem-block using the public key in DER format
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	// encode the public pem-block to memory (you could encode directly to an io.Writer here...)
	pemPublic = pem.EncodeToMemory(block)
	// fmt.Println("PEM Encoded Public EC KEY:")
	// fmt.Println(string(pemPublic))

	return public, pemPublic, nil
}
