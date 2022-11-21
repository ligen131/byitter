package externalsort_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"

	. "externalsort"
)

func sortA(a *[]uint64, flag chan<- bool) {
	sort.Slice(*a, func(i, j int) bool {
		return (*a)[i] < (*a)[j]
	})
	flag <- true
}

func TestExternalSort(t *testing.T) {
	const FILENAME = "test"
	const FILESIZE = 8 * (1 << 22)
	const MEMORY_LIMIT = 1 << 15
	file, err := os.OpenFile(DATA_FOLDER+FILENAME+IN_FILE_SUFFIX,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Error(err, "Open test writing file failed.")
		return
	}

	const n = FILESIZE / 8
	tmp := make([]byte, 8)
	a := []uint64{}
	for i := 0; i < n; i++ {
		x := rand.Uint64()
		a = append(a, x)
		binary.LittleEndian.PutUint64(tmp, x)
		file.Write(tmp)
	}
	flag := make(chan bool)
	go sortA(&a, flag)
	file.Close()

	s := ExternalSort{}
	err = s.Sort(FILENAME, FILESIZE, MEMORY_LIMIT, 2)
	if err != nil {
		t.Error(err, "Sort failed.")
		return
	}
	<-flag
	file, err = os.OpenFile(DATA_FOLDER+FILENAME+OUT_FILE_SUFFIX,
		os.O_RDONLY, 0644)
	if err != nil {
		t.Error(err, "Open answer file failed.")
		return
	}
	for i := 0; i < n; i++ {
		_, err := file.Read(tmp)
		if err != nil {
			t.Error(err, "Read answer file failed.")
			return
		}
		x := binary.LittleEndian.Uint64(tmp)
		if x != a[i] {
			t.Log(fmt.Sprintf("Wrong answer at %d number, read 0x%x, expected 0x%x.", i, x, a[i]))
			t.FailNow()
		}
	}
}
