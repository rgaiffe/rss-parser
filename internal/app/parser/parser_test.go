package parser

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetFeedsFromRssParser(t *testing.T) {
	xmlData := `
	<rss>
		<channel>
			<item>
				<title>Test item</title>
				<link>this-is-a-test-link</link>
			</item>
		</channel>
	</rss>`

	reader := ioutil.NopCloser(strings.NewReader(xmlData))
	feed, _ := getFeedsFromRssParser(reader)

	t.Run("get link from items", func(t *testing.T) {
		if feed.Items[0].Link != "this-is-a-test-link" {
			t.Errorf("Expected link to be %s, got %s", "this-is-a-test-link", feed.Items[0].Link)
		}
	})
}
