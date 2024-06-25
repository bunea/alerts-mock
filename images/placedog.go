package images

import (
	"fmt"
	"sync"
)

const (
	maxPlaceholderID = 238
	urlFormat        = "https://placedog.net/800/640?id=%d"
)

var (
	currentPlaceholderID = 1
	placeholders         = map[string]string{}
	lock                 sync.Mutex
)

func GetImageURL(uid string) string {
	lock.Lock()
	defer lock.Unlock()

	url, found := placeholders[uid]
	if !found {
		url = fmt.Sprintf(urlFormat, currentPlaceholderID)
		placeholders[uid] = url

		currentPlaceholderID++
		if currentPlaceholderID > maxPlaceholderID {
			currentPlaceholderID = 1
		}
	}

	return url
}
