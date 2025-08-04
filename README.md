# NJUPTautoAuthorize4Windows

Are you tired of having to log in to the internet every time you turn on your computer? This is an auto-login script for NJUPT written in Go.

## Build

1. install golang
2. fill in the 4 constant value according to your config
3. run `go mod tidy` to install dependencies
4. run `go build` to compile the executable file

## Use

1. press `Win+R` to open the **Run dialog box**
2. enter `shell:startup` to open startup folder
3. move the compiled binary to this folder