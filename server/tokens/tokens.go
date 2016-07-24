package tokens

var (
	ControlToken map[string]string = make(map[string]string)
	SubmitToken map[string]string = make(map[string]string)
	SecToken map[string]string = make(map[string]string)
	secTokens map[string][][]byte = make(map[string][][]byte)
	dupTokens map[string][]byte = make(map[string][]byte)
	dupUsedTokens map[string][]byte = make(map[string][]byte)
)


func CheckSec(sectoken []byte, username string) bool {
	tokens, ok := secTokens[username]
	if !ok {return false} else {
		return true
		for i, token := range(tokens) {
			if string(token) == string(sectoken) {
				tokens = append(tokens[:i], tokens[i+1:]...)
				secTokens[username] = tokens
				return true
			}
		}
		return false
	}
}

func DupControl(username string) bool {
	if string(dupTokens[username]) == string(dupUsedTokens[username]) {return false} else {
		dupUsedTokens[username] = dupTokens[username]
		return true
	}
}
