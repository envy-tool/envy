package runs

// Runner is the main type in this package, used to run either individual
// commands or a persistent background agent.
type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}
