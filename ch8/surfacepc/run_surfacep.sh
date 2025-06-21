# !/bin/bash
go run ./ch8/surfacepc/main.go > sin.svg
go run ./ch8/surfacepc/main.go -type=eggbox > eggbox.svg
go run ./ch8/surfacepc/main.go -type=moguls > moguls.svg
go run ./ch8/surfacepc/main.go -type=saddle > saddle.svg