# Actionsfmt

This is a small program that reads Go Test JSON output from its Stdin and formats it according to GitHub Actions rule. 

## Installation and Usage

go install marwan.io/actionsfmt@latest
In an actions workflow:

1. `go install marwan.io/actionsfmt@latest` 
2. run `go test -json | actionsfmt`

Or you can do it all in one line: `go test -json | go run marwan.io/actionsfmt@latest` 
