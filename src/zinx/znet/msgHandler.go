package znet

import (
	"go_code/src/zinx/utils"
	"go_code/src/zinx/ziface"
)

type MsgHandler struct {
	CmdsHandler map[uint32]ziface.IRouter

	// 业务工作池的worker数量
	WorkerPoolSize uint32

	// 业务工作worker队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		CmdsHandler:    make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalVar.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalVar.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.CmdsHandler[request.GetMsgID()]
	if !ok {
		router, _ = mh.CmdsHandler[0]
	}

	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断是否存在
	if _, ok := mh.CmdsHandler[msgID]; ok {
		panic("repeat router")
	}
	mh.CmdsHandler[msgID] = router
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalVar.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 轮询负载均衡
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}
