package points_finder

import (
	"container/heap"
	"math"
)

type EarthPoint struct {
	LatDegrees, LonDegrees float64
	latRadian              float64 // latitude at radian
	lonRadian              float64 // longitude at radian
	Id                     int64
}

func (p *EarthPoint) Dim() int {
	return 2
}

func (p *EarthPoint) GetValue(dim int) float64 {
	switch dim {
	case 0:
		return p.latRadian
	case 1:
		return p.lonRadian
	default:
		panic("out of dimension")
	}
}

func (p *EarthPoint) Distance(other *EarthPoint) float64 {
	clat1 := math.Cos(p.latRadian)
	clat2 := math.Cos(other.GetValue(0))
	slat1 := math.Sin(p.latRadian)
	slat2 := math.Sin(other.GetValue(0))
	delta := other.GetValue(1) - p.lonRadian
	cdelta := math.Cos(delta)
	sdelta := math.Sin(delta)

	y := math.Sqrt(math.Pow(clat2*sdelta, 2) + math.Pow(clat1*slat2-slat1*clat2*cdelta, 2))
	x := slat1*slat2 + clat1*clat2*cdelta

	return math.Atan2(y, x)
}

func (p *EarthPoint) PlaneDistance(val float64, dim int) float64 {
	var other *EarthPoint
	switch dim {
	case 0:
		other = NewPoint(val, p.GetValue(1), -1)
		break
	case 1:
		other = NewPoint(p.GetValue(0), val, -1)
		break
	default:
		panic("out of dimension")
	}
	return p.Distance(other)
}

func NewPoint(latDegrees, lonDegrees float64, id int64) *EarthPoint {
	return &EarthPoint{
		LatDegrees: latDegrees,
		LonDegrees: lonDegrees,
		latRadian:  latDegrees * math.Pi / 180,
		lonRadian:  lonDegrees * math.Pi / 180,
		Id:         id,
	}
}

func FindNearestPoints(points []*EarthPoint, target *EarthPoint, numOfNearest int) []*EarthPoint {
	hp := newHeap(numOfNearest)
	for _, point := range points {
		if hp.Len() < numOfNearest || (*hp)[0].distance > target.Distance(point) {
			heap.Push(hp, &heapNode{p: point, distance: target.Distance(point)})

			if hp.Len() > numOfNearest {
				heap.Pop(hp)
			}
		}
	}

	res := make([]*EarthPoint, numOfNearest)
	for i := numOfNearest - 1; i >= 0; i-- {
		res[i] = heap.Pop(hp).(*heapNode).p
	}
	return res
}

type heapHelper []*heapNode

func newHeap(numOf int) *heapHelper {
	ret := make(heapHelper, 0, numOf)
	return &ret
}

func (h heapHelper) Len() int {
	return len(h)
}

func (h heapHelper) Less(i, j int) bool {
	return h[i].distance > h[j].distance
}

func (h heapHelper) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapHelper) Push(x interface{}) {
	item := x.(*heapNode)
	*h = append(*h, item)
}

func (h *heapHelper) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type heapNode struct {
	p        *EarthPoint
	distance float64
}
