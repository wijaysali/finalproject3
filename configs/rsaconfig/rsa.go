package rsaconfig

import (
	"MyGram/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

var PrivateKey *rsa.PrivateKey

func CreateRsaKey() {
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	// Write private key to file.
	if err := ioutil.WriteFile(viper.GetString("RSA_PRIVATE_KEY_PATH"), keyPEM, 0700); err != nil {
		panic(err)
	}

	// Write public key to file.
	if err := ioutil.WriteFile(viper.GetString("RSA_PUBLIC_KEY_PATH"), pubPEM, 0755); err != nil {
		panic(err)
	}
}

func ConfigRSA() {
	if !utils.IsFileExists(viper.GetString("RSA_PRIVATE_KEY_PATH")) ||
		!utils.IsFileExists(viper.GetString("RSA_PUBLIC_KEY_PATH")) {
		CreateRsaKey()
	}

	priv, err := ioutil.ReadFile(viper.GetString("RSA_PRIVATE_KEY_PATH"))
	if err != nil {
		panic(errors.New("no RSA private key found"))
	}

	privPem, _ := pem.Decode(priv)

	if privPem.Type != "RSA PRIVATE KEY" {
		panic(fmt.Errorf("RSA private key is of the wrong type : %v", privPem.Type))
		//return errors.New(fmt.Sprintf("RSA private key is of the wrong type", privPem.Type))
	}

	privPemBytes := privPem.Bytes
	var parsedKey interface{}
	//PKCS1
	//x509.ParsePKCS8PrivateKey()
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		//If what you are sitting on is a PKCS#8 encoded key
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			panic(fmt.Errorf("unable to parse RSA private key: %v", err))
		}
	}

	var ok bool
	PrivateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		panic(errors.New("unable to parse RSA private key"))
	}

	pub, err := ioutil.ReadFile(viper.GetString("RSA_PUBLIC_KEY_PATH"))
	if err != nil {
		panic(fmt.Errorf("no RSA public key found: %v", err))
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		panic(errors.New("RSA public key not in pem format"))
	}

	if pubPem.Type != "RSA PUBLIC KEY" {
		panic(fmt.Errorf("RSA public key is of the wrong type: %v", pubPem.Type))
	}

	if parsedKey, err = x509.ParsePKCS1PublicKey(pubPem.Bytes); err != nil {
		panic(fmt.Errorf("unable to parse RSA public key: %v", err))
	}

	var pubKey *rsa.PublicKey
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		panic(errors.New("unable to parse RSA public key"))
	}

	PrivateKey.PublicKey = *pubKey
}

// GenRSA returns a new RSA key of bits length
func GenRSA(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	return key, err
}
