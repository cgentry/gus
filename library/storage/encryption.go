// This package provides field encryption for any driver. The encryption
// should be enabled using the driver options
package storage

type Encrypter interface {
	Encrypt( string ) string
	Decrypt( string ) string

	SetKey( key string )
}
