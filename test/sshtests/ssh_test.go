package tests
// Run source ../setenv.sh before executing this test

import (
    "testing"
    "github.com/rraks/remocc/pkg/controller"
)

func TestSSH(t *testing.T) {
    controller.AddKey("ssh-rsa rakshit@rbc")
}
