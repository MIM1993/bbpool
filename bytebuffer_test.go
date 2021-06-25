/*
@Time : 2021/6/24 上午9:40
@Author : MuYiMing
@File : bytebuffer_test
@Software: GoLand
*/
package bytebufferpool

import (
	"fmt"
	"testing"
)

func TestByteBuffer_Set(t *testing.T) {
	arr := make([]byte,5,10)

	arr[0]='a'
	arr[1]='b'
	arr[2]='c'

	fmt.Println(string(arr))
	fmt.Printf("%d  ==> %d\n",len(arr),cap(arr))

	p := []byte{'x','y'}

	arr = append(arr[:0],p...)

	fmt.Println(string(arr))
	fmt.Printf("%d  ==> %d\n",len(arr),cap(arr))




}
