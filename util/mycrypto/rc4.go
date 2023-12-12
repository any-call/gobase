package mycrypto

import "crypto/rc4"

func RC4Encode(key, inputData []byte) []byte {
	rceObj1, _ := rc4.NewCipher(key)
	outputData := make([]byte, len(inputData))
	rceObj1.XORKeyStream(outputData, inputData)
	return outputData
}

func RC4Decode(key, inputData []byte) []byte {
	rceObj1, _ := rc4.NewCipher(key)
	outputData := make([]byte, len(inputData))
	rceObj1.XORKeyStream(outputData, inputData)
	return outputData
}
