package pool

import (
	"time"
	"sync"
	"container/list"
	"errors"
)

type Conn interface {
	Close()
}

type idleConn struct {
	c Conn
	t time.Time     // 放入链表中的连接对象对应的时间
}

type ConnPool struct {
	// 创建连接对象的函数
	Dial         func() (Conn, error)

	// 测试空闲连接对象健康情况
	TestOnBorrow func(c Conn, t time.Time) error

	MaxActive   int             // 允许最大连接数
	MaxIdle     int             // 允许最大空闲连接数
	IdleTimeout time.Duration   // 空闲连接存活时间
	Wait        bool

	mu     sync.Mutex
	cond   *sync.Cond  // 条件变量
	closed bool        // 连接池关闭标志
	active int         // 活动的连接数
	idle   list.List   // 存放空闲连接的链表，空闲连接对象类型为idleConn
}

func (p *ConnPool) Get() (Conn, error) {
	p.mu.Lock()

	// 从空闲连接链表中清除过期的空闲连接
	if timeout := p.IdleTimeout; timeout > 0 {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Back()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			if ic.t.Add(timeout).After(time.Now()) {
				break
			}
			p.idle.Remove(e)
			p.release()
			p.mu.Unlock()
			ic.c.Close()
			p.mu.Lock()
		}
	}

	for {
		// 从空闲连接链表中获取可用的连接
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Front()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			p.idle.Remove(e)
			test := p.TestOnBorrow
			p.mu.Unlock()
			if test == nil || test(ic.c, ic.t) == nil {
				return ic.c, nil
			}
			ic.c.Close()
			p.mu.Lock()
			p.release()
		}

		if p.closed {
			p.mu.Unlock()
			return nil, errors.New("get on closed pool")
		}

		// 获取新的连接对象
		if p.MaxActive == 0 || p.active < p.MaxActive {
			dial := p.Dial
			p.active += 1
			p.mu.Unlock()
			c, err := dial()
			if err != nil {
				p.mu.Lock()
				p.release()
				p.mu.Unlock()
				c = nil
			}
			return c, err
		}

		if !p.Wait {
			p.mu.Unlock()
			return nil, errors.New("connection pool exhausted")
		}

		if p.cond == nil {
			p.cond = sync.NewCond(&p.mu)
		}
		p.cond.Wait()  // 等待通知信号
	}
}

func (p *ConnPool) Put(c Conn) {
	p.mu.Lock()
	if !p.closed {  // 当连接池关闭时，跳过对连接池相关操作
		p.idle.PushFront(idleConn{c: c, t: time.Now()})
		if p.idle.Len() > p.MaxIdle {
			oldConn := p.idle.Remove(p.idle.Back()).(idleConn)
			oldConn.c.Close()      // 当空闲连接数满时，删除旧的连接对象
		}
	}

	if p.cond != nil {
		p.cond.Signal()
	}
	p.mu.Unlock()
}

func (p *ConnPool) Release(conn Conn) {
	p.active--
	//conn.Close()
}

func (p *ConnPool) Close() {
	p.mu.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.active -= idle.Len()
	if p.cond != nil {
		p.cond.Broadcast() // 通知所有等待的get()函数已经关闭了连接池信号
	}
	p.mu.Unlock()

	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(idleConn).c.Close()   // 释放所有的空闲连接对象
	}
}

func (p *ConnPool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

func (p *ConnPool) release() {
	p.active -= 1
	if p.cond != nil {
		p.cond.Signal()
	}
}