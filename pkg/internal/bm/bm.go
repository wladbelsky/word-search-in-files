package bm

import (
	"bufio"
)

// Оставлю свою наработку по алгоритму Бойера-Мура на случай если понадобится

func Exists(s *bufio.Reader, substr string, size int) bool {
	return Index(s, substr, size) >= 0
}

func Index(s *bufio.Reader, substr string, size int) int {
	d := calculateSlideTable(substr)
	return indexWithTable(&d, s, substr, size)
}

func indexWithTable(d *[256]int, reader *bufio.Reader, substr string, ls int) int {
	lsub := len(substr)
	s := make([]byte, ls)
	readBuffer := make([]byte, lsub)
	switch {
	case lsub == 0:
		return -1
	case lsub > ls:
		return -1
	case lsub == ls:
		_, err := reader.Read(s)
		if err != nil {
			return -1
		}
		if string(s) == substr {
			return 0
		}
		return -1
	}
	i := 0
	lastRead := 0
	for i+lsub-1 < ls {
		j := lsub - 1
		if i+lsub > lastRead {
			bytesRead, err := reader.Read(readBuffer)
			if err != nil {
				return -1
			}
			if lastRead+lsub > ls {
				copy(s[lastRead:ls], readBuffer[:bytesRead])
			} else {
				copy(s[lastRead:lastRead+lsub], readBuffer)
			}
			lastRead += lsub
		}
		for ; j >= 0 && s[i+j] == substr[j]; j-- {
		}
		if j < 0 {
			return i
		}

		slid := j - d[s[i+j]]
		if slid < 1 {
			slid = 1
		}
		i += slid
	}
	return -1
}

func calculateSlideTable(substr string) [256]int {
	var d [256]int
	for i := 0; i < 256; i++ {
		d[i] = -1
	}
	for i := 0; i < len(substr); i++ {
		d[substr[i]] = i
	}
	return d
}
