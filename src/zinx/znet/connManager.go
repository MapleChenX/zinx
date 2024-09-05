package znet

import (
	"errors"
	"fmt"
	"go_code/src/zinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // connection map
	connLock    sync.RWMutex                  // lock
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) error {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if _, exists := cm.connections[conn.GetConnID()]; exists {
		return fmt.Errorf("connection with ID %d already exists", conn.GetConnID())
	}

	cm.connections[conn.GetConnID()] = conn
	return nil
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
}
