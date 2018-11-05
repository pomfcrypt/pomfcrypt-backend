package util

import (
	"crypto/rand"
	"github.com/sirupsen/logrus"
	"io"
)

var AvailableCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func Generate(length int) string {
	uuid := make([]byte, length)
	random := make([]byte, length+(length/4))
	charLength := byte(len(AvailableCharacters))
	maximumRb := byte(256 - (256 % len(AvailableCharacters)))
	i := 0
	for {
		// Read from rand.Reader
		if _, err := io.ReadFull(rand.Reader, random); err != nil {
			logrus.Fatal("Failed to generate UUID: ", err)
		}
		// Iterate over the random slice
		for _, c := range random {
			// Check if larger than the maximum of the rb
			if c >= maximumRb {
				continue
			}
			// Add a random char to the uuid
			uuid[i] = AvailableCharacters[c%charLength]
			i++
			// Check if the output is as long as specified
			if i == length {
				return string(uuid)
			}
		}
	}
}
