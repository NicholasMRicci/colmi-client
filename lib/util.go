package lib

import "log"

func Must(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func Must1[T any](v T, err error) T {
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return v
}
