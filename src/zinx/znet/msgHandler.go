package znet

import "go_code/src/zinx/ziface"

type MsgHandler struct {
	CmdsHandler map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		CmdsHandler: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.CmdsHandler[request.GetMsgID()]
	if !ok {
		return
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
