package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	StartOneWorker(workerID int, taskQueue chan IRequest)
	// 发送给工作池的消息队列处理
	SendMsgToTaskQueue(request IRequest)
}
