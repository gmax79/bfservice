package buckets

// AttemptsCounter - struct, which counts login attempts
type AttemptsCounter struct {
}

// CreateCounter - create main object
func CreateCounter() *AttemptsCounter {
	var c AttemptsCounter
	return &c
}

// CheckAndCount - main function to count attempts and collect it in buckets
func (fw *AttemptsCounter) CheckAndCount(login, password, hostip string) (bool, string, error) {
	return true, "", nil
}
