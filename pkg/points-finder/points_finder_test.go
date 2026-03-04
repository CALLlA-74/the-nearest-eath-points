package points_finder

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"sort"
	"testing"
)

func TestEarthPoint_Distance(t *testing.T) {
	p1 := NewPoint(1, 1, 1)
	p2 := NewPoint(1, -1, 2)
	assert.Equal(t, p1.Distance(p2), p2.Distance(p1))
	assert.Less(t, float64(0), p1.Distance(p2))
	assert.Less(t, float64(0), p2.Distance(p1))

	p3 := NewPoint(1, -179, 3)
	p4 := NewPoint(1, 179, 4)
	assert.Equal(t, p3.Distance(p4), p4.Distance(p3))
	assert.LessOrEqual(t, math.Abs(p1.Distance(p2)-p3.Distance(p4)), 1e-6)
	assert.Less(t, float64(0), p3.Distance(p4))
	assert.Less(t, float64(0), p4.Distance(p3))
}

func TestEarthPoint_Distance2(t *testing.T) {
	p1 := NewPoint(1, 0, 1)
	p2 := NewPoint(1, 0, 2)
	assert.Equal(t, p1.Distance(p2), p2.Distance(p1))
	assert.LessOrEqual(t, float64(0), p1.Distance(p2))
	assert.LessOrEqual(t, float64(0), p2.Distance(p1))

	p3 := NewPoint(1, -180, 3)
	p4 := NewPoint(1, 180, 4)
	assert.Equal(t, p3.Distance(p4), p4.Distance(p3))
	assert.LessOrEqual(t, math.Abs(p1.Distance(p2)-p3.Distance(p4)), 1e-6)
	assert.LessOrEqual(t, float64(0), p3.Distance(p4))
	assert.LessOrEqual(t, float64(0), p4.Distance(p3))
}

func TestEarthPoint_Distance3(t *testing.T) {
	p1 := NewPoint(90, 0, 1)
	p2 := NewPoint(90, -180, 2)
	assert.LessOrEqual(t, math.Abs(p1.Distance(p2)), 1e-6)
	assert.LessOrEqual(t, float64(0), p1.Distance(p2))
	assert.LessOrEqual(t, float64(0), p2.Distance(p1))

	p3 := NewPoint(-90, 0, 3)
	p4 := NewPoint(-90, 180, 4)
	assert.LessOrEqual(t, math.Abs(p3.Distance(p4)), 1e-6)
	assert.LessOrEqual(t, float64(0), p3.Distance(p4))
	assert.LessOrEqual(t, float64(0), p4.Distance(p3))
}

func TestEarthPoint_Distance4(t *testing.T) {
	p1 := NewPoint(-1, -1, 1)
	p2 := NewPoint(1, 1, 2)
	p3 := NewPoint(1, -1, 3)
	p4 := NewPoint(-1, 1, 4)

	assert.LessOrEqual(t, float64(0), p1.Distance(p2))
	assert.LessOrEqual(t, float64(0), p2.Distance(p1))
	assert.LessOrEqual(t, float64(0), p3.Distance(p4))
	assert.LessOrEqual(t, float64(0), p4.Distance(p3))

	assert.LessOrEqual(t, float64(0), p1.Distance(p3))
	assert.LessOrEqual(t, float64(0), p3.Distance(p1))
	assert.LessOrEqual(t, float64(0), p2.Distance(p4))
	assert.LessOrEqual(t, float64(0), p4.Distance(p2))

	assert.LessOrEqual(t, float64(0), p1.Distance(p4))
	assert.LessOrEqual(t, float64(0), p2.Distance(p3))
	assert.LessOrEqual(t, float64(0), p3.Distance(p2))
	assert.LessOrEqual(t, float64(0), p4.Distance(p1))

	assert.LessOrEqual(t, math.Abs(p1.Distance(p2)-p3.Distance(p4)), 1e-6)
	assert.LessOrEqual(t, math.Abs(p1.Distance(p3)-p2.Distance(p4)), 1e-6)
	assert.LessOrEqual(t, math.Abs(p1.Distance(p4)-p2.Distance(p3)), 1e-6)
}

func TestEarthPoint_PlaneDistance(t *testing.T) {
	p1 := NewPoint(90, 1, 1)
	p2 := NewPoint(0, 0, 2)
	p3 := NewPoint(-90, -1, 3)

	assert.Equal(t, p1.Distance(p2), p3.Distance(p2))
}

