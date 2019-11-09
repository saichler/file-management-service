package message_handlers

import (
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerHandlers struct {
	handlers map[string]IMessageHandler
}

func (mh *FileManagerHandlers) Init(service IService) {
	mh.handlers = make(map[string]IMessageHandler)
	mh.addHanlder(NewRlsMH(service))
	mh.addHanlder(NewCpMH(service))
}

func (mh *FileManagerHandlers) Handlers() []IMessageHandler {
	result := make([]IMessageHandler, 0)
	for _, h := range mh.handlers {
		result = append(result, h)
	}
	return result
}

func (mh *FileManagerHandlers) addHanlder(handler IMessageHandler) {
	mh.handlers[handler.Topic()] = handler
}

func (mh *FileManagerHandlers) Handler(topic string) IMessageHandler {
	return mh.handlers[topic]
}
