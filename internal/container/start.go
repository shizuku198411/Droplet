package container

func StartContainer(opt StartOption) error {
	fifo := fifoPath(opt.ContainerId)

	// write fifo
	if err := writeFifo(fifo); err != nil {
		return err
	}

	// remove fifo
	if err := removeFifo(fifo); err != nil {
		return err
	}

	return nil
}
