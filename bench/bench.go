package main

import (
	"bufio"
	"fmt"
	"gopolygonjudger"
	"io"
	"os"
	"time"
)

func LoadTest(filename string) ([]*gopolygonjudger.Point, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var lat float64
	var lng float64
	var cid int
	var points []*gopolygonjudger.Point
	var results []int

	count, err := fmt.Fscanf(reader, "%f %f %d\n", &lat, &lng, &cid)
	for err == nil && count == 3 {
		points = append(points, &gopolygonjudger.Point{lat, lng})
		results = append(results, cid)
		count, err = fmt.Fscanf(reader, "%f %f %d\n", &lat, &lng, &cid)
	}
	if err != io.EOF {
		return nil, nil, err
	}
	return points, results, nil
}

func main() {
	arg_num := len(os.Args)
	if arg_num != 3 {
		fmt.Printf("%s city-info test-points\n", os.Args[0])
		return
	}

	var load_start int64
	var load_end int64
	var check_start int64
	var check_end int64

	var query []*gopolygonjudger.Point
	var result []int
	var queryCount int

	var right int
	var wrong int

	var err error

	// load testdata
	query, result, err = LoadTest(os.Args[2])
	if err != nil {
		fmt.Printf("Load TestData %s failed:%v", os.Args[2], err)
		return
	}
	queryCount = len(query)

	defer func(load_start *int64, load_end *int64, check_start *int64, check_end *int64, right *int, wrong *int, queryCount int) {
		load := (float64)((*load_end)-(*load_start)) / (float64)(time.Millisecond)
		check := (float64)((*check_end)-(*check_start)) / (float64)(time.Millisecond)
		toq := check / (float64)(queryCount)
		fmt.Printf("LoadTime:%f ProcessTime:%f TOQ:%f Count:%d Right:%d Wrong:%d\n",
			load, check, toq, queryCount, *right, *wrong)
	}(&load_start, &load_end, &check_start, &check_end, &right, &wrong, queryCount)

	load_start = time.Now().UnixNano()
	judger, err := gopolygonjudger.NewJudger(os.Args[1])
	if err != nil {
		fmt.Printf("Load %s failed:%v", os.Args[1], err)
		return
	}
	load_end = time.Now().UnixNano()
	check_start = time.Now().UnixNano()
	for idx, point := range query {
		ret := judger.FindCityId(*point)
		if ret == result[idx] {
			right += 1
		} else {
			wrong += 1
		}
    fmt.Printf("p:%v ret:%d ex:%d\n", *point, ret, result[idx])
	}
	check_end = time.Now().UnixNano()
}
