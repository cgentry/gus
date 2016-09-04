// Package storage  provides field encryption for any driver. The encryption
// should be enabled using the driver options
package storage

// Encrypter provides the interface that storage classes need to support encryption
type Encrypter interface {
	Encrypt(string) string
	Decrypt(string) string

	SetKey(key string)
}
