package impl

import (
	"errors"
	"regexp"
	"strings"

	"d1y.io/neovideo/common/json"
	"github.com/tidwall/gjson"
)

var jiexiURLFuzzyReg = regexp.MustCompile(`(?i)(jiexi|jiexiurl|url)?=`)
var jiexiURLReg = regexp.MustCompile(`^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
var jiexiURLAndNoteReg = regexp.MustCompile(`^(\S*:?)\s*(https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]=)`)
var jiexiIgnoreReg = regexp.MustCompile(`^(//|;)`)
var symbolReg = regexp.MustCompile(`[,.:：，。]*$`)

func gjsonGGGetString(g gjson.Result, s ...string) (string, error) {
	for _, v := range s {
		s := g.Get(v).String()
		if len(s) >= 1 {
			return s, nil
		}
	}
	return "", errors.New("not match")
}

func parseJiexiWithLines(raw string) []JiexiParse {
	var result = make([]JiexiParse, 0)
	lines := strings.Split(raw, "\n")
	for _, item := range lines {
		s := strings.TrimSpace(item)
		if len(s) <= 6 || !jiexiURLFuzzyReg.MatchString(s) || jiexiIgnoreReg.MatchString(s) {
			continue
		}
		var i JiexiParse
		if jiexiURLReg.MatchString(s) {
			pl := strings.Split(s, "")
			if pl[len(pl)-1] == "=" {
				i.URL = s
			}
		} else {
			subs := jiexiURLAndNoteReg.FindStringSubmatch(s)
			if len(subs) == 3 {
				n, u := subs[1], subs[2]
				n = symbolReg.ReplaceAllString(n, "")
				i.URL = u
				i.Name = n
			}
		}
		result = append(result, i)
	}
	return result
}

func parseJiexiWithJSON(raw string) []JiexiParse {
	var result = make([]JiexiParse, 0)
	gs := gjson.Parse(raw)
	if gs.IsArray() {
		for _, item := range gs.Array() {
			if item.IsObject() {
				t, e1 := gjsonGGGetString(item, "name", "title")
				u, e2 := gjsonGGGetString(item, "url", "jiexi_url", "jiexi_url", "jiexiUrl") /* 不校验 url */
				if e1 != nil && e2 != nil {
					continue
				}
				result = append(result, JiexiParse{URL: u, Name: t})
			} else {
				if jiexiURLReg.MatchString(item.String()) {
					result = append(result, JiexiParse{URL: item.String()})
				}
			}
		}
	} /* else if gs.IsObject() {
	}*/
	return result
}

func ParseJiexi(raw string) []JiexiParse {
	var m = make(map[string]JiexiParse)
	if json.VerifyStringIsJSON(raw) && gjson.Valid(raw) {
		for _, item := range parseJiexiWithJSON(raw) {
			if len(item.URL) >= 1 {
				m[item.URL] = item
			}
		}
	} else {
		for _, item := range parseJiexiWithLines(raw) {
			if len(item.URL) >= 1 {
				m[item.URL] = item
			}
		}
	}
	var result = make([]JiexiParse, 0)
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func ParseMaccms(raw string) []MacCMSParse {
	var result = make([]MacCMSParse, 0)
	if json.VerifyStringIsJSON(raw) && gjson.Valid(raw) {
		// TODO: impl json
	} else {
		// TODO: impl lines
	}
	return result
}
