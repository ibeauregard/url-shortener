package main

import (
	"os"
)

var AppHost = os.Getenv("APP_HOST")

const AppScheme = "http"

// This alphabet will be used to generate the paths of the shortened URLs.
// It consists of the decimal digits and of the uppercase and lowercase letters, plus some special characters.
// Characters that could cause ambiguity or generate offensive words were removed.
const alphabet = "23456789BCDFGHJKLMNPQRSTVWXYZbcdfghjkmnpqrstvwxyz-_~!$&=@"
const alphabetLength = uint(len(alphabet))
