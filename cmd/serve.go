package main

type Service interface {
	Serve(req []byte) ([]byte, error)
}
