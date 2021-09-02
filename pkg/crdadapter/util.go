package crdadapter

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/casbin/casbin/v2/util"
	"strings"
)

func removeStringAndLineBreaks(old string) string {
	tmp := strings.Replace(old, "\n", "", -1)
	tmp = strings.Replace(tmp, " ", "", -1)
	return tmp
}

func policyToString(ptype string, rule []string) string {
	var tmp bytes.Buffer
	tmp.WriteString(ptype + ", ")
	tmp.WriteString(util.ArrayToString(rule))
	return removeStringAndLineBreaks(tmp.String())
}

func generatePolicyName(policy string) string {
	d := []byte(policy)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
