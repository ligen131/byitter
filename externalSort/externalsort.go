package externalsort

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"

	"externalsort/util"
)

const (
	DATA_FOLDER                  = "./data/"
	IN_FILE_SUFFIX               = ".in"
	OUT_FILE_SUFFIX              = ".out"
	INTERAL_SORT_TMP_FILE_SUFFIX = ".in.tmp"
	MERGE_SORT_TMP_FILE_SUFFIX   = ".me.tmp"
	READ_FILE_BUFFER_SIZE        = 4096
	WRITE_FILE_BUFFER_SIZE       = 4096
)

type ExternalSort struct {
	blockNum int
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

func (s *ExternalSort) chanWriteToBuffer(fout *bufio.Writer, c <-chan uint64) (err error) {
	for x := range c {
		err = s.writeToBuffer(fout, x)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ExternalSort) readByteToUint64(fin *bufio.Reader,
	size uint64, c chan<- uint64) {
	// defer close(c)
	for i := uint64(0); i < size; i += 8 {
		tmp := make([]byte, 8)
		_, err := fin.Read(tmp)
		if err != nil {
			close(c)
			return
		}
		c <- binary.LittleEndian.Uint64(tmp)
	}
	close(c)
}

func (s *ExternalSort) readToChan(fin *bufio.Reader,
	size uint64) (c <-chan uint64) {
	ch := make(chan uint64, 1024)
	go s.readByteToUint64(fin, size, ch)
	return ch
}

func (s *ExternalSort) internalSort(inChan <-chan uint64, fout *bufio.Writer) (err error) {
	a := []uint64{}
	for x := range inChan {
		a = append(a, x)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	for _, x := range a {
		err := s.writeToBuffer(fout, x)
		if err != nil {
			return err
		}
	}
	fmt.Printf("%d 0x%x 0x%x\n", len(a), a[0], a[len(a)-1])
	return nil
}

func (s *ExternalSort) Sort(filename string, fileSize uint64,
	memLimit uint64, k int) (err error) {
	if memLimit%8 != 0 {
		return errors.New("memory limit should be divisible by 8")
	}

	inputFile, err := os.Open(DATA_FOLDER + filename + IN_FILE_SUFFIX)
	if err != nil {
		util.ErrorPrint(err, nil, "Open input file failed.")
		return err
	}

	inputBuffer := bufio.NewReaderSize(inputFile, int(memLimit))
	s.blockNum = int(math.Ceil(float64(fileSize) / float64(memLimit)))
	for i := 0; i < s.blockNum; i++ {
		outputFile, err := os.OpenFile(fmt.Sprintf("%s%s.%d%s",
			DATA_FOLDER, filename, i, INTERAL_SORT_TMP_FILE_SUFFIX),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			util.ErrorPrint(err, nil, fmt.Sprintf("Open output file %d failed.", i))
			return err
		}
		outputBuffer := bufio.NewWriter(outputFile)

		inChan := s.readToChan(inputBuffer, memLimit)
		err = s.internalSort(inChan, outputBuffer)
		if err != nil {
			util.ErrorPrint(err, nil, fmt.Sprintf("Internal sort %d failed", i))
			return err
		}
		outputBuffer.Flush()
		outputFile.Close()
	}
	inputFile.Close()

	// fileIn, err := os.Open(DATA_FOLDER + filename + IN_FILE_SUFFIX)
	// if err != nil {
	// 	util.ErrorPrint(err, nil, "Open reading file failed.")
	// 	return err
	// }
	// defer fileIn.Close()
	// s.fin = bufio.NewReaderSize(fileIn, READ_FILE_BUFFER_SIZE)

	// fileOut, err := os.OpenFile(DATA_FOLDER+filename+OUT_FILE_SUFFIX,
	// 	os.O_WRONLY|os.O_CREATE, 0644)
	// if err != nil {
	// 	util.ErrorPrint(err, nil, "Open writing file failed.")
	// 	return err
	// }
	// defer fileOut.Close()
	// s.fout = bufio.NewWriterSize(fileOut, WRITE_FILE_BUFFER_SIZE)
	// defer s.fout.Flush()

	// p := make([]byte, READ_FILE_BUFFER_SIZE)
	// for _, err := s.fin.Read(p); err == nil; _, err = s.fin.Read(p) {
	// 	_, err = s.fout.Write(p)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
