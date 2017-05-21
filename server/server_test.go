package main

import "testing"

func TestMain(m *testing.M) {
	// Listen on localhost on a random port
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Errorf("Failed to listen: %v\n", err)
	}
	os.Exit(m.Run())
}

func TestGetPinDirection(t *testing.T) {
	//Feature not supported in library
}
