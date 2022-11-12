package externalsort_test

import (
	"encoding/binary"
	"math/rand"
	"os"
	"testing"

	. "externalsort"
)

func TestExternalSort(t *testing.T) {
	const FILENAME = "test"
	file, err := os.OpenFile(DATA_FOLDER+FILENAME+IN_FILE_SUFFIX,
		os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Error(err, "Open test writing file failed.")
		return
	}
	defer file.Close()

	const n = 102400
	tmp := make([]byte, 8)
	for i := 0; i < n; i++ {
		binary.BigEndian.PutUint64(tmp, rand.Uint64())
		file.Write(tmp)
	}

	s := ExternalSort{}
	err = s.Sort(FILENAME)
	if err != nil {
		t.Error(err, "Sort failed.")
		return
	}
}
