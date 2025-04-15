package domainapp

type IDomainApp interface {
	GetArgs() []string
	GetInputFileName() (string, error)
	Run() (err error)
}
