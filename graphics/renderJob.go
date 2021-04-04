package graphics

var (
	renderJobs = make(chan *RenderJob, 100)
)

// Async opengl sounds likes a great idea :)

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

func AddJobBlock(ro RenderObject, block func(r RenderObject)) bool {
	jobWrapper := func(job *RenderJob) []interface{} {
		block((job.params[0]).(RenderObject))
		return nil
	}
	return AddJob(&RenderJob{
		jobFunc: jobWrapper,
		params:  []interface{}{ro},
	})
}

func CreateDefaultRenderObjectJob(ro *DefaultRenderObject, texture string, elements int, callback func(...interface{})) {
	AddJob(&RenderJob{
		callback: callback,
		jobFunc:  callCreateDefaultRenderObject,
		params:   []interface{}{ro, texture, elements},
	})
}

func CreateBaseRenderObjectJob(ro *BaseRenderObject, texture string, elements int, callback func(...interface{})) {
	AddJob(&RenderJob{
		callback: callback,
		jobFunc:  callCreateBaseRenderObject,
		params:   []interface{}{ro, texture, elements},
	})
}

func callCreateDefaultRenderObject(job *RenderJob) []interface{} {
	ro := CreateDefaultRenderObject((job.params[1]).(string), (job.params[2]).(int))
	*(job.params[0]).(*DefaultRenderObject) = *ro
	ro.async = true
	return []interface{}{ro}
}

func callCreateBaseRenderObject(job *RenderJob) []interface{} {
	ro := CreateBaseRenderObject((job.params[1]).(string), (job.params[2]).(int))
	*(job.params[0]).(*BaseRenderObject) = *ro
	ro.async = true
	return []interface{}{ro}
}

func UpdateBuffersJob(ro RenderObject, callback func(...interface{})) {
	AddJob(&RenderJob{
		callback: callback,
		jobFunc:  callUpdateBuffers,
		params:   []interface{}{ro},
	})
}

func callUpdateBuffers(job *RenderJob) []interface{} {
	ro := job.params[0].(RenderObject)
	ro.UpdateBuffers()
	return nil
}
