package container

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignalMap(t *testing.T) {
	// == arrange ==

	// == act ==
	sigTerm := signalMap["TERM"]
	sigKill := signalMap["KILL"]
	sigInt := signalMap["INT"]
	sigHup := signalMap["HUP"]
	sigQuit := signalMap["QUIT"]
	sigUsr1 := signalMap["USR1"]
	sigUsr2 := signalMap["USR2"]
	sigStop := signalMap["STOP"]
	sigCont := signalMap["CONT"]

	// == assert ==
	assert.Equal(t, syscall.SIGTERM, sigTerm)
	assert.Equal(t, syscall.SIGKILL, sigKill)
	assert.Equal(t, syscall.SIGINT, sigInt)
	assert.Equal(t, syscall.SIGHUP, sigHup)
	assert.Equal(t, syscall.SIGQUIT, sigQuit)
	assert.Equal(t, syscall.SIGUSR1, sigUsr1)
	assert.Equal(t, syscall.SIGUSR2, sigUsr2)
	assert.Equal(t, syscall.SIGSTOP, sigStop)
	assert.Equal(t, syscall.SIGCONT, sigCont)
}
