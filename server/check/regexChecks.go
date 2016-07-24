package check

func CheckDBPrefix(s string) bool {
	return Validate("^[A-Za-z0-9]+$",s)
}

func CheckDomainName(s string) bool {
	return Validate("^(["+URL+"]+[\\.\\/])+["+URL+"]+$",s)
}

func CheckLogin(s string) bool {
	return Validate("^["+CHAR+"]["+CHAR+NUM+SPECIAL+"]*$", s)
}
