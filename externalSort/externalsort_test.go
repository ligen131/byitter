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
	const FILESIZE = 8 * (1 << 20) // 1 << (3 + 20)
	const MEMORY_LIMIT = 1 << 15
	file, err := os.OpenFile(DATA_FOLDER+FILENAME+IN_FILE_SUFFIX,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Error(err, "Open test writing file failed.")
		return
	}
	defer file.Close()

	const n = FILESIZE / 8
	tmp := make([]byte, 8)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint64(tmp, rand.Uint64())
		file.Write(tmp)
	}

	s := ExternalSort{}
	err = s.Sort(FILENAME, FILESIZE, MEMORY_LIMIT, 2)
	if err != nil {
		t.Error(err, "Sort failed.")
		return
	}
}
