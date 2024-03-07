package searcher

import (
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
)

func TestSearcher_Search(t *testing.T) {
	type fields struct {
		FS fs.FS
	}
	type args struct {
		word string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "Ok",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "World"},
			wantFiles: []string{"file1.txt", "file3.txt"},
			wantErr:   false,
		},
		{
			name: "NotFound",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "World2"},
			wantFiles: nil,
			wantErr:   false,
		},
		{
			name: "Empty",
			fields: fields{
				FS: fstest.MapFS{},
			},
			args:      args{word: "World"},
			wantFiles: nil,
			wantErr:   false,
		},
		{
			name: "Error",
			fields: fields{
				FS: os.DirFS("testdata"), // Тестовая директория не существует
			},
			args:      args{word: "World"},
			wantFiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				FS: tt.fields.FS,
			}
			gotFiles, err := s.Search(tt.args.word)
			sort.Strings(gotFiles) // В задаче не указано что список файл должен быть отсортирован, не сортируем для максимальной производительности
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("Search() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}
