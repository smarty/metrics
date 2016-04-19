package metrics

import "sync"

type SynchronizedHistogram struct {
	inner     Histogram
	readLock  sync.Locker
	writeLock sync.Locker
}

func NewSynchronizedHistogram(inner Histogram, readLock, writeLock sync.Locker) *SynchronizedHistogram {
	return &SynchronizedHistogram{inner: inner, readLock: readLock, writeLock: writeLock}
}

func (this *SynchronizedHistogram) RecordValue(v int64) error {
	this.writeLock.Lock()
	err := this.inner.RecordValue(v)
	this.writeLock.Unlock()
	return err
}

func (this *SynchronizedHistogram) Min() int64 {
	this.readLock.Lock()
	min := this.inner.Min()
	this.readLock.Unlock()
	return min
}

func (this *SynchronizedHistogram) Max() int64 {
	this.readLock.Lock()
	max := this.inner.Max()
	this.readLock.Unlock()
	return max
}

func (this *SynchronizedHistogram) Mean() float64 {
	this.readLock.Lock()
	mean := this.inner.Mean()
	this.readLock.Unlock()
	return mean
}

func (this *SynchronizedHistogram) StdDev() float64 {
	this.readLock.Lock()
	deviation := this.inner.StdDev()
	this.readLock.Unlock()
	return deviation
}

func (this *SynchronizedHistogram) TotalCount() int64 {
	this.readLock.Lock()
	total := this.inner.TotalCount()
	this.readLock.Unlock()
	return total
}

func (this *SynchronizedHistogram) ValueAtQuantile(q float64) int64 {
	this.readLock.Lock()
	value := this.inner.ValueAtQuantile(q)
	this.readLock.Unlock()
	return value
}
