package main

import "strconv"

func stringToValue(s string) any {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	} else if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}

func arrayToValues(a []string) any {
	var d []any
	for _, v := range a {
		d = append(d, stringToValue(v))
	}
	return d
}
