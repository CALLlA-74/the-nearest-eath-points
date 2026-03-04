package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	pointsFinder "github.com/callla-74/the-nearst-earth-points/pkg/points-finder"
	"math"
	"os"
	"strings"
)

func main() {
	menu()
}

const (
	enterPath    = "Введите пуль к файлу с описанием точек."
	ifWantToExit = "Если хотите завершить программу, введите: exit"
	enterInput   = "Введите координаты целевой точки (широта, долгота) и количество точек через пробел."
	example      = "Пример: 1.23 45.6 15"
	exitCmd      = "exit"
	bye          = "Bye!"
)

func menu() {
	var (
		path, inpS, errS     string
		targetLat, targetLon float64
		numOfNearestPoints   int
		err                  error
	)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(enterPath)
		fmt.Println(ifWantToExit)
		if scanner.Scan() {
			inpS = scanner.Text()
			if inpS == exitCmd {
				fmt.Println(bye)
				return
			}
			path = inpS
			break
		}
	}

	inputPoints := ScanPoints(path)

	for {
		err = nil
		fmt.Println(enterInput)
		fmt.Println(ifWantToExit)
		fmt.Println(example)
		if !scanner.Scan() {
			fmt.Println("")
		} else if inpS = scanner.Text(); inpS == exitCmd {
			fmt.Println(bye)
			return
		} else if _, err = fmt.Sscanf(inpS, "%f %f %d", &targetLat, &targetLon, &numOfNearestPoints); err != nil {
			fmt.Println(err)
		} else if errS = validate(targetLat, targetLon, numOfNearestPoints); len(errS) > 0 {
			fmt.Println(errS)
		} else {
			nearestPoints := pointsFinder.FindNearestPoints(
				inputPoints,
				pointsFinder.NewPoint(targetLat, targetLon, -1),
				numOfNearestPoints,
			)
			PrintResults(nearestPoints)
		}
	}
}

const (
	invalidLatFormat = "широта не должна превышать 90 по модулю: |%f| > 90\n"
	invalidLonFormat = "долгота не должна превышать 180 по модулю: |%f| > 180\n"
	invalidNumFormat = "количество точек не должно быть отрицательным:  %d < 0\n"
)

func validate(lat, lon float64, num int) string {
	var builder strings.Builder
	if math.Abs(lat) > 90 {
		builder.WriteString(fmt.Sprintf(invalidLatFormat, lat))
	}

	if math.Abs(lon) > 180 {
		builder.WriteString(fmt.Sprintf(invalidLonFormat, lon))
	}

	if num < 0 {
		builder.WriteString(fmt.Sprintf(invalidNumFormat, num))
	}

	return builder.String()
}

func PrintResults(result []*pointsFinder.EarthPoint) {
	fmt.Println("Ближайшие точки к заданной:")
	for i, p := range result {
		fmt.Printf("%d) {%f; %f}\n", i+1, p.LatDegrees, p.LonDegrees)
	}
}

type InputPoint struct {
	Lat float64 `json:"lat"` // latitude at degrees
	Lon float64 `json:"lon"` // longitude at degrees
}

func ScanPoints(path string) []*pointsFinder.EarthPoint {
	data, e := os.ReadFile(path)
	if e != nil {
		panic(e)
	}

	var points []InputPoint
	if e = json.Unmarshal(data, &points); e != nil {
		panic(e)
	}

	resPoints := make([]*pointsFinder.EarthPoint, 0, len(points))
	for i, p := range points {
		resPoints = append(resPoints, pointsFinder.NewPoint(p.Lat, p.Lon, int64(i)))
	}

	return resPoints
}
