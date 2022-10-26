package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := wrapWithDone(in, done)

	for _, stage := range stages {
		if stage == nil {
			continue
		}

		out = stage(wrapWithDone(out, done))
	}
	return out
}

func wrapWithDone(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
