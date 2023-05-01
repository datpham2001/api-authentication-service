package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	ACCOUNT          = "ACCOUNT"
	ACCOUNT_LENGTH   = 6
	STRING_TO_GEN_ID = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenNanoID(alphabet string, length int) string {
	res, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		return ""
	}

	return res
}

func GenAccountID() string {
	genResult := "ACC" + GenNanoID(STRING_TO_GEN_ID, ACCOUNT_LENGTH)
	return genResult
}
