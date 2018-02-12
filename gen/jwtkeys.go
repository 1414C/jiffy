package gen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 2048
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
	RSABits uint // 0
	// ECDSACurve string   // "P256" || "P384" || "P521" // deprecated
	ECDSA     []string // {"256", "384", "521"}
	RSA       []string // {"256", "384", "512"},
	TargetDir string   // target directory to write keys to
}

// GenerateJWTKeys creates some ECDSA and RSA keys
func (keyc *KeyConfig) GenerateJWTKeys(conf *Config) error {

	var privatePEM []byte
	var publicPEM []byte

	// check for and create directory
	chkCrtDir := func(baseDir, subDir string) string {
		dir := ""
		if subDir == "" {
			dir = baseDir
		} else {
			dir = baseDir + subDir
		}
		_, err := os.Stat(dir)
		if err != nil {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				panic("unable to create jwt key-pair directories" + err.Error())
			}
		}
		return dir
	}

	// if the jwt key-pair dirs don't exist, create them
	chkCrtDir(keyc.TargetDir, "")

	// generate ECDSA key-pairs
	for _, b := range keyc.ECDSA {

		curve := "P" + b

		priv, privPEM, err := privateECDSAKey(curve)
		if err != nil {
			return err
		}

		_, pubPEM, err := publicECDSAKey(priv)
		if err != nil {
			return err
		}

		privatePEM = privPEM
		publicPEM = pubPEM

		// write ECDSA key files to the filesystem
		td := chkCrtDir(keyc.TargetDir, fmt.Sprintf("/ecdsa%s", b))
		privFile := fmt.Sprintf("%s/ec%s.priv.pem", td, b)
		err = ioutil.WriteFile(privFile, privatePEM, 0755)
		if err != nil {
			return err
		}
		log.Println("generated JWT private key:", privFile)

		pubFile := fmt.Sprintf("%s/ec%s.pub.pem", td, b)
		err = ioutil.WriteFile(pubFile, publicPEM, 0755)
		if err != nil {
			return err
		}
		log.Println("generated JWT public key:", pubFile)

		switch b {
		case "256":
			conf.ECDSA256PrivKeyFile = privFile
			conf.ECDSA256PubKeyFile = pubFile
		case "384":
			conf.ECDSA384PrivKeyFile = privFile
			conf.ECDSA384PubKeyFile = pubFile
		case "521":
			conf.ECDSA521PrivKeyFile = privFile
			conf.ECDSA521PubKeyFile = pubFile
		default:
			// do nothing
		}
	}

	// generate RSA key-pairs. no need to look at the shaxxx hash here
	if keyc.RSABits == 0 {
		log.Printf("RSA bits set to zero - defaulting to 2048\n")
		keyc.RSABits = 2048
	}

	for _, b := range keyc.RSA {

		privateRSAKey, privateRSAPEM, err := privateRSAKey(keyc.RSABits)
		if err != nil {
			return err
		}

		_, publicRSAPEM, err := publicRSAKey(privateRSAKey)
		if err != nil {
			return err
		}

		// write RSA key files to the filesystem
		td := chkCrtDir(keyc.TargetDir, fmt.Sprintf("/rsa%s", b))
		privFile := fmt.Sprintf("%s/rsa%s.priv.pem", td, b)
		err = ioutil.WriteFile(privFile, privateRSAPEM, 0755)
		if err != nil {
			return err
		}
		log.Println("generated JWT private key:", privFile)

		pubFile := fmt.Sprintf("%s/rsa%s.pub.pem", td, b)
		err = ioutil.WriteFile(pubFile, publicRSAPEM, 0755)
		if err != nil {
			return err
		}
		log.Println("generated JWT public key:", pubFile)

		switch b {
		case "256":
			conf.RSA256PrivKeyFile = privFile
			conf.RSA256PubKeyFile = pubFile
		case "384":
			conf.RSA384PrivKeyFile = privFile
			conf.RSA384PubKeyFile = pubFile
		case "512":
			conf.RSA512PrivKeyFile = privFile
			conf.RSA512PubKeyFile = pubFile
		default:
			// do nothing
		}
	}
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
	return public, pemPublic, nil
}

// privateRSAKey returns a private RSA key of the prescribed length
func privateRSAKey(l uint) (priv *rsa.PrivateKey, pemPriv []byte, err error) {

	// check the length
	if l != 1024 &&
		l != 2048 &&
		l != 4096 {
		return nil, nil, fmt.Errorf("rsa keys must be 1024, 2048 or 4096 bits.  %v is not a valid length", l)
	}

	// generate a private key
	priv, err = rsa.GenerateKey(rand.Reader, int(l))
	if err != nil {
		return nil, nil, err
	}

	// marshal the private key into DER format
	derKey := x509.MarshalPKCS1PrivateKey(priv)
	if derKey == nil {
		return nil, nil, fmt.Errorf("unable to marshall RSA private-key into DER format")
	}

	// create a pem-block using the private key in DER format
	keyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derKey,
	}

	// encode the private pem-block to memory (you could encode directly to an io.Writer here...)
	pemPriv = pem.EncodeToMemory(keyBlock)
	return priv, pemPriv, nil
}

// publicECDSAKey returns a public ecdsa key based on the private-key
func publicRSAKey(priv *rsa.PrivateKey) (public *rsa.PublicKey, pemPublic []byte, err error) {

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
	return public, pemPublic, nil
}
