package osType

type OsType interface {
	Commands() []string
	InstallDocker() []string
}
