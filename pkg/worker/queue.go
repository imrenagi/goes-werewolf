package worker

var defaultWorkQueue = make(chan Job, MAX_QUEUE)
