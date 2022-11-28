# Actionsfmt

This is a small program that reads Go Test JSON output from its Stdin and formats it according to GitHub Actions rule. 

## Install

go install marwan.io/actionsfmt@latest

### Usage

In an actions workflow, run `go test -json | actionsfmt`
