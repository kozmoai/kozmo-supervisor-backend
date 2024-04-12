package tokenvalidator

import (
	"encoding/base64"
	"fmt"
	"sort"

	"crypto/md5"

	"github.com/kozmoai/kozmo-supervisor-backend/src/utils/config"
)

type RequestTokenValidator struct {
	Config *config.Config
}

func NewRequestTokenValidator() *RequestTokenValidator {
	return &RequestTokenValidator{
		Config: config.GetInstance(),
	}
}

func (r *RequestTokenValidator) GenerateValidateToken(input ...string) string {
	return r.GenerateValidateTokenBySliceParam(input)
}

func (r *RequestTokenValidator) GenerateValidateTokenBySliceParam(input []string) string {
	var concatr string
	sort.Strings(input)
	for _, str := range input {
		concatr += str
	}
	concatr += r.Config.GetSecretKey()
	fmt.Printf("[DUMP] GenerateValidateTokenBySliceParam r.Config.GetSecretKey(): %+v\n", r.Config.GetSecretKey())
	hash := md5.Sum([]byte(concatr))
	var hashConverted []byte = hash[:]

	return base64.StdEncoding.EncodeToString(hashConverted)
}
