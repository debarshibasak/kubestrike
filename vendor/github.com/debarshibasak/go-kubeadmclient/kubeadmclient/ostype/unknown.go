package osType

type Unknown struct {
}

func (u *Unknown) Commands() []string {
	panic("unknown os type detected")
	return []string{}
}

func (u *Unknown) InstallDocker() []string {
	panic("unknown os type detected while installing docker")
	return []string{}
}
