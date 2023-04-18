package dao

// MockDAOCall holds data to allow mocking a DAO query method.
type MockDAOCall struct {
	ShouldCall bool
	Args       []interface{}
	Return     []interface{}
}

// Bind sets up assertions based on the data in the MockDAOCall.
func (c MockDAOCall) Bind(m *MockQueryProvider, name string) {
	if c.ShouldCall {
		m.On(name, c.Args...).Return(c.Return...)
	}
}
