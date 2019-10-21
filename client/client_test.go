package client_test

import (
	"github.com/quicsocks/quicsocks/client"
	"testing"
)

func TestNewClient(t *testing.T) {
	err := client.NewClient("0.0.0.0:1080", "")
	t.Fatal(err)
}
