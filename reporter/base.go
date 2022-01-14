package reporter

import "io"

type Reporter interface {
	Report(writer io.Writer)
}

type name struct {
}
