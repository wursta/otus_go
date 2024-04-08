package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for i := 0; i < len(stages); i++ {
		out = runStage(stages[i], out, done)
	}
	return out
}

func runStage(stage Stage, in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			default:
				for v := range stage(in) {
					select {
					case <-done:
						return
					default:
						out <- v
					}
				}
				return
			}
		}
	}()

	return out
}
