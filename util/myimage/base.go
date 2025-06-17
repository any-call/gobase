package myimage

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

func DecodeBase64PNG(data string) (image.Image, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(decoded))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func DecodeBase64JPG(data string) (image.Image, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	img, err := jpeg.Decode(bytes.NewReader(decoded))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func DecodeBase64Auto(data string) (image.Image, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(decoded)

	// 检查前几个字节判断格式
	head := make([]byte, 512)
	_, err = buf.Read(head)
	if err != nil {
		return nil, err
	}

	// 重置 buffer 指针
	buf.Seek(0, 0)

	switch {
	case bytes.HasPrefix(head, []byte("\x89PNG")):
		return png.Decode(buf)
	case bytes.HasPrefix(head, []byte("\xff\xd8\xff")):
		return jpeg.Decode(buf)
	default:
		return nil, errors.New("不支持的图片格式")
	}
}
