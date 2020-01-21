package hostslist

type subnet struct {
	ip   uint32
	mask uint32
}

func (m *subnet) Parse(subnet string) error {
	return nil
}

func (m *subnet) Check(host ip) bool {
	return false
}
