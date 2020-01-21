package hostslist

type subnet struct {
	address ip
	mask    ip
}

func (m *subnet) Parse(subnet string) error {
	return nil
}

func (m *subnet) Check(host ip) bool {
	return false
}
