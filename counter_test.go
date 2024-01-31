/*
* @Author:  老杨
* @Email:   xcapp1314@gmail.com
* @Date:    2024/1/31 22:35:03 星期三
* @Explain: ...
 */

package counter

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	cc := NewCounter()
	cc.Add(1)
	fmt.Println(cc.Count(time.Second))
	fmt.Println(cc.Count(time.Second))
	fmt.Println(cc.Count(time.Second))
}
