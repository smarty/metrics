package metrics

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SynchronizedHistogramFixture struct {
	*gunit.Fixture

	outer     *SynchronizedHistogram
	inner     *FakeHistogram
	readLock  *FakeLock
	writeLock *FakeLock
}

func (this *SynchronizedHistogramFixture) Setup() {
	this.readLock = &FakeLock{}
	this.writeLock = &FakeLock{}
	this.inner = &FakeHistogram{}
	this.outer = NewSynchronizedHistogram(this.inner, this.readLock, this.writeLock)
}

func (this *SynchronizedHistogramFixture) TestReadBehavior_Max_IsProtectedWithTheReadLock() {
	this.outer.Max()
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "Max")
}
func (this *SynchronizedHistogramFixture) TestReadBehavior_Min_IsProtectedWithTheReadLock() {
	this.outer.Min()
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "Min")
}
func (this *SynchronizedHistogramFixture) TestReadBehavior_Mean_IsProtectedWithTheReadLock() {
	this.outer.Mean()
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "Mean")
}
func (this *SynchronizedHistogramFixture) TestReadBehavior_StandardDeviation_IsProtectedWithTheReadLock() {
	this.outer.StdDev()
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "StdDev")
}
func (this *SynchronizedHistogramFixture) TestReadBehavior_TotalCount_IsProtectedWithTheReadLock() {
	this.outer.TotalCount()
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "TotalCount")
}
func (this *SynchronizedHistogramFixture) TestReadBehavior_ValueAtQuantile_IsProtectedWithTheReadLock() {
	value := this.outer.ValueAtQuantile(99.9)
	this.So([]time.Time{this.readLock.locked, this.inner.instant, this.readLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "ValueAtQuantile")
	this.So(value, should.Equal, 99)
}
func (this *SynchronizedHistogramFixture) TestWriteBehavior_RecordValue_IsProtectedWithTheWriteLock() {
	err := this.outer.RecordValue(42)
	this.So([]time.Time{this.writeLock.locked, this.inner.instant, this.writeLock.unlocked}, should.BeChronological)
	this.So(this.inner.invoked, should.Equal, "RecordValue")
	this.So(err, should.Resemble, errors.New("42"))
}

////////////////////////////////////////////////////////////////////////////

type FakeHistogram struct {
	invoked string
	instant time.Time
}

func (this *FakeHistogram) recordFirstInvocation(name string) {
	if this.invoked != "" {
		panic("Unexpected invocation.")
	}
	this.instant = time.Now()
	this.invoked = name
}
func (this *FakeHistogram) RecordValue(v int64) error {
	this.recordFirstInvocation("RecordValue")
	return fmt.Errorf("%d", v)
}
func (this *FakeHistogram) Min() int64 {
	this.recordFirstInvocation("Min")
	return 12345
}
func (this *FakeHistogram) Max() int64 {
	this.recordFirstInvocation("Max")
	return 54321
}
func (this *FakeHistogram) Mean() float64 {
	this.recordFirstInvocation("Mean")
	return 123.45
}
func (this *FakeHistogram) StdDev() float64 {
	this.recordFirstInvocation("StdDev")
	return 54.321
}
func (this *FakeHistogram) TotalCount() int64 {
	this.recordFirstInvocation("TotalCount")
	return 99999
}
func (this *FakeHistogram) ValueAtQuantile(q float64) int64 {
	this.recordFirstInvocation("ValueAtQuantile")
	return int64(q)
}

////////////////////////////////////////////////////////////////////////////

type FakeLock struct {
	locked   time.Time
	unlocked time.Time
}

func (this *FakeLock) Lock()   { this.locked = time.Now() }
func (this *FakeLock) Unlock() { this.unlocked = time.Now() }
