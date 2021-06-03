package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

/*
eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQyNzQ1NTQuNzM4OTI4LCJpYXQiOjE2MjI3Mzg1NTQuNzM4OTMyLCJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIxMjM0NTY3Nzg5IiwiUm9sZXMiOlsiQURNSU4iXX0.btlXINXk0co2munSzWNQtumxzkAz4dKvsjTNHLzQpa7lsbCUvIDo8gQZ9Z5r6QXJAHihuhJfAX66SOibcQIymE-UD_Iu_Hy0HbYpkYm3WZ0NgHHQAdSjqCqqxmWUby6ERErrYNwttf38HDk-MaCoIE5LCmutPSIPFmtBlQLT3aR-EAdKOv9Odpt1j8JzHdGkY8V42HlH55SjGMRo_-e-kcwWYtiojfU9vbRI5uf0Z0tdgLaoSEuKUedVvmq-9P8_tog9Hu4FJMJGrWiLPIRFAbhTVDmxz_Zkn8SKv3CxAiqrL7H_WHu9jM2urvfvCb5tItrwWMOxPRFYQOocGlkUKA
*/

func main() {
	err := GenToken() //GenKey()

	if err != nil {
		log.Println(err)
	}
}

func GenToken() error {

	pkf, err := os.Open("zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")
	if err != nil {
		return errors.Wrap(err, "opening PEM private key file")
	}
	defer pkf.Close()
	privatePEM, err := io.ReadAll(io.LimitReader(pkf, 1024*1024))
	if err != nil {
		return errors.Wrap(err, "reading PEM private key file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return errors.Wrap(err, "parsing PEM into private key")
	}

	claims := struct {
		jwt.StandardClaims
		Roles []string
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "1234567789",
			ExpiresAt: jwt.At(time.Now().Add(8760 * time.Hour)),
			IssuedAt:  jwt.Now(),
		},
		Roles: []string{"USER"},
	}

	method := jwt.GetSigningMethod("RS256")
	if method == nil {
		return errors.Errorf("unknown algorithm %v", "RS256")
	}

	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	str, err := token.SignedString(privateKey)
	if err != nil {
		return errors.Wrap(err, "signing token")
	}

	fmt.Println("======================================")
	fmt.Println(str)
	fmt.Println("======================================")

	// ----------------------------------------------------------------------

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}), jwt.WithAudience("student"))

	var clm struct {
		jwt.StandardClaims
		Roles []string
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// kid := t.Header["kid"].(string)
		return &privateKey.PublicKey, nil
	}

	tkn, err := parser.ParseWithClaims(str, &clm, keyFunc)
	if err != nil {
		return err
	}

	if !tkn.Valid {
		return errors.New("invalid token")
	}

	fmt.Println("Validated", tkn.Claims, tkn.Header)

	return nil
}

func GenKey() error {

	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a file for the private key information in PEM form.
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer privateFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	// -----------------------------------------------------------------

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "marshaling public key")
	}

	// Create a file for the public key information in PEM form.
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return errors.Wrap(err, "creating public file")
	}
	defer publicFile.Close()

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return errors.Wrap(err, "encoding to public file")
	}

	fmt.Println("private and public key files generated")
	return nil
}
