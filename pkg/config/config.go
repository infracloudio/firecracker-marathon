package config

type Runtime struct {
	Environment  Environment
	FunctionUUID string
}

type RunConfig struct {
	Language   ProgrammingLanguage
	SocketPath string
	KernelPath string
	KernelOpts string
	RootDrive  string
}

type Network struct {
	TapDeviceName string
	MAC           string
	IP            string
}

type Environment struct {
	Language ProgrammingLanguage
	Version  string
}

type ProgrammingLanguage string

const (
	Go     ProgrammingLanguage = "go"
	Python ProgrammingLanguage = "python"
)
