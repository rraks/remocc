package scheduler

import (
    "time"
    "reflect"
)


type Scheduler interface {
    InitScheduler(fn func())
    Start() (err error)
    Stop() (err error)
}

type Sched struct {
    Duration time.Duration
    Timer *time.Timer
    Task interface{}
    Args []interface{}
    C chan bool
}

func (sch *Sched) InitScheduler(duration time.Duration, fn interface{}, args ...interface{}) {
    sch.Duration = duration
    sch.Task = fn
    sch.Args = args
    sch.C = make(chan bool)
}


func  (sch *Sched) CallFunc() {
    vf := reflect.ValueOf(sch.Task)
    vxs := []reflect.Value{}
    for i:=0; i<len(sch.Args); i++ {
        vxs = append(vxs, reflect.ValueOf(sch.Args[i]))
    }
    vf.Call(vxs)
}

func (sch *Sched) Start() {
    sch.Timer = time.NewTimer(sch.Duration)
    go func() {
        <-sch.Timer.C
        sch.CallFunc()
        sch.C<-true
    }()
}

func (sch *Sched) Stop() {
    sch.Timer.Stop()
    sch.C<-true
}
