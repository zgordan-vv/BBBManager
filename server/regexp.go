package main

const (
	SPECFILE string ="-_,.@#$%^!*()+="
	OP string = string('"')
	AP string = "'"
	CHAR string = "A-Za-zа-яА-ЯёЁ"
	NUM string = "0-9"
	SPECIAL string = " -=.!@#$%*()_+|':;,?`~"
	URL string = "A_Za-z0-9-!*()_'~"
)

var (
	LOGIN = "^["+CHAR+"]["+CHAR+NUM+SPECIAL+"]*$"
	DESC = "^["+CHAR+NUM+SPECIAL+"]*$"
	D = "\\d"
	P = "\\."
	SEG = D+"?"+D+"?"+D
	IP = "^"+SEG+P+SEG+P+SEG+P+SEG+"$"
	DIGITS = "^"+D+"+$"
	CHARNUM = "^["+CHAR+NUM+"]+$"
	DOMAIN = "^(["+URL+"]+"+P+")*"+"["+URL+"]+$"
)
