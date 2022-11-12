package externalsort

import (
	"bufio"
	"os"

	"externalsort/util"
)

const (
	DATA_FOLDER            = "./data/"
	IN_FILE_SUFFIX         = ".in"
	OUT_FILE_SUFFIX        = ".out"
	READ_FILE_BUFFER_SIZE  = 4096 // 3 times buffer in memory (2 read + 1 write)
	WRITE_FILE_BUFFER_SIZE = 4096
)

type ExternalSort struct {
	fin  *bufio.Reader
	fout *bufio.Writer
}

func (s *ExternalSort) Sort(filename string) (err error) {
	file, err := os.Open(DATA_FOLDER + filename + IN_FILE_SUFFIX)
	if err != nil {
		util.ErrorPrint(err, nil, "Open reading file failed.")
		return err
	}
	defer file.Close()
	s.fin = bufio.NewReaderSize(file, READ_FILE_BUFFER_SIZE)

	file, err = os.OpenFile(DATA_FOLDER+filename+OUT_FILE_SUFFIX,
		os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		util.ErrorPrint(err, nil, "Open writing file failed.")
		return err
	}
	defer file.Close()
	s.fout = bufio.NewWriterSize(file, WRITE_FILE_BUFFER_SIZE)
	defer s.fout.Flush()
	p := make([]byte, READ_FILE_BUFFER_SIZE)
	for _, err := s.fin.Read(p); err == nil; _, err = s.fin.Read(p) {
		_, err = s.fout.Write(p)
		if err != nil {
			return err
		}
	}

	return nil
}
