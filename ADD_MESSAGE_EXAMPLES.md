# Add Message Examples

You can obtain an application-token from the apps tab inside the UI or using the REST-API (`GET /application`).

NOTE: Assuming Gotify is running on `http://localhost:8008`.

### Python

```python
import requests #pip install requests
resp = requests.post('http://localhost:8008/message?token=<token-from-application>', json={
    "message": "Well hello there.",
    "priority": 2,
    "title": "This is my title"
})
```

### Golang

```go
package main

import (
        "net/http"
        "net/url"
)

func main() {
    http.PostForm("http://localhost:8008/message?<token-from-application>",
        url.Values{"message": {"My Message"}, "title": {"My Title"}})
}
```