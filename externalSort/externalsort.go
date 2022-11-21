package externalsort

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"sync"

	"externalsort/util"
)

const (
	DATA_FOLDER            = "./data/"
	TMP_FOLDER             = "./data/tmp/"
	IN_FILE_SUFFIX         = ".in"
	OUT_FILE_SUFFIX        = ".out"
	TMP_FILE_SUFFIX        = ".tmp"
	READ_FILE_BUFFER_SIZE  = 4096
	WRITE_FILE_BUFFER_SIZE = 4096
)

type ExternalSort struct {
	filename     string
	blockNum     int
	blockSize    []uint64
	tmpCount     int
	tmpCountLock sync.Mutex
}

func (s *ExternalSort) writeToBuffer(fout *bufio.Writer, x uint64) (err error) {
	tmp := make([]byte, 8)
	binary.LittleEndian.PutUint64(tmp, x)
	_, err = fout.Write(tmp)
	if err != nil {
		util.ErrorPrint(err, nil, "Write to buffer failed.")
		return err
	}
	return nil
}

func (s *ExternalSort) writeToFile(fout *os.File, x uint64) (err error) {
	tmp := make([]byte, 8)
	binary.LittleEndian.PutUint64(tmp, x)
	_, err = fout.Write(tmp)
	if err != nil {
		util.ErrorPrint(err, nil, "Write to file failed.")
		return err
	}
	return nil
}

func (s *ExternalSort) createReaderBuffer(filename string, size uint64) (fin *bufio.Reader, fileClose func() error, err error) {
	inputFile, err := os.Open(filename)
	if err != nil {
		util.ErrorPrint(err, nil, "Open input file failed.")
		return nil, nil, err
	}
	return bufio.NewReaderSize(inputFile, int(size)), inputFile.Close, nil
}

func (s *ExternalSort) createWriterBuffer(filename string, size uint64) (fout *bufio.Writer, fileClose func() error, err error) {
	s.tmpCountLock.Lock()
	s.tmpCount++
	tmpCount := s.tmpCount
	s.tmpCountLock.Unlock()
	outputFile, err := os.OpenFile(fmt.Sprintf("%s%s.%d%s",
		TMP_FOLDER, filename, tmpCount, TMP_FILE_SUFFIX),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		util.ErrorPrint(err, nil, fmt.Sprintf("Open output file %d failed.", tmpCount))
		return nil, nil, err
	}
	return bufio.NewWriterSize(outputFile, int(size)), outputFile.Close, nil
}

func (s *ExternalSort) readOneUint64FromBuffer(fin *bufio.Reader) (x uint64, err error) {
	tmp := make([]byte, 8)
	n, err := fin.Read(tmp)
	if n < 8 {
		_, err = fin.Read(tmp[n:])
	}
	return binary.LittleEndian.Uint64(tmp), err
}

func (s *ExternalSort) readOneUint64FromFile(fin *os.File) (x uint64, err error) {
	tmp := make([]byte, 8)
	_, err = fin.Read(tmp)
	return binary.LittleEndian.Uint64(tmp), err
}

func (s *ExternalSort) internalSort(fin *os.File, blockSize uint64) (err error) {
	s.tmpCount++
	outputFile, err := os.OpenFile(fmt.Sprintf("%s%s.%d%s",
		TMP_FOLDER, s.filename, s.tmpCount, TMP_FILE_SUFFIX), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		util.ErrorPrint(err, nil, "Open pre-writeen tmp file failed.")
		return err
	}
	a := []uint64{}
	for i := 0; i < int(blockSize); i += 8 {
		x, err := s.readOneUint64FromFile(fin)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				util.ErrorPrint(err, nil, "Reading input file failed.")
				return err
			}
		}
		a = append(a, x)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	for _, x := range a {
		err := s.writeToFile(outputFile, x)
		if err != nil {
			return err
		}
	}
	s.blockSize = append(s.blockSize, blockSize)
	fmt.Printf("Internal Sort: length = %d, from 0x%x to 0x%x\n", len(a), a[0], a[len(a)-1])
	outputFile.Close()
	return nil
}

