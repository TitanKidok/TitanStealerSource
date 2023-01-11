package main

import "strings"

var (
	m = map[string]string{}
)

func MagicEncode(str string) string {
	setmap()
	liststr := []string{}
	for _, symbol := range str {
		if m[string(symbol)] != "" {
			liststr = append(liststr, m[string(symbol)])
		}
	}
	return Join(liststr, " ")
}

func MagicDecode(str string) string {
	setmap()
	decoded := []string{}
	strl := strings.Split(str, " ")
	for _, s := range strl {
		if m[s] != "" {
			decoded = append(decoded, m[s])
		}
	}
	return Join(decoded, "")
}

func setmap() {
	m["0"] = "0110"
	m["1"] = "0122"
	m["2"] = "1012"
	m["3"] = "1011"
	m["4"] = "2022"
	m["5"] = "2143"
	m["6"] = "2134"
	m["7"] = "5424"
	m["8"] = "1353"
	m["9"] = "4524"
	m["."] = "1354"
	m[":"] = "111"
	m["0110"] = "0"
	m["0122"] = "1"
	m["1012"] = "2"
	m["1011"] = "3"
	m["2022"] = "4"
	m["2143"] = "5"
	m["2134"] = "6"
	m["5424"] = "7"
	m["1353"] = "8"
	m["4524"] = "9"
	m["1354"] = "."
	m["111"] = ":"
}

func Join(s []string, sep string) (sss string) {
	for _, ss := range s {
		sss += ss
		sss += sep
	}
	return
}
