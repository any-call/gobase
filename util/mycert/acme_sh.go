package mycert

import "strings"

type acmeSh struct {
	autoCreateNginxConf bool
	domainList          []string
	*strings.Builder
}

func NewAcmeSH() *acmeSh {
	return &acmeSh{
		Builder: &strings.Builder{},
	}
}
