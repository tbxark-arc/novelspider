package loader

import (
	"testing"
)

func TestSksw_Book(t *testing.T) {
	sksw := NewSksw("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	title, content, err := sksw.Book("http://www.4ksw.com/49/49383/18339525.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(title)
	t.Log(content)
}

func TestSksw_Category(t *testing.T) {
	sksw := NewSksw("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	title, category, err := sksw.Category("http://www.4ksw.com/49/49383/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(title)
	t.Log(category)
}
