package ziface

type IConnManager interface {
	Add(conn IConnection) error             // add connection
	Remove(conn IConnection)                // remove connection
	Get(connID uint32) (IConnection, error) // get connection by connID
	Len() int                               // get number of connections
	ClearConn()                             // clear all connections
}
