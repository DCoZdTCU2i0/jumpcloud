package encoding

import "crypto/sha512"
import "encoding/base64"

func Conversion(input string) string {
	bytes := []byte(input)
	checksum := sha512.Sum512(bytes)
	output := base64.StdEncoding.EncodeToString(checksum[:])
	return output
}
