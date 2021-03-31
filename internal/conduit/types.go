/*
 * conduit - current practices for Go web servers
 *
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package conduit defines the types exposed to the client by the web server.
package conduit

// all types are derived from https://github.com/gothinkster/realworld/blob/9686244365bf5681e27e2e9ea59a4d905d8080db/api/swagger.json

type Article struct {
	Slug           string `json:"slug"`           // "slug": "how-to-train-your-dragon"
	Title          string `json:"title"`          // "title": "How to train your dragon"
	Description    string `json:"description"`    // "description": "Ever wonder how?"
	Body           string `json:"body"`           // "body": "It takes a Jacobian"
	TagList        string `json:"tagList"`        // "tagList": ["dragons", "training"]
	CreatedAt      string `json:"createdAt"`      // "createdAt": "2016-02-18T03:22:56.637Z"
	UpdatedAt      string `json:"updatedAt"`      // "updatedAt": "2016-02-18T03:48:35.824Z"
	Favorited      bool   `json:"favorited"`      // "favorited": false
	FavoritesCount int    `json:"favoritesCount"` // "favoritesCount": 0
	Author         Author `json:"author"`
}

type Author struct {
	Username  string `json:"username"`  // "username": "jake"
	Bio       string `json:"bio"`       // "bio": "I work at statefarm"
	Image     string `json:"image"`     // "image": "https://i.stack.imgur.com/xHWG8.jpg"
	Following bool   `json:"following"` // "following": false
}

type Comment struct {
	Id        int    `json:"id"`        // "id": 1,
	CreatedAt string `json:"createdAt"` // "createdAt": "2016-02-18T03:22:56.637Z",
	UpdatedAt string `json:"updatedAt"` // "updatedAt": "2016-02-18T03:22:56.637Z",
	Body      string `json:"body"`      // "body": "It takes a Jacobian",
	Author    Author `json:"author"`
}

type Profile struct {
	Username  string  `json:"username"`  // "username": "jake"
	Bio       *string `json:"bio"`       // "bio": "I work at statefarm"
	Image     *string `json:"image"`     // "image": "https://static.productionready.io/images/smiley-cyrus.jpg"
	Following bool    `json:"following"` // "following": false
}

type User struct {
	Id        int     `json:"id"`
	Email     string  `json:"email"`           // "email": "jake@jake.jake"
	Username  string  `json:"username"`        // "username": "jake"
	CreatedAt string  `json:"createdAt"`       // "createdAt": "2016-02-18T03:22:56.637Z",
	UpdatedAt string  `json:"updatedAt"`       // "createdAt": "2016-02-18T03:22:56.637Z",
	Token     string  `json:"token,omitempty"` // "token": "jwt.token.here"
	Bio       *string `json:"bio"`             // "bio": "I work at statefarm" // API requires this to be nullable
	Image     *string `json:"image"`           // "image": null // API requires this to be nullable
}
