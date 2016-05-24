package main

import "testing"

//These tests currently just call the functions as if they are being run by the program itself.
//This is a basic check, but to properly test the program some refactoring would be required
//to allow the functions to take args representing behaviour such as reading from folders.
//Then I could mock this behaviour in these tests, perhaps with temp dirs and files (reading and
//writing to the GOROOT path for tests is not very nice).
func TestlistGos(t *testing.T) {
	listGos()
}

func Testgoroot(t *testing.T) {
	goroot("")
}
