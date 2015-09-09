package gopolygonjudger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"path"
)

const (
	InBoxName          = "inbox.txt"
	OutBoxName         = "outbox.txt"
	CityEdgeNamePrefix = "city_"
	CityEdgeNameSuffix = ".txt"
)

type Point struct {
	Lat float64
	Lng float64
}

type Polygon []*Point

type Rectangle struct {
	LB Point
	RT Point
}

type AreaJudger struct {
	InBox    []*Rectangle
	OutBox   []*Rectangle
	Edges    []*Polygon
	DataPath string
}

func NewJudger(dataDirPath string) (areaJudger *AreaJudger, err error) {
	inBox, err := parseRectangle(path.Join(dataDirPath, InBoxName))
	if err != nil {
		return
	}

	outBox, err := parseRectangle(path.Join(dataDirPath, OutBoxName))
	if err != nil {
		return
	}

	if len(inBox) != len(outBox) {
		return
	}

	edges, err := parseCityEdges(dataDirPath, len(inBox))
	if err != nil {
		return
	}

	areaJudger = &AreaJudger{inBox, outBox, edges, dataDirPath}
	return
}

func (judger *AreaJudger) ToStdout() {
	for _, item := range judger.InBox {
		fmt.Printf("\t%v\n", *item)
	}
	for _, item := range judger.OutBox {
		fmt.Printf("\t%v\n", *item)
	}
	for idx, item := range judger.Edges {
		fmt.Printf("\tCity %d\n", idx + 1)
		for _, point := range *item {
			fmt.Printf("\t\t%v\n", *point)
		}
	}
}

func (judger *AreaJudger) FindCityId(point Point) int {
	id := judger.inBoxMatch(point)
	if id >= 0 {
		return id
	}

	ids := judger.outBoxMatch(point)
	if len(ids) <= 0 {
		return -1
	}

	for _, id := range ids {
		if judger.polygonMatch(id, point) {
			return id
		}
	}

	return -1
}

func (judger *AreaJudger) inBoxMatch(point Point) int {
	for idx, rec := range judger.InBox {
		if InRectangle(rec, point) {
			return idx
		}
	}
	return -1
}

func (judger *AreaJudger) outBoxMatch(point Point) []int {
	ids := make([]int, 0)
	for idx, rec := range judger.OutBox {
		if InRectangle(rec, point) {
			ids = append(ids, idx)
		}
	}
	return ids
}

func (judger *AreaJudger) polygonMatch(id int, point Point) bool {
	if id >= len(judger.Edges) {
		return false
	}
	edges := *judger.Edges[id]
	if len(edges) <= 2 {
		return false
	}

	from := edges[0]
	crossCount := 0
	idx := 1
	for ; idx < len(edges); idx += 1 {
		to := edges[idx]
		if ((point.Lat <= to.Lat && point.Lat > from.Lat) || (point.Lat > to.Lat && point.Lat <= from.Lat)) && (point.Lng >= from.Lng || point.Lng >= to.Lng) {
			if from.Lat == to.Lat {
				if from.Lat < point.Lat {
					crossCount += 1
				} else if from.Lat == point.Lat {
					return false
				}
			} else {
				crossLng := (to.Lng-from.Lng)*(point.Lat-from.Lat)/(to.Lat-from.Lat) + from.Lng
				if crossLng < point.Lng {
					crossCount += 1
				}
			}
		}
		from = to
	}
	return (crossCount & 1) != 0
}

func InRectangle(rec *Rectangle, point Point) bool {
	if point.Lat >= rec.LB.Lat && point.Lng >= rec.LB.Lng && point.Lat <= rec.RT.Lat && point.Lng <= rec.RT.Lng {
		return true
	}
	return false
}

func parseCityEdges(dataDirPath string, count int) (edges []*Polygon, err error) {
	tmpEdges := make([]*Polygon, 0)
	var polygon *Polygon
	for idx := 0; idx < count; idx += 1 {
		filename := CityEdgeNamePrefix + strconv.Itoa(idx+1) + CityEdgeNameSuffix
		polygon, err = parsePolygon(path.Join(dataDirPath, filename))
		if err != nil {
			return
		}
		tmpEdges = append(tmpEdges, polygon)
	}
	edges = tmpEdges
	return
}

func parseRectangle(filename string) (recList []*Rectangle, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var lbLat, lbLng float64
	var rtLat, rtLng float64
	var count int

	tmpRecList := make([]*Rectangle, 0)
	count, err = fmt.Fscanf(reader, "%f,%f,%f,%f\n", &lbLng, &lbLat, &rtLng, &rtLat)
	for err == nil && count == 4 {
		tmpRecList = append(tmpRecList, &Rectangle{Point{lbLat, lbLng}, Point{rtLat, rtLng}})
		count, err = fmt.Fscanf(reader, "%f,%f,%f,%f\n", &lbLng, &lbLat, &rtLng, &rtLat)
	}
	if err == io.EOF {
		recList = tmpRecList
		err = nil
	}
	return
}

func parsePolygon(filename string) (polygon *Polygon, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var lat, lng float64
	var count int

	tmpPoints := make([]*Point, 0)
	count, err = fmt.Fscanf(reader, "%f,%f\n", &lng, &lat)
	for err == nil && count == 2 {
		tmpPoints = append(tmpPoints, &Point{lat, lng})
		count, err = fmt.Fscanf(reader, "%f,%f\n", &lng, &lat)
	}
	if err == io.EOF {
		tPolygon := Polygon(tmpPoints)
		polygon = &tPolygon
		err = nil
	}
	return
}
