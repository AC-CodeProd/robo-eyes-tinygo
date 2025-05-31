module github.com/AC-CodeProd/robo-eyes-tinygo/example

go 1.23.0

toolchain go1.23.4

require (
	github.com/AC-CodeProd/robo-eyes-tinygo v0.0.0
	tinygo.org/x/drivers v0.31.0
)

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect

replace github.com/AC-CodeProd/robo-eyes-tinygo => ../
