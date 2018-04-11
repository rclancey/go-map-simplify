package mapSimplify

import (
	"errors"
	"sort"
)

type Map struct {
	Polygons map[int]*Polygon
	pointIndex map[Point]map[int]bool
	simplified bool
}

func NewMap() *Map {
	return &Map{
		Polygons: make(map[int]*Polygon),
		pointIndex: make(map[Point]map[int]bool),
		simplified: false,
	}
}

func (m *Map) AddPolygon(id int, points []Point) (*Polygon, error) {
	poly, err := NewPolygon(id, points)
	if err != nil {
		return nil, err
	}
	m.Polygons[id] = poly
	m.simplified = false
	n := len(poly.Points) - 1
	var pt Point
	var ok bool
	var x map[int]bool
	for i := 0; i < n; i++ {
		pt = points[i]
		x, ok = m.pointIndex[pt]
		if ok {
			x[id] = true
		} else {
			x = make(map[int]bool)
			x[id] = true
			m.pointIndex[pt] = x
		}
	}
	return poly, nil
}

func (m *Map) polyIds() []int {
	polyIds := make([]int, len(m.Polygons))
	i := 0
	for id := range m.Polygons {
		polyIds[i] = id
		i++
	}
	sort.Ints(polyIds)
	return polyIds
}

func (m *Map) Simplify(threshold float64) ([]*Polygon, error) {
	polyIds := m.polyIds()
	out := make([]*Polygon, len(polyIds))
	for _, poly := range m.Polygons {
		poly.index()
	}
	for i, id := range polyIds {
		orig := m.Polygons[id]
		simple, err := orig.Simplify(m, threshold)
		if err != nil {
			return nil, err
		}
		out[i] = simple
		i++
	}
	m.simplified = true
	return out, nil
}

func (m *Map) Segments() ([]*Segment, error) {
	if !m.simplified {
		return nil, errors.New("Map not yet simplified")
	}
	polyIds := m.polyIds()
	out := []*Segment{}
	for _, id := range polyIds {
		orig := m.Polygons[id]
		for _, seg := range orig.segments {
			if seg.ref1 == orig {
				out = append(out, seg)
			}
		}
	}
	return out, nil
}

