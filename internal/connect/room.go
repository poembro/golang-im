package connect

import (
    "container/list"
    "golang-im/pkg/protocol"
    "sync"
    "time"
)

var RoomsManager sync.Map

// SubscribedRoom 订阅房间 (目前一个连接句柄无法同时订阅多个房间)
func SubscribedRoom(conn *Conn, roomId string) {
    if roomId == conn.RoomId {
        return //该连接句柄 已经订阅过该房间了 直接return
    }

    oldRoomId := conn.RoomId
    // 该连接曾经订阅过其他房间 则先取消
    if oldRoomId != "" {
        value, ok := RoomsManager.Load(oldRoomId)
        if !ok {
            return
        }
        room := value.(*Room)
        room.Unsubscribe(conn)

        if room.Conns.Front() == nil {
            RoomsManager.Delete(oldRoomId)
        }
        return
    }

    // 订阅 1.先用该roomid 创建一个room结构对象
    //      2.再将其放入全局map key:roomid val:room对象
    //      3.将当前用户的连接句柄 放入到 room对象下Conns链表中
    if roomId != "" {
        value, ok := RoomsManager.Load(roomId)
        var room *Room
        if !ok {
            room = NewRoom(roomId)
            RoomsManager.Store(roomId, room)
        } else {
            room = value.(*Room)
        }
        room.Subscribe(conn)
        return
    }
}

// PushRoom 从全局Map中 找到对应roomid对应的Room结构体对象, 该对象下有 所有用户的连接句柄
func PushRoom(roomId string, p *protocol.Proto) {
    //RoomsManager.Range(func(key, value interface{}) bool {
    //    fmt.Println("当前有房间 %+v ", value.(*Room))
    //    return true
    //})
    //fmt.Println("推送到房间  1111   ", roomId)
    value, ok := RoomsManager.Load(roomId)
    if !ok {
        return
    }
    value.(*Room).Push(p)
}

type Room struct {
    RoomId string     // 房间ID
    Conns  *list.List // 订阅房间消息的连接
    lock   sync.RWMutex
}

func NewRoom(roomId string) *Room {
    return &Room{
        RoomId: roomId,
        Conns:  list.New(), //初始化1个空链表
    }
}

// Subscribe 订阅房间
func (r *Room) Subscribe(conn *Conn) {
    r.lock.Lock()
    defer r.lock.Unlock()
    // 将conn指针对象追加到链表,返回链表元素指针 该返回值可以用来删除链表中指定元素
    conn.Element = r.Conns.PushBack(conn)
    conn.RoomId = r.RoomId
}

// Unsubscribe 取消订阅
func (r *Room) Unsubscribe(conn *Conn) {
    r.lock.Lock()
    defer r.lock.Unlock()

    r.Conns.Remove(conn.Element)
    conn.Element = nil
    conn.RoomId = ""
}

// Push 推送消息到房间 (推送消息到该房间下的所有句柄)
func (r *Room) Push(p *protocol.Proto) {
    r.lock.RLock()
    defer r.lock.RUnlock()
 
    timeout := 2 * time.Minute // 1 分钟
    timeoutConns := make([]*Conn, 0) // 超时的连接 
    element := r.Conns.Front()
    for {
        c := element.Value.(*Conn)
        //fmt.Printf("下发消息 %+v", c)
        if time.Now().Sub(c.LastHeartbeatTime) > timeout {
            timeoutConns = append(timeoutConns, c) 
        } else {
            c.Send(p)
        }
        
        element = element.Next()
        if element == nil {
            break
        }
    }
 
    for i := range timeoutConns {
        //fmt.Printf("关闭连接111 %+v", timeoutConns[i])
        timeoutConns[i].Close()  // 超时情况下，这里要从链表中删除， 由Close 方法中另起一个groutine去做
    }
}
