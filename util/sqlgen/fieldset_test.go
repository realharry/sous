package sqlgen

import "testing"

func TestRowDedupe(t *testing.T) {
	fs := NewFieldset()

	fs.Row(func(rd RowDef) { rd.CF("?", "testcol", 1) })
	fs.Row(func(rd RowDef) { rd.CF("?", "testcol", 1) })
	fs.Row(func(rd RowDef) { rd.CF("?", "testcol", 2) })

	sql := fs.InsertSQL("testtable", "")
	expectedSQL := "insert into testtable (testcol) values ($1),\n($2) "
	if sql != expectedSQL {
		t.Errorf("Expected fieldset to generate %q; got %q", expectedSQL, sql)
	}

	vals := fs.InsertValues()
	if len(vals) != 2 {
		t.Errorf("Expected 2 values, got %d", len(vals))
	}

	if vals[0] != 1 {
		t.Errorf("Expected first value to be 1, got %v", vals[0])
	}
	if vals[1] != 2 {
		t.Errorf("Expected first value to be 1, got %v", vals[1])
	}
}

func TestRowDedupe_without_candidates(t *testing.T) {
	fs := NewFieldset()

	fs.Row(func(rd RowDef) { rd.FD("?", "testcol", 1) })
	fs.Row(func(rd RowDef) { rd.FD("?", "testcol", 1) })
	fs.Row(func(rd RowDef) { rd.FD("?", "testcol", 2) })

	sql := fs.InsertSQL("testtable", "")
	expectedSQL := "insert into testtable (testcol) values ($1),\n($2),\n($3) "
	if sql != expectedSQL {
		t.Errorf("Expected fieldset to generate %q; got %q", expectedSQL, sql)
	}

	vals := fs.InsertValues()
	if len(vals) != 3 {
		t.Errorf("Expected 3 values, got %d", len(vals))
	}
}

func TestRowDedupe_mixed_fields(t *testing.T) {
	fs := NewFieldset()

	fs.Row(func(rd RowDef) {
		rd.CF("?", "candidate", 1)
		rd.FD("?", "other", "something")
	})
	fs.Row(func(rd RowDef) {
		rd.CF("?", "candidate", 1)
		// note that non-candidate fields aren't considered for de-duplications
		// that's what a candidate field means in this context, though.
		rd.FD("?", "other", "or other")
	})
	fs.Row(func(rd RowDef) {
		rd.CF("?", "candidate", 2)
		rd.FD("?", "other", "also this")
	})

	sql := fs.InsertSQL("testtable", "")
	expectedSQL := "insert into testtable (candidate,other) values ($1,$2),\n($3,$4) "
	if sql != expectedSQL {
		t.Errorf("Expected fieldset to generate %q; got %q", expectedSQL, sql)
	}

	vals := fs.InsertValues()
	if len(vals) != 4 {
		t.Errorf("Expected 4 values, got %d", len(vals))
	}
	firstC := vals[0]
	expectedC := 1
	if firstC != expectedC {
		t.Errorf("Expected first candidate to be %v; got %v", expectedC, firstC)
	}

	firstNC := vals[1]
	expectedNC := "something"
	if firstNC != expectedNC {
		t.Errorf("Expected first non-candidate to be %v; got %v", expectedNC, firstNC)
	}

	secondNC := vals[3]
	expectedNC = "also this"
	if secondNC != expectedNC {
		t.Errorf("Expected other non-candidate to be %v; got %v", expectedNC, secondNC)
	}
}
