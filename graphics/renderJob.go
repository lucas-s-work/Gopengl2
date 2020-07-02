package graphics

// Job Types
const (
	CreateRenderObjectJobType = iota
)

var (
	renderJobs = make(chan *RenderJob)
)

type RenderJob struct {
	callback func(interface{})
	jobFunc  func(*RenderJob) interface{}
	params   []interface{}
}

func (j *RenderJob) execute() {
	j.callback(j.jobFunc(j))
}

func PerformJobs() {
	select {
	case job := <-jobs:
		job.execute()
	}
}

func CreateRenderObjectJob(callback func(interface{}), texture string, elements int) {
	renderJobs <- &RenderJob{
		callback: callback,
		jobFunc:  callCreateRenderObject,
		params:   []interface{}{texture, elements},
	}
}

func callCreateRenderObject(job *RenderJob) interface{} {
	return CreateRenderObject((job.params[0]).(string), (job.params[1]).(int))
}
