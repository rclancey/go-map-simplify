package mapSimplify

import (
	"errors"
	"fmt"
	"sort"
	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/reducers"
)

type Polygon struct {
	Id int
	Points []Point
	pointIndex map[Point]int
	segments map[SegmentId]*Segment
	simplified bool
}

func NewPolygon(id int, points []Point) (*Polygon, error) {
	n := len(points) - 1
	if n <= 3 {
		return nil, fmt.Errorf("Polygon %d too small", id)
	}
	if points[0] != points[n] {
		return nil, fmt.Errorf("Polygon %d not closed", id)
	}
	poly := &Polygon{
		Id: id,
		Points: points,
		simplified: false,
	}
	return poly, nil
}

func (p *Polygon) index() error {
	p.pointIndex = make(map[Point]int)
	p.segments = make(map[SegmentId]*Segment)
	p.simplified = false
	var pt Point
	var ok bool
	n := len(p.Points) - 1
	for i := 0; i < n; i++ {
		pt = p.Points[i]
		if _, ok = p.pointIndex[pt]; ok {
			return fmt.Errorf("Choke point %d %s in polygon %d", i, pt, p.Id)
		}
		p.pointIndex[pt] = i
	}
	return nil
}

func (p *Polygon) makeSegments(m *Map, thresh float64) error {
	var start, end int
	var neighbor *Polygon
	var neighborStart, neighborEnd int
	var pt, next Point
	var shared, nextShared map[int]bool
	var simple []Point
	var ok bool
	n := len(p.Points) - 1
	i := 0
	var j int
	var segId SegmentId
	var seg *Segment
	for i < n {
		neighbor = nil
		start = i
		pt = p.Points[i]
		i++
		next = p.Points[i]
		shared = m.pointIndex[pt]
		nextShared = m.pointIndex[next]
		for id := range shared {
			if id != p.Id {
				if _, ok = nextShared[id]; ok {
					neighbor = m.Polygons[id]
					break
				}
			}
		}
		if neighbor == nil {
			// edge is not shared.  make a segment from here
			// to next shared vertex.
			for i <= n {
				end = i
				next = p.Points[i]
				nextShared = m.pointIndex[next]
				if len(nextShared) != 1 {
					break
				}
				i++
			}
			segId = SegmentId{start, end}
			simple = p.simplifySegment(start, end, thresh)
			p.segments[segId] = &Segment{
				Points: simple,
				ref1: p,
				ref2: nil,
			}
		} else {
			// edge is shared.  segment should be following edges
			// shared by the same set of polygons.  If we encounter
			// a vertex shared by more than two polygons, stop there
			j = neighbor.pointIndex[next]
			end = i
			neighborStart = j
			for i <= n && j >= 0 {
				next = p.Points[i]
				if next != neighbor.Points[j] {
					i--
					break
				}
				end = i
				neighborStart = j
				if len(m.pointIndex[next]) != 2 {
					break
				}
				i++
				j--
			}
			neighborEnd = neighborStart + (end - start)
			segId = SegmentId{start, end}
			if _, ok = p.segments[segId]; !ok {
				simple = p.simplifySegment(start, end, thresh)
				seg = &Segment{
					Points: simple,
					ref1: p,
					ref2: neighbor,
				}
				p.segments[segId] = seg
				segId = SegmentId{neighborStart, neighborEnd}
				neighbor.segments[segId] = seg
			}
		}
	}
	return nil
}

func (p *Polygon) simplifySegment(start, end int, thresh float64) []Point {
	if end - start <= 1 {
		return p.Points[start:end+1]
	}
	xydata := make([][2]float64, end+1-start)
	for i := start; i < end+1; i++ {
		xydata[i-start] = [2]float64{p.Points[i][0], p.Points[i][1]}
	}
	path := geo.NewPathFromXYData(xydata)
	simplePath := reducers.Visvalingam(path, thresh, 2)
	simple := make([]Point, simplePath.Length())
	for i, pt := range simplePath.Points() {
		simple[i] = Point{pt[0], pt[1]}
	}
	return simple
}

func (p *Polygon) fullSimplifiedPath() ([]Point, error) {
	segIds := make([]SegmentId, len(p.segments))
	i := 0
	for segId := range p.segments {
		segIds[i] = segId
		i++
	}
	sort.Sort(segmentIdList(segIds))
	if segIds[0][0] != 0 {
		return nil, errors.New("first point missing!")
	}
	prev := 0
	path := []Point{}
	var seg *Segment
	var rev []Point
	for i, segId := range segIds {
		if segId[0] < prev {
			return nil, fmt.Errorf("overlapping segments (%d, %d)", prev, segId[0])
		}
		if segId[0] > prev {
			return nil, fmt.Errorf("disjointed segments (%d, %d)", prev, segId[0])
		}
		prev = segId[1]
		seg = p.segments[segId]
		if seg.ref1 == p {
			if i == 0 {
				path = append(path, seg.Points...)
			} else {
				path = append(path, seg.Points[1:]...)
			}
		} else if seg.ref2 == p {
			m := len(seg.Points)
			rev = make([]Point, m)
			m--
			for j, pt := range seg.Points {
				rev[m - j] = pt
			}
			if i == 0 {
				path = append(path, rev...)
			} else {
				path = append(path, rev[1:]...)
			}
		} else {
			return nil, errors.New("unreferenced segment")
		}
	}
	if prev != len(p.Points) - 1 {
		return nil, errors.New("last point missing!")
	}
	return path, nil
}

func (p *Polygon) Simplify(m *Map, threshold float64) (*Polygon, error) {
	err := p.makeSegments(m, threshold)
	if err != nil {
		return nil, err
	}
	points, err := p.fullSimplifiedPath()
	if err != nil {
		return nil, err
	}
	return &Polygon{
		Id: p.Id,
		Points: points,
	}, nil
}

