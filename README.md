# CS465-Lab1
AES Encryption Cipher Lab

## Setup and Operation
1. If you do not yet have GoLang installed follow the instructions [here](https://golang.org/doc/install) to set up your GoLang environment
1. run `go get github.com/dcorey28/CS465-Lab1` to dowload the repository to your go workspace
1. run `cd ${GOPATH}/src/github.com/dcorey28/CS465-Lab1` to go to the project
1. `go run main.go` will run the program and produce debugging output matching that found in FIPS 197 Appendix C

## Unit Tests
Unit tests are included in the project and can be run as described below:
1. `cd ${GOPATH}/src/github.com/dcorey28/CS465-Lab1/aes`
1. `go test`

### Statement of resources and authenticity
I only used the resources prescribed in the lab outline found [here](https://cs465.internet.byu.edu/fall-2018/projects/project1) and used no other resources related to AES and looked at no other source code. This work represents my (David Corey's) own work and should not be copied or used for any purpose.

### Appendix C Test Cases
My code passes all test cases in Appendix C of FIPS 197 and those test cases are included in my source code under `aes/aes_test.go` and can run and verified independently if desired. Simply follow the instructions above to run the unit tests.
