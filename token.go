package bseg

type Token struct {
	name  string
	count int
}

type Tokens []Token

func (ts Tokens) Len() int {
	return len(ts)
}

func (ts Tokens) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Tokens) Less(i, j int) bool {
	return ts[i].count < ts[j].count
}
