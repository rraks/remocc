package test

// Assumes environment is set from source ../setenv.sh 
import (
    "github.com/rraks/remocc/pkg/scheduler"
    "log"
    "testing"
    "time"
)

func hellow(s string, i int) {
    log.Println("Hellow\t", s, "\tNum\t", i)
}

func TestSched(t *testing.T) {
    sch := new(scheduler.Sched)
    sch.InitScheduler(time.Second*2, hellow, "world", 2)
    sch.Start()
    if (<-sch.C) {
        log.Println("Done")
    }
}
