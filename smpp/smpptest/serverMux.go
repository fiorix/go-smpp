package main

import (
	"github.com/fiorix/go-smpp/smpp/pdu"
	"sync"
)

//The HandlerFunc type is an adapter to allow the use of ordinary functions as SMPP handlers.
//If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(pdu.Body, ResponseWriter)

//ServeSMPP is method for implementing Handler intreface.
func (hf HandlerFunc) ServeSMPP(p pdu.Body, rw ResponseWriter) {
	hf(p, rw)
}

//ServerMux is handler for smpp server. It route requets.
type ServerMux struct {
	mu sync.RWMutex
	m  map[pdu.ID]Handler
}

//Handle set Handler to id.
func (sm *ServerMux) Handle(id pdu.ID, hndl Handler) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[id] = hndl
}

//Handler returns Handler for id. If Handler not defined result is nil.
func (sm *ServerMux) Handler(id pdu.ID) Handler {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.m[id]
}

//ServeSMPP is method for implementing Handler intreface.
func (sm *ServerMux) ServeSMPP(p pdu.Body, rw ResponseWriter) {
	sm.mu.RLock()
	handler := sm.m[p.Header().ID]
	sm.mu.RUnlock()
	if handler == nil {
		return
	}
	handler.ServerSMPP(p, rw)
}
