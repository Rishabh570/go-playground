package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"shortener/internal/redisPkg"
	"strings"
)

type RequestBody struct {
	URL string
}

var shortLinkMapping = make(map[string]string)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var ctx = context.Background()

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func HandleShortenURL(clientConn *redisPkg.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var urls []RequestBody

		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, url := range urls {
			fmt.Println("url:", url.URL)
			slug := RandSeq(10)
			fmt.Println("slug created:", slug)

			// store mapping in redis
			// shortLinkMapping[slug] = url.URL
			err := clientConn.Client.Set(ctx, slug, url.URL, 0).Err()
			if err != nil {
				panic(err)
			}
		}

		for ok, mapping := range shortLinkMapping {
			fmt.Println("mapping", ok, mapping)
		}

		w.Write([]byte("HELLO"))
	})
}

func HandleRedirection(clientConn *redisPkg.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Split the request URL path by '/'
		parts := strings.Split(r.URL.Path, "/")
		fmt.Println("parts:", parts)
		fmt.Println("parts[0]:", parts[0])

		for ok, mapping := range shortLinkMapping {
			fmt.Println("mapping", ok, mapping)
		}

		for _, slug := range parts[1:] {
			fmt.Println("slug:", slug)

			val, err := clientConn.Client.Get(ctx, slug).Result()
			if err != nil {
				panic(err)
			}

			// mappingFound, ok := shortLinkMapping[slug]
			// if !ok {
			// 	fmt.Printf("no mapping present for %s", mappingFound)
			// 	return
			// }

			fmt.Println("slug:", slug)
			http.Redirect(w, r, val, http.StatusTemporaryRedirect)
		}
	})
}
