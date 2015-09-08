package gopolygonjudger

import (
	"testing"
)

func TestNewFail(t *testing.T) {
	judger, err := NewJudger("./testdata0")
	if err != nil {
		t.Error("init failed, err:%v", err)
	}

	//judger.ToStdout()
	// city 2
	point1 := Point{24.152680, 115.107652}
	// city 2 miss
	point2 := Point{24.150305, 114.528813}

	// city 1
	point3 := Point{34.512752, 117.059471}
	// city 1 miss
	point4 := Point{34.522305, 117.215448}
	// city 1 hit
	point5 := Point{34.157152, 117.203769}
	//cityi 1 hit
	point6 := Point{34.237524, 117.153766}

	var id int
	id = judger.FindCityId(point1)
	if id != 1 {
		t.Errorf("1 find city failed, point:%v, ret:%d, ex:1", point1, id)
	}

	id = judger.FindCityId(point2)
	if id != -1 {
		t.Fatalf("2 find city faild, point:%v, ret:%d, ex:-1", point2, id)
	}

	id = judger.FindCityId(point3)
	if id != 0 {
		t.Errorf("3 find city failed, point:%v, ret:%d, ex:0", point3, id)
	}

	id = judger.FindCityId(point4)
	if id != -1 {
		t.Errorf("4 find city failed, point:%v, ret:%d, ex:-1", point4, id)
	}

	id = judger.FindCityId(point5)
	if id != 0 {
		t.Errorf("5 find city failed, point:%v, ret:%d, ex:0", point5, id)
	}

	id = judger.FindCityId(point6)
	if id != 0 {
		t.Errorf("6 find city failed, point:%v, ret:%d, ex:0", point6, id)
	}
}
