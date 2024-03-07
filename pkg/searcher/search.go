package searcher

import (
	"bufio"
	"io/fs"
	"sync"
	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func (s *Searcher) lookupInFile(file string, word string) (bool, error) {
	// В задании написано что будет плюсом реализовать поиск с алгоритмической сложностью O(1)
	// Однако я считаю что такая алгоритмическая сложность невозможна для поиска слова в файле
	// Для такой сложности можно было бы использовать поиск в хэш-таблице, но в данном случае
	// это не подходит, так как для построения хэш-таблицы нужно прочитать весь файл,
	// разбить его на слова и добавить в хэш-таблицу, что приведет к сложности O(n)
	// Я пробовал лгоритм Бройера-Мура, однако он не учитывает границы слов
	// Поэтому я реализовал поиск с использованием простого разбиения файла на слова

	f, err := s.FS.Open(file)
	defer f.Close()
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return false, scanner.Err()
		}
		// strings.ToLower(scanner.Text()) == strings.ToLower(word) - для поиска без учета регистра
		if scanner.Text() == word {
			return true, nil
		}
	}
	return false, nil
}

func (s *Searcher) lookupInFileAsync(file string, word string, resultChan chan string, errorChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	found, err := s.lookupInFile(file, word)
	if err != nil {
		errorChan <- err
	}
	if found {
		resultChan <- file
	}
}

func (s *Searcher) Search(word string) (files []string, err error) {
	files, err = dir.FilesFS(s.FS, ".")
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	fileListChan := make(chan string, len(files))
	errorChan := make(chan error, len(files))
	const chunkSize = 10
	chunkedFiles := chunkBy(files, chunkSize)
	for _, chunk := range chunkedFiles {
		for _, file := range chunk {
			wg.Add(1)
			go s.lookupInFileAsync(file, word, fileListChan, errorChan, &wg)
		}
		wg.Wait()
	}
	close(fileListChan)
	close(errorChan)
	for err := range errorChan {
		return nil, err
	}
	filesContainigWord := make([]string, 0)
	for file := range fileListChan {
		filesContainigWord = append(filesContainigWord, file)
	}
	if len(filesContainigWord) == 0 {
		return nil, nil
	}
	return filesContainigWord, err
}
