package model

import (
	"bytes"
	"strconv"
	"sync"
	"time"
)

type SyncReport struct {
	exists     []*FileDescriptor
	copied     []*FileDescriptor
	errored    map[string][]*FileDescriptor
	sizeDiff   []*FileDescriptor
	hashDiff   []*FileDescriptor
	lastReport int64
	started    int64
	lock       *sync.Mutex
}

func NewSyncReport() *SyncReport {
	sr := &SyncReport{}
	sr.lastReport = time.Now().Unix()
	sr.started = sr.lastReport
	sr.exists = make([]*FileDescriptor, 0)
	sr.copied = make([]*FileDescriptor, 0)
	sr.errored = make(map[string][]*FileDescriptor, 0)
	sr.sizeDiff = make([]*FileDescriptor, 0)
	sr.hashDiff = make([]*FileDescriptor, 0)
	sr.lock = &sync.Mutex{}
	return sr
}

func (sr *SyncReport) AddExist(descriptor *FileDescriptor) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	sr.exists = append(sr.exists, descriptor)
}

func (sr *SyncReport) AddCopied(descriptor *FileDescriptor) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	sr.copied = append(sr.copied, descriptor)
}

func (sr *SyncReport) AddErrored(msg string, descriptor *FileDescriptor) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	_, ok := sr.errored[msg]
	if !ok {
		sr.errored[msg] = make([]*FileDescriptor, 0)
	}
	sr.errored[msg] = append(sr.errored[msg], descriptor)
}

func (sr *SyncReport) AddSizeDiff(descriptor *FileDescriptor) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	sr.sizeDiff = append(sr.sizeDiff, descriptor)
}

func (sr *SyncReport) AddHashDiff(descriptor *FileDescriptor) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	sr.hashDiff = append(sr.hashDiff, descriptor)
}

func (sr *SyncReport) Exists() []*FileDescriptor {
	return sr.clone(sr.exists)
}

func (sr *SyncReport) Copied() []*FileDescriptor {
	return sr.clone(sr.copied)
}

func (sr *SyncReport) SizeDiff() []*FileDescriptor {
	return sr.clone(sr.sizeDiff)
}

func (sr *SyncReport) HashDiff() []*FileDescriptor {
	return sr.clone(sr.hashDiff)
}

func (sr *SyncReport) clone(arr []*FileDescriptor) []*FileDescriptor {
	result := make([]*FileDescriptor, 0)
	sr.lock.Lock()
	defer sr.lock.Unlock()
	for _, d := range arr {
		result = append(result, d)
	}
	return result
}

func (sr *SyncReport) Report(ignoreTimeout bool) string {
	if time.Now().Unix()-sr.lastReport > 30 || ignoreTimeout {
		sr.lastReport = time.Now().Unix()
		buff := bytes.Buffer{}
		buff.WriteString("Sync Report (" + strconv.Itoa(int(time.Now().Unix()-sr.started)) + "):\n")
		sr.lock.Lock()
		defer sr.lock.Unlock()
		buff.WriteString("Exists: ")
		buff.WriteString(strconv.Itoa(len(sr.exists)))
		buff.WriteString("\n")
		buff.WriteString("Copied: ")
		buff.WriteString(strconv.Itoa(len(sr.copied)))
		buff.WriteString("\n")
		buff.WriteString("SizeDiff: ")
		buff.WriteString(strconv.Itoa(len(sr.sizeDiff)))
		buff.WriteString("\n")
		buff.WriteString("HashDiff: ")
		buff.WriteString(strconv.Itoa(len(sr.hashDiff)))
		buff.WriteString("\n")
		for msg, list := range sr.errored {
			buff.WriteString(msg)
			buff.WriteString(" Errored: ")
			buff.WriteString(strconv.Itoa(len(list)))
			buff.WriteString("\n")
		}
		return buff.String()
	}
	return ""
}
