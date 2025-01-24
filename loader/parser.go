package loader

type Parser interface {
	Category(link string) (title string, category []string, err error)
	Book(link string) (title, content string, err error)
}

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0.1 Safari/605.1.15"
