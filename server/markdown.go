package main

import (
	"regexp"
	"github.com/russross/blackfriday"
)

func MDtoHTML (mdstring string) string {
	md := []byte(imgReplace(mdstring))
	return string(blackfriday.MarkdownCommon(md))
}

func imgReplace(s string) string {
	re := regexp.MustCompile("<<(["+CHAR+NUM+SPECFILE+"]*)>><<<(left|right|center|floatleft|floatright)>>>")
	return re.ReplaceAllString(s, "<img src="+OP+"./static/$1"+OP+" class="+OP+"$2"+OP+">")
}

func imgRemove(s string) string {
	re := regexp.MustCompile("<<(["+CHAR+NUM+SPECFILE+"]*)>><<<(left|right|center|floatleft|floatright)>>>")
	return re.ReplaceAllString(s, "")
}
