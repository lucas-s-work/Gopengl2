package graphics

var (
	renderJobs = make(chan *RenderJob, 100)
)

// Async opengl sounds likes a great idea...

type RenderJob struct {
	callback func(...interface{})
	jobFunc  func(*RenderJob) []interface{}
	params   []interface{}
}

func (j *RenderJob) execute() {
	r := j.jobFunc(j)
	if j.callback != nil {
		j.callback(r...)
	}
}

func performJobs() {
	select {
	case job := <-renderJobs:
		job.execute()
		return
	default:
		return
	}
}

// Non-blocking send, this is working on a single go routine with nested sends
// So this is expected to block sometimes.
func AddJob(job *RenderJob) bool {
	select {
	case renderJobs <- job:
		return true
	default:
		return false
	}
}

func (ro *RenderObject) AddJobBlock(block func(ro *RenderObject)) bool {
	jobWrapper := func(job *RenderJob) []interface{} {
		block((job.params[0]).(*RenderObject))
		return nil
	}
	return AddJob(&RenderJob{
		jobFunc: jobWrapper,
		params:  []interface{}{ro},
	})
}

func CreateRenderObjectJob(texture string, elements int, callback func(...interface{})) {
	AddJob(&RenderJob{
		callback: callback,
		jobFunc:  callCreateRenderObject,
		params:   []interface{}{texture, elements},
	})
}

func callCreateRenderObject(job *RenderJob) []interface{} {
	ro := CreateRenderObject((job.params[0]).(string), (job.params[1]).(int))
	ro.async = true
	return []interface{}{ro}
}

func (ro *RenderObject) UpdateBuffersJob(callback func(...interface{})) {
	AddJob(&RenderJob{
		callback: callback,
		jobFunc:  callUpdateBuffers,
		params:   []interface{}{ro},
	})
}

func callUpdateBuffers(job *RenderJob) []interface{} {
	ro := job.params[0].(*RenderObject)
	ro.vao.UpdateBuffers()
	return nil
}
