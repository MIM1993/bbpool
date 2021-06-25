/*
@Time : 2021/6/23 下午5:36
@Author : MuYiMing
@File : poll
@Software: GoLand
*/
package bytebufferpool

import (
	"sort"
	"sync"
	"sync/atomic"
)

const (
	minBitSize = 6 // 2**6=64 is cpu cache line size
	steps      = 20

	minSize = 1 << minBitSize               //pool min size level
	maxSize = 1 << (minBitSize + steps - 1) //pool max size level

	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

//Pool represents byte buffer pool.
//根据缓存容量的大小将资源分配到不同的级别的pool中，从pool获取缓存空间时根据需要尺寸
//的大小分配对应level的缓存资源
type Pool struct {
	//保存每个level的资源数量
	calls [steps]uint64

	//一把🔒
	calibrating uint64

	//每个level pool的资源默认数量
	defaultSize uint64
	//每个level pool的资源最大数量
	maxSize uint64

	pool sync.Pool
}

//Get  ByteBuffer 从pool中获取资源
func (p *Pool) Get() *ByteBuffer {
	v := p.pool.Get()
	if v != nil {
		return v.(*ByteBuffer)
	}

	return NewByteBuffer(atomic.LoadUint64(&p.defaultSize))
}

//put  ByteBuffer 将使用后的资源返还给对应level的pool
func (p *Pool) Put(bb *ByteBuffer) {
	idx := index(len(bb.Buf))

	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold {
		//单一level资源数量超过界限，进行资源整理 calibrate
	}

	maxSize := atomic.LoadUint64(&p.maxSize)
	if maxSize == 0 || cap(bb.Buf) < int(maxSize) {
		bb.ReSet()
		p.pool.Put(bb)
	}

}

//calibrate
func (p *Pool) calibrate() {
	//加锁
	if atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}

	a := make(callSizes, steps)
	var callsSum uint64
	for i, _ := range p.calls {
		call := atomic.SwapUint64(&p.calls[i], 0)
		a = append(a, callSize{
			calls: call,
			size:  maxSize << i,
		})
		callsSum += call
	}

	sort.Sort(a)

	defaultSize := a[0].size
	maxSize := defaultSize

	maxSum := uint64(float64(callsSum) * maxPercentile)
	callsSum = 0
	for i := 0; i < steps; i++ {
		if callsSum > maxSum {
			break
		}
		callsSum += a[i].calls
		size := a[i].size
		if size > maxSize {
			maxSize = size
		}
	}

	atomic.StoreUint64(&p.defaultSize,defaultSize)
	atomic.StoreUint64(&p.maxSize,maxSize)

	//解锁
	atomic.StoreUint64(&p.calibrating,1)
}

type callSize struct {
	calls uint64
	size  uint64
}

type callSizes []callSize

func (ci callSizes) Len() int {
	return len(ci)
}

func (ci callSizes) Less(i, j int) bool {
	return ci[i].calls > ci[j].calls
}

func (ci callSizes) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}

//分级索引
func index(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		idx = steps - 1
	}
	return idx
}
