package main

var output map[bool][]byte = map[bool][]byte{true:[]byte{'t','r','u','e'},false:[]byte{'f','a','l','s','e'}}

var (
	PORT, DBPREFIX, DOMAINNAME string
	installed bool
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
