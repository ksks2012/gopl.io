# !/bin/bash
go run ./ch3/surfacep/main.go > sin.svg
go run ./ch3/surfacep/main.go -type=eggbox > eggbox.svg
go run ./ch3/surfacep/main.go -type=moguls > moguls.svg
go run ./ch3/surfacep/main.go -type=saddle > saddle.svg