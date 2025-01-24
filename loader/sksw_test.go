package loader

import (
	"testing"
)

func TestSksw_Book(t *testing.T) {
	sksw := NewSksw()
	title, content, err := sksw.Book("http://www.4ksw.com/49/49383/18339525.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(title)
	t.Log(content)
}

func TestSksw_Category(t *testing.T) {
	sksw := NewSksw()
	title, category, err := sksw.Category("http://www.4ksw.com/49/49383/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(title)
	t.Log(category)
}
