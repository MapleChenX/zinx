package znet

import "go_code/src/zinx/ziface"

// 基类，根据需要进行重写
type BaseRouter struct{}

// 为什么方法全部为空？因为有时候我们不需要实现PreHandle或者PostHandle，但是抽象需要被全部实现才能用，所以我们实现，但全是空方法，这样就规避了这个问题

// 处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// 处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
