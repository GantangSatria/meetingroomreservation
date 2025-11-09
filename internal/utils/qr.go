package utils

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return Base64Encode(png), nil
}
