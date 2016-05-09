package main

func checkDBPrefix(s string) bool {
	return validate("^[A-Za-z0-9]+$",s)
}

func checkDomainName(s string) bool {
	return validate("^(["+URL+"]+[\\.\\/])+["+URL+"]+$",s)
}

func checkLogin(s string) bool {
	return validate("^["+CHAR+"]["+CHAR+NUM+SPECIAL+"]*$", s)
}

func checkPwd(username,pwd string) bool {
	user, ok := getUser(username)
	if ok {return user.Keyword == passEncrypt(pwd)} else {return false}
}
