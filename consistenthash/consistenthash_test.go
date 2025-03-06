package consistenthash

import (
	"strconv"
	"testing"
)

func Test_consistenthash(t *testing.T) {
	m := NewMap(3, func(b []byte) uint32 {
		n, _ := strconv.Atoi(string(b))
		return uint32(n)
	})
	//02:2,12:2,22:2
	//04:4,14:4;24:4
	//06:6;16:6;26:6
	//2 4 6 12 14 16 22 24 26
	m.Add("2", "4", "6")
	testcases := map[string]string{
		"25": "6",
		"37": "2",
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testcases {
		if m.Get(k) != v {
			t.Errorf("get %s expect %s but %s", k, v, m.Get(k))
		}
	}
	//8:8,18:8;28:8
	m.Add("8")
	testcases["27"] = "8"
	for k, v := range testcases {
		if m.Get(k) != v {
			t.Errorf("get %s expect %s but %s", k, v, m.Get(k))
		}
	}
}
