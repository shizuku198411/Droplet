package container

type mockAppArmorManager struct {
	aaprofileName string
	aaprofileErr  error

	aaprofileOnexecName string
	aaprofileOnexecErr  error
}

func (m *mockAppArmorManager) ApplyAAProfile(profile string) error {
	m.aaprofileName = profile
	return m.aaprofileErr
}

func (m *mockAppArmorManager) ApplyAAProfileOnExec(profile string) error {
	m.aaprofileOnexecName = profile
	return m.aaprofileOnexecErr
}
