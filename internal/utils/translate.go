package utils

import (
	"fmt"
	"github.com/mozillazg/go-unidecode"
	"strings"
)

func Transliterate(input string) string {
	fmt.Println(input)
	fmt.Println(strings.TrimSpace(unidecode.Unidecode(input)))

	return strings.TrimSpace(unidecode.Unidecode(input))
}
