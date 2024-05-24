/*
@author: sk
@date: 2023/2/25
*/
package test

import (
	"fmt"
	"strings"
	"testing"
)

func Test11(t *testing.T) {
	tag := "<object id=\"%d\" class=\"grid\" x=\"%d\" y=\"%d\">\n   <point/>\n  </object>"
	res := strings.Builder{}
	for i := 0; i < 20; i++ {
		for j := 0; j < 9; j++ {
			res.WriteString(fmt.Sprintf(tag, i+j*20+1, i*64, j*64+32))
		}
	}
	fmt.Println(res.String())
}
