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
)

func TestName(t *testing.T) {
	Add("test", 1)
	fmt.Println(Count("test", 1))
}
