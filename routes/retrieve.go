package routes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func (ctl *Controller) Retrieve(c iris.Context) {
	// Retrieve and validate password
	password := c.PostValue("password")
	if len(password) < 1 {
		NewAPIError("Failed to get password", errors.New("password too short")).Throw(c, 400)
		return
	}
	key := pbkdf2.Key([]byte(password), []byte(ctl.Settings.Salt), 4096, 16, sha1.New)

	fname := c.Params().Get("name")
	ciphertext, err := ioutil.ReadFile(ctl.BuildPath(fname))
	if err != nil {
		logrus.Warn(err)
		if os.IsNotExist(err) {
			NewAPIError("File not found", err).Throw(c, 404)
			return
		}
		NewAPIError("Failed to read file", err).Throw(c, 400)
		return
	}
	if string(ciphertext) == "deleted\n" || string(ciphertext) == "deleted" {
		NewAPIError("File deleted due to legal reasons", errors.New("file was deleted")).Throw(c, 410)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		NewAPIError("Failed to build AES block", err).Throw(c, 500)
		return
	}
	if len(ciphertext) < aes.BlockSize {
		NewAPIError("ciphertext is too short", err).Throw(c, 500)
		return
	}
	h := hmac.New(sha256.New, key)
	h.Write(ciphertext[32:])
	mac := h.Sum(nil)
	filemac := ciphertext[:32]
	iv := ciphertext[32 : 32+aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize+32:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	size := binary.BigEndian.Uint16(ciphertext[:2])
	filename := ciphertext[2 : size+2]
	content := ciphertext[size+2:]
	if !(hmac.Equal(mac, filemac)) {
		NewAPIError("Failed to validate HMAC", errors.New("hmac mismatch")).Throw(c, 500)
		return
	}
	c.Header("Content-Disposition", "inline; filename='"+string(filename[:])+"'")
	c.Header("Content-Length", strconv.Itoa(len(content)))
	c.Header("Content-Type", http.DetectContentType(content[:512]))
	c.Write(content)

}
