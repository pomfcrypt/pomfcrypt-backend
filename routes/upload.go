package routes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/pomfcrypt/pomfcrypt-backend/model"
	"github.com/pomfcrypt/pomfcrypt-backend/util"
	"golang.org/x/crypto/pbkdf2"
	"html"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// Upload doc
// @Summary Upload route
// @Description Upload a new file
// @Accept mpfd
// @Accept json
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} model.FileResponse
// @Failure 500 {object} routes.APIErrorMessage
// @Failure 400 {object} routes.APIErrorMessage
// @Router /data [put]
// @Param  file body string true "Uploaded file"
func (ctl *Controller) Upload(c iris.Context) {
	// Check if a password and a file is given, first
	password := c.PostValue("password")
	if password == "" {
		NewAPIError("Failed to retrieve password", errors.New("password was not set")).Throw(c, 400)
		return
	}
	// Get the posted file
	_, postFileHeader, err := c.FormFile("file")
	if err != nil {
		NewAPIError("Failed to receive file", err).Throw(c, 400)
		return
	}
	if postFileHeader.Size > ctl.Settings.MaxSize {
		NewAPIError("File is too large", err).Throw(c, 400)
		return
	}
	file, err := postFileHeader.Open()
	if err != nil {
		NewAPIError("Failed to open file", err).Throw(c, 500)
		return
	}
	fileExists := true
	// Find a random filename that is not already allocated
	var fileName string
	for fileExists == true {
		fileName = util.Generate(ctl.Settings.FilenameLength)
		// Check if file actually exists
		_, err := os.Stat(ctl.BuildPath(fileName))
		if err != nil && !os.IsNotExist(err) {
			NewAPIError("Failed to retrieve information about existing file", err).Throw(c, 500)
			return
		} else {
			fileExists = false
		}
	}
	// Close the file after the upload is done
	defer file.Close()
	// Sanitize input string
	postFileHeader.Filename = html.EscapeString(postFileHeader.Filename)

	// Create the output file stream
	outFile, err := os.Create(ctl.BuildPath(fileName))
	if err != nil {
		NewAPIError("Failed to create file", err).Throw(c, 500)
		return
	}

	// Close the outFile after the upload is done
	defer outFile.Close()

	// Create the length header
	lengthHeader := make([]byte, 2)
	// Put the (sanitized) file name into the length header
	binary.BigEndian.PutUint16(lengthHeader, uint16(len([]byte(postFileHeader.Filename))))

	// Create a headers variable containing all relevant headers for the file
	var headers []byte
	// Append the generated lengthHeader to the headers first
	headers = append(headers, lengthHeader...)
	// Append the name of the file
	headers = append(headers, []byte(postFileHeader.Filename)...)

	// Create the integrity-ensuring pbkdf2 key out of the salt and the given password utilizing SHA256
	key := pbkdf2.Key([]byte(password), []byte(ctl.Settings.Salt), 4096, 16, sha256.New)

	// Create a new AES cypher
	aesBlock, err := aes.NewCipher(key)

	if err != nil {
		NewAPIError("Failed to create AES key", err).Throw(c, 500)
		return
	}

	// Create the deriving HMAC
	keyHmac := hmac.New(sha256.New, key)

	// Create the IV key from the AES BlockSize
	iv := make([]byte, aes.BlockSize)
	// Create random from rand.Reader into IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		NewAPIError("Failed to read random", err).Throw(c, 500)
		return
	}
	// Seek 32 bytes for the headers
	if _, err := outFile.Seek(32, io.SeekStart); err != nil {
		NewAPIError("Failed to seek in output file", err).Throw(c, 500)
		return
	}

	// Create the cipher stream
	cipherStream := cipher.NewCFBEncrypter(aesBlock, iv)

	// Create a writer for the key HMAC
	streamOut := io.MultiWriter(keyHmac, outFile)
	// Write to the created writer
	if _, err := streamOut.Write(iv); err != nil {
		NewAPIError("Failed to write to stream", err).Throw(c, 500)
		return
	}
	// Create a StreamWriter for the cipherStream
	streamWriter := &cipher.StreamWriter{S: cipherStream, W: streamOut}
	// Create the CFB Reader
	cfbReader := io.MultiReader(bytes.NewReader(headers), file)
	// Copy from the streamWriter to the cfbReader
	if _, err := io.Copy(streamWriter, cfbReader); err != nil {
		NewAPIError("Failed to write to cfbReader", err).Throw(c, 500)
		return
	}
	// Get the calculated HMAC sum
	sum := keyHmac.Sum(nil)
	// Write the sum to the file
	if _, err := outFile.Write(sum); err != nil {
		NewAPIError("Failed to write to output file", err).Throw(c, 500)
		return
	}

	// Stat the output file
	fileInfo, err := outFile.Stat()
	if err != nil {
		NewAPIError("Could not stat output file", err).Throw(c, 500)
		return
	}

	// Check the file size
	if fileInfo.Size() > ctl.Settings.MaxSize {
		NewAPIError("The file size is too large!", errors.New("the file size exceeds "+strconv.Itoa(int(ctl.Settings.MaxSize))))
		return
	}

	// Read the final output
	finalBytes, err := ioutil.ReadAll(outFile)
	if err != nil {
		NewAPIError("Failed to read output file", err).Throw(c, 500)
		return
	}

	c.StatusCode(201)
	hashSum := sha256.Sum256(finalBytes)
	c.JSON(model.FileResponse{
		Filename:   fileName,
		UploadedAt: fileInfo.ModTime().Unix(),
		Hash:       fmt.Sprintf("%x", hashSum),
	})
}
