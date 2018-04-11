package mapSimplify

import (
	"testing"
)

func TestBasicSimplify(t *testing.T) {
	m := NewMap()
	_, err := m.AddPolygon(1, []Point{
		Point{0.0, 0.0},
		Point{1.0, 0.0},
		Point{0.999, 0.5},
		Point{1.0, 1.0},
		Point{0.0, 1.0},
		Point{0.0, 0.0},
	})
	if err != nil {
		t.Error("Error adding polygon:", err)
	}
	_, err = m.AddPolygon(2, []Point{
		Point{1.0, 0.0},
		Point{2.0, 0.0},
		Point{2.0, 1.0},
		Point{1.0, 1.0},
		Point{0.999, 0.5},
		Point{1.0, 0.0},
	})
	polys, err := m.Simplify(0.1)
	if err != nil {
		t.Error("Error simplifying map:", err)
		return
	}
	if len(polys) != 2 {
		t.Errorf("Expected 2 polygons, got %d", len(polys))
		return
	}
	if len(polys[0].Points) != 5 {
		t.Errorf("Expected 5 points in first polygon, got %d", len(polys[0].Points))
	}
	if len(polys[1].Points) != 5 {
		t.Errorf("Expected 5 points in first polygon, got %d", len(polys[1].Points))
	}
}

func TestBasicSegments(t *testing.T) {
	m := NewMap()
	_, err := m.AddPolygon(1, []Point{
		Point{0.0, 0.0},
		Point{1.0, 0.0},
		Point{0.999, 0.5},
		Point{1.0, 1.0},
		Point{0.0, 1.0},
		Point{0.0, 0.0},
	})
	if err != nil {
		t.Error("Error adding polygon:", err)
	}
	_, err = m.AddPolygon(2, []Point{
		Point{1.0, 0.0},
		Point{2.0, 0.0},
		Point{2.0, 1.0},
		Point{1.0, 1.0},
		Point{0.999, 0.5},
		Point{1.0, 0.0},
	})
	_, err = m.Simplify(0.1)
	if err != nil {
		t.Error("Error simplifying map:", err)
		return
	}
	segs, err := m.Segments()
	if err != nil {
		t.Error("Error getting simplified segments:", err)
		return
	}
	if len(segs) != 4 {
		t.Errorf("Expected 4 segments, got %d", len(segs))
		return
	}
}