/*
1. 1 центр. +- 20град по осям. 1000 точек. запросить 20 ближайших к центру.
в тесте чекнуть: если точка не входит в множество ближайших, то расстояние до нее
должно быть >= расстояния до самой дальней (чекаем с eps 1e-6)
*/
func TestFindNearestPoints(t *testing.T) {
	centers := []*EarthPoint{
		NewPoint(10.5, 90, 1),
	}
	testFindingByCenters(t, centers, 10, 100)
}

var centers = []*EarthPoint{
	NewPoint(0, 0, 1),
	NewPoint(0, 180, 2),
	NewPoint(90, 180, 3),
	NewPoint(-90, 0, 4),
	NewPoint(-7.63657, 154.00466, 5),
	NewPoint(-6.08252, -147.06346, 6),
	NewPoint(-34.139821940712395, 129.1656234317047, 7),
	NewPoint(-52.43359558106289, 89.42477495458306, 8),
	NewPoint(-68.74566685663875, 29.957330060173316, 9),
	NewPoint(-27.282307793440893, 23.019759781288922, 10),
	NewPoint(-24.579986363626066, 2.3492128937321404, 11),
	NewPoint(-16.36912903601095, -25.025680139865845, 12),
	NewPoint(-29.71143856224785, -52.44616530589945, 13),
	NewPoint(27.06698948142664, -94.56576104675804, 14),
	NewPoint(6.543793827541402, -103.4786845531056, 15),
	NewPoint(32.576311809635676, -140.481673375893, 16),
	NewPoint(-20.79962802312857, -115.83685532895203, 17),
	NewPoint(56.915855002856645, 76.21232545409316, 18),
	NewPoint(32.29170387727744, 58.5367970390593, 19),
	NewPoint(6.421949732453414, 93.25919453735887, 20),
}

/*
2. 20 центров. +- 10град по осям. (центры находятся минимум на расстоянии 20 градусов по 1 из осей)
по 5к точек на центр. запросить для каждой ближайшие 20 точек. в id у точек выделить 2 младших
разряда под id центра. в тесте чекать, что найденные точки действительно принадлежат центру
*/
func TestFindNearestPoints2(t *testing.T) {
	testFindingByCenters(t, centers, 10, 100000)
}

func testFindingByCenters(t *testing.T, centers []*EarthPoint, delta float64, numOf int) {
	points := generatePoints(centers, delta, numOf)
	for _, center := range centers {
		nears := FindNearestPoints(points, center, 20)
		mp := make(map[int64]*EarthPoint, len(nears))

		for _, near := range nears {
			mp[near.Id] = near
			assert.Equal(t, center.Id, near.Id%100)
		}
		assert.Equal(t, len(nears), len(mp))
		sort.Slice(nears, func(i, j int) bool {
			return nears[i].Distance(center) > nears[j].Distance(center)
		})

		for _, point := range points {
			if _, ok := mp[point.Id]; !ok {
				assert.LessOrEqual(t, nears[0].Distance(center), (*point).Distance(center))
			}
		}
	}
}

func generatePoints(centers []*EarthPoint, delta float64, numOf int) []*EarthPoint {
	numOfPerCenter := numOf / len(centers)

	ret := make([]*EarthPoint, 0, numOfPerCenter*len(centers))
	id := 1
	for _, center := range centers {
		for i := 0; i < numOfPerCenter; i++ {
			ret = append(
				ret,
				NewPoint(
					normalize(
						(center.LatDegrees-math.Abs(delta))+rand.Float64()*math.Abs(delta)*2,
						90,
					),
					normalize(
						(center.LonDegrees-math.Abs(delta))+rand.Float64()*math.Abs(delta)*2,
						180,
					),
					int64(id*100)+center.Id,
				),
			)
			id++
		}
	}
	return ret
}

func normalize(v float64, norm int64) float64 {
	if math.Abs(v) > float64(norm) {
		if v < float64(norm) {
			v += (math.Abs(math.Abs(v) - float64(norm))) * 2
		} else {
			v -= (math.Abs(math.Abs(v) - float64(norm))) * 2
		}

		if norm == 180 {
			v *= -1
		}
	}

	return v
}

func BenchmarkFindNearestPoints3(b *testing.B) {
	points := generatePoints(centers, 10, 1e6)
	for i := 0; i < b.N; i++ {
		FindNearestPoints(points, centers[rand.Int()%len(centers)], 100)
	}
}
