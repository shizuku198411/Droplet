package container

type mockCotainerFifoHandler struct {
	// createFifo()
	createFifoCallFlag bool
	createFifoPath     string
	createFifoErr      error

	// removeFifo()
	removeFifoCallFlag bool
	removeFifoPath     string
	removeFifoErr      error

	// readFifo()
	readFifoCallFlag bool
	readFifoPath     string
	readFifoErr      error

	// writeFifo()
	writeFifoCallFlag bool
	writeFifoPath     string
	writeFifoErr      error
}

func (m *mockCotainerFifoHandler) createFifo(path string) error {
	m.createFifoCallFlag = true
	m.createFifoPath = path
	return m.createFifoErr
}

func (m *mockCotainerFifoHandler) removeFifo(path string) error {
	m.removeFifoCallFlag = true
	m.removeFifoPath = path
	return m.removeFifoErr
}

func (m *mockCotainerFifoHandler) readFifo(path string) error {
	m.readFifoCallFlag = true
	m.readFifoPath = path
	return m.readFifoErr
}

func (m *mockCotainerFifoHandler) writeFifo(path string) error {
	m.writeFifoCallFlag = true
	m.writeFifoPath = path
	return m.writeFifoErr
}
