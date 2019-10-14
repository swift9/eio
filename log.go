package eio

import "log"

type ILog interface {
	Info(args ...interface{})
	Infow(template string, args ...interface{})
	Error(args ...interface{})
	Errorw(template string, args ...interface{})
}

type SysLog struct {
}

func (l *SysLog) Info(args ...interface{}) {
	log.Println("INFO ", args)
}

func (l *SysLog) Infow(template string, args ...interface{}) {
	log.Println("INFO "+template, args)
}

func (l *SysLog) Error(args ...interface{}) {
	log.Println("ERROR ", args)
}

func (l *SysLog) Errorw(template string, args ...interface{}) {
	log.Println("ERROR "+template, args)
}
