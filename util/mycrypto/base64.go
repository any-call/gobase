package mycrypto

import "encoding/base64"

func Base64Encode(inputData []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(inputData)
	return encodeString
}

func Base64DecodeWithString(inputStr string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(inputStr)
	if err != nil {
		return nil
	}

	return decodeBytes
}

func Base64DecodeWithBytes(inputBytes []byte) []byte {
	desBytes := make([]byte, len(inputBytes)*2)
	_, err := base64.StdEncoding.Decode(desBytes, inputBytes)
	if err != nil {
		return nil
	}

	return desBytes
}
