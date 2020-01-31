package buckets

// RatesLimits - maximums for rates, after their, check limited
type RatesLimits struct {
	Login    int
	Password int
	Host     int
}

// AttemptsCounter - struct, which counts login attempts
type AttemptsCounter struct {
	limits RatesLimits
}

// CreateCounter - create main object
func CreateCounter(l RatesLimits) *AttemptsCounter {
	var c AttemptsCounter
	c.limits = l
	return &c
}

// CheckAndCount - main function to count attempts and collect it in buckets
func (fw *AttemptsCounter) CheckAndCount(login, password, hostip string) (bool, string, error) {
	return true, "", nil
}

// Reset - reset login+host from counter buckets
func (fw *AttemptsCounter) Reset(login, hostip string) (bool, error) {
	return false, nil
}
