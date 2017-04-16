package main

type SearchConst struct {
	In string
	Out string
}

func (SearchConst) Search() error {
	return nil
}
