package mapSimplify

import (
	"fmt"
)

type SegmentId [2]int

type segmentIdList []SegmentId
func (l segmentIdList) Len() int { return len(l) }
func (l segmentIdList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l segmentIdList) Less(i, j int) bool { return l[i][0] < l[j][0] }

type Segment struct {
	Points []Point
	ref1 *Polygon
	ref2 *Polygon
}

func (s *Segment) LeftId() int {
	if s.ref1 == nil {
		return 0
	}
	return s.ref1.Id
}

func (s *Segment) RightId() int {
	if s.ref2 == nil {
		return 0
	}
	return s.ref2.Id
}

func (s *Segment) String() string {
	pts := make([]string, len(s.Points))
	for i, pt := range s.Points {
		pts[i] = pt.String()
	}
	return fmt.Sprintf("[%s]:(%d, %d)", pts, s.LeftId(), s.RightId())
}
