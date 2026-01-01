package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerRun_Success(t *testing.T) {
	dummySpecLoader := &dummyFileSpecLoader{}
	dummyFifoHandler := &dummyFifoHandler{}
	dummyCommandFactory := &dummyCommandFactory{
		cmd: &dummyCmd{
			pid: 11111,
		},
	}
	dummyContainerStart := &ContainerStart{
		fifoHandler: dummyFifoHandler,
	}
	dummyContainerRun := &ContainerRun{
		specLoader:     dummySpecLoader,
		fifoCreator:    dummyFifoHandler,
		commandFactory: dummyCommandFactory,
		containerStart: dummyContainerStart,
	}

	result := dummyContainerRun.Run(RunOption{ContainerId: "123456"})

	// assert
	assert.Equal(t, nil, result)
}