func (s *ExternalSort) mergeSort(indexBegin int, k int, bufferSize uint64,
	finishFlag chan<- int) (err error) {
	outputBuffer, closeOutputFunc, err := s.createWriterBuffer(s.filename, bufferSize)
	if err != nil {
		return err
	}
	type inputType struct {
		buf       *bufio.Reader
		close     func() error
		leng      uint64
		completed bool
	}
	input := []inputType{}
	num := []uint64{}
	minIndex := -1
	for i := 0; i < k; i++ {
		buf, fun, err := s.createReaderBuffer(
			fmt.Sprintf("%s%s.%d%s", TMP_FOLDER, s.filename, indexBegin+i, TMP_FILE_SUFFIX), bufferSize)
		if err != nil {
			return err
		}

		x, err := s.readOneUint64FromBuffer(buf)
		if err != nil {
			if err == io.EOF {
				fun()
			} else {
				util.ErrorPrint(err, buf, "Read file failed for the first time in merge sort.")
				return err
			}
		}
		input = append(input, inputType{buf, fun, s.blockSize[indexBegin+i-1]/8 - 1, false})
		num = append(num, x)
		if minIndex == -1 || x < num[minIndex] {
			minIndex = i
		}
	}

	leng := 0
	min := num[minIndex]
	max := uint64(0)
	for minIndex > -1 {
		s.writeToBuffer(outputBuffer, num[minIndex])
		max = num[minIndex]
		leng++

		if input[minIndex].leng > 0 {
			num[minIndex], err = s.readOneUint64FromBuffer(input[minIndex].buf)
			input[minIndex].leng--
			if err != nil {
				if err == io.EOF {
					input[minIndex].close()
					input[minIndex].leng = 0
					input[minIndex].completed = true
					minIndex = -1
				} else {
					util.ErrorPrint(err, input[minIndex], "Read file failed in merge sort.")
					return err
				}
			} else {
				if input[minIndex].leng == 0 {
					input[minIndex].close()
				}
			}
		} else {
			input[minIndex].completed = true
			minIndex = -1
		}
		for i := 0; i < k; i++ {
			if !input[i].completed && (minIndex == -1 || num[i] < num[minIndex]) {
				minIndex = i
			}
		}
	}

	s.blockSize = append(s.blockSize, uint64(leng*8))
	fmt.Printf("Merge Sort: length = %d, index from %d to %d, num from 0x%x to 0x%x\n",
		leng, indexBegin, indexBegin+k-1, min, max)
	outputBuffer.Flush()
	closeOutputFunc()
	finishFlag <- 1
	return nil
}

func (s *ExternalSort) Sort(filename string, fileSize uint64,
	memLimit uint64, k int) (err error) {
	if memLimit%8 != 0 {
		return errors.New("memory limit should be divisible by 8")
	}

	s.filename = filename
	os.Mkdir(TMP_FOLDER, 0644)

	inputFile, err := os.Open(DATA_FOLDER + s.filename + IN_FILE_SUFFIX)
	if err != nil {
		util.ErrorPrint(err, nil, "Read input file failed.")
		return err
	}
	s.blockNum = int((fileSize-1)/memLimit) + 1
	s.tmpCount = 0
	for i := 0; i < s.blockNum; i++ {
		err = s.internalSort(inputFile, memLimit)
		if err != nil {
			util.ErrorPrint(err, nil, fmt.Sprintf("Internal sort %d failed", i))
			return err
		}
	}
	inputFile.Close()

	nowIndex := 1
	for s.blockNum > 1 {
		var ioBufferSize uint64
		thisIndex := s.tmpCount
		var finishFlag chan int
		t := 0
		if s.blockNum%k == 1 {
			ioBufferSize = uint64(math.Max(
				float64(memLimit/uint64(s.blockNum-1+s.blockNum/k)), 1))
			t = s.blockNum / k
			finishFlag = make(chan int, t)
		} else {
			ioBufferSize = uint64(math.Max(
				float64(memLimit/uint64(s.blockNum+(s.blockNum-1)/k+1)), 1))
			t = (s.blockNum-1)/k + 1
			finishFlag = make(chan int, t)
		}
		for nowIndex < thisIndex {
			if thisIndex-nowIndex+1 < k {
				go s.mergeSort(nowIndex, thisIndex-nowIndex+1, ioBufferSize, finishFlag)
				nowIndex = thisIndex + 1
			} else {
				go s.mergeSort(nowIndex, k, ioBufferSize, finishFlag)
				nowIndex += k
			}
		}

		flagCount := 0
		for x := range finishFlag {
			flagCount += x
			if flagCount == t {
				break
			}
		}
		close(finishFlag)
		s.blockNum = s.tmpCount - nowIndex + 1
	}

	err = os.Rename(fmt.Sprintf("%s%s.%d%s", TMP_FOLDER, s.filename,
		s.tmpCount, TMP_FILE_SUFFIX), DATA_FOLDER+s.filename+OUT_FILE_SUFFIX)
	if err != nil {
		util.ErrorPrint(err, nil, fmt.Sprintf("Failed to rename output file %s%s.%d%s", TMP_FOLDER,
			s.filename, s.tmpCount, TMP_FILE_SUFFIX))
		return nil
	}

	err = os.RemoveAll(TMP_FOLDER)
	if err != nil {
		util.ErrorPrint(err, nil, "Error occur while remove tmp file.")
	}

	return nil
}
