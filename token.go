package bseg

type Token struct {
	Name  string
	Count int
}

type Tokens []Token

func (ts Tokens) Len() int {
	return len(ts)
}

func (ts Tokens) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Tokens) Less(i, j int) bool {
	if ts[i].Count != ts[j].Count {
		return ts[i].Count > ts[j].Count
	}
	return ts[i].Name < ts[j].Name
}
