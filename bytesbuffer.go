/*
@Time : 2021/6/23 下午5:36
@Author : MuYiMing
@File : bytebuffer
@Software: GoLand
*/
package bytebufferpool

import "io"

//用于读写的[]byte，可以理解位用于缓存的空间
type ByteBuffer struct {
	Buf []byte
}

//生成一个ByteBuffer数据结构
func NewByteBuffer(size uint64) *ByteBuffer {
	return &ByteBuffer{
		Buf: make([]byte, 0, size),
	}
}

//Len
func (b *ByteBuffer) Len() int {
	return len(b.Buf)
}

//read from io.reader
func (b *ByteBuffer) ReadFrom(r io.Reader) (int, error) {
	p := b.Buf
	bStart := len(p)
	nMax := cap(p)
	n := bStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, 64) //无空间创建
	} else {
		p = p[:nMax] //有空则直接全部占用
	}

	for {
		if n == nMax {
			//开辟空间
			nMax *= 2
			newCap := make([]byte, nMax)
			copy(newCap, p)
			p = newCap
		}

		nn, err := r.Read(p[n:])
		n += nn //指针移动
		if err != nil {
			b.Buf = p[:n]
			n -= nn
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

//WriteTo
func (b *ByteBuffer) WriteTo(w io.Writer) (int, error) {
	return w.Write(b.Buf)
}

//ReSet buf is empty
func (b *ByteBuffer) ReSet() {
	b.Buf = b.Buf[0:]
}

//Byte return all byte
func (b *ByteBuffer) Bytes() []byte {
	return b.Buf
}

//String returns string representation of buf
func (b *ByteBuffer) String() string {
	return BytesToString(b.Buf)
}

//Write  implements io.Writer - it appends p to ByteBuffer.B
func (b *ByteBuffer) Write(p []byte) (int, error) {
	b.Buf = append(b.Buf, p...)
	return len(p), nil
}

//WriteByte  appends the byte c to the buffer.
func (b *ByteBuffer) WriteByte(c byte) error {
	b.Buf = append(b.Buf, c)
	return nil
}

// WriteString appends s to ByteBuffer.B.
func (b *ByteBuffer) WriteString(s string) (int, error) {
	b.Buf = append(b.Buf, s...)
	return len(s), nil
}

// Set sets ByteBuffer.B to p.
func (b *ByteBuffer) Set(p []byte) (int, error) {
	b.Buf = append(b.Buf[:0], p...)
	return len(p), nil
}

// SetString sets ByteBuffer.B to s.
func (b *ByteBuffer) SetString(s string) (int, error) {
	b.Buf = append(b.Buf[:0], s...)
	return len(s), nil
}
