package twoFA

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	qr "github.com/sec51/qrcode"
)

// CreateSecret generate random SHA1 secret
func CreateSecret() []byte {
	key := make([]byte, crypto.SHA1.Size())
	rand.Seed(time.Now().UnixNano())
	rand.Read(key)

	return key
}

// GenerateQRCode generate QR code encoded as base64
func GenerateQRCode(link string) (string, error) {
	// generate qr code
	code, err := qr.Encode(link, qr.Q)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(code.PNG())), nil
}

// GenerateLink generate URI in follow format:
// Example: otpauth://totp/rightnao:qwe_login?secret=qweqwe123&issuer=rightnao&digits=6&period=30&algorithm=SHA1
func GenerateLink(secret []byte, login string) string {
	// generate link
	secretEnc := base32.StdEncoding.EncodeToString(secret)
	u := url.URL{}
	v := url.Values{}
	u.Scheme = "otpauth"
	u.Host = "totp"
	u.Path = fmt.Sprintf("%s:%s", url.QueryEscape("rightnao"), login)
	v.Add("secret", secretEnc)
	// v.Add("counter", fmt.Sprintf("%d", otp.getIntCounter()))
	v.Add("issuer", "rightnao")
	v.Add("digits", "6")
	v.Add("period", "30")
	v.Add("algorithm", "SHA1")

	u.RawQuery = v.Encode()
	return u.String()
}

// CheckCode ...
func CheckCode(secret []byte, code string) bool {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64((time.Now().Unix() / 30))) // 30 is interval

	hash := hmac.New(sha1.New, secret)
	hash.Write(bs)
	h := hash.Sum(nil)
	o := (h[19] & 15)
	var header uint32
	r := bytes.NewReader(h[o : o+4])
	binary.Read(r, binary.BigEndian, &header)
	h12 := (int(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(int(h12))

	if func(s string) string {
		if len(s) == 6 {
			return s
		}
		for i := (6 - len(s)); i > 0; i-- {
			s = "0" + s
		}
		return s
	}(otp) != code {
		return false
	}

	return true
}

// GetKey returns secret encoded in base32
func GetKey(secret []byte) string {
	key := base32.StdEncoding.EncodeToString(secret)
	return key
}
