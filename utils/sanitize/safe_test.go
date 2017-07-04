package sanitize

import "testing"

func TestSafe(t *testing.T) {
	safeString := map[string]string{
		"'":                  "&#39;",
		"&":                  "&amp;",
		"http://exemple.com": "http://exemple.com",
	}
	for key, val := range safeString {
		safe := Safe(key)
		if string(safe) != val {
			t.Errorf("Safe doesn't escape the right values, expected result %s, got %s", key, val)
		}
	}
}

func TestSafeText(t *testing.T) {
	safeString := map[string]string{
		"'":                                    "&#39;",
		"&":                                    "&amp;",
		"http://exemple.com":                   "http://exemple.com",
		"<em>test</em><script>lol();</script>": "&lt;em&gt;test&lt;/em&gt;&lt;script&gt;lol();&lt;/script&gt;",
	}
	for key, val := range safeString {
		safe := Safe(key)
		if string(safe) != val {
			t.Errorf("Safe doesn't escape the right values, expected result %s, got %s", key, val)
		}
	}
}
