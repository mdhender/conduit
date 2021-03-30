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

package conduit

// all request types are derived from https://github.com/gothinkster/realworld/blob/9686244365bf5681e27e2e9ea59a4d905d8080db/api/swagger.json

type ArticleCreateRequest struct {
	Article struct {
		Body        string   `json:"body"`        // "body": "You have to believe" // required
		Description string   `json:"description"` // "description": "Ever wonder how?" // required
		TagList     []string `json:"tagList"`     // "tagList": ["reactjs", "angularjs", "dragons"] // optional
		Title       string   `json:"title"`       // "title": "How to train your dragon" // required
	} `json:"article"`
}

type ArticleUpdateRequest struct {
	Article struct {
		Body        string `json:"body"`        // "body": "You have to believe" // optional
		Description string `json:"description"` // "description": "Ever wonder how?" // optional
		Title       string `json:"title"`       // "title": "How to train your dragon" // optional
	} `json:"article"`
}

type CommentAddRequest struct {
	Comment struct {
		Body string `json:"body"` // "body": "His name was my name too." // required
	} `json:"comment"`
}

type LoginUserRequest struct {
	User LoginUser `json:"user"`
}

type LoginUser struct {
	Email    string `json:"string"`   // "email": "jake@jake.jake" // required
	Password string `json:"password"` // "password": "jakejake" // required
}

type NewUserRequest struct {
	User NewUser `json:"user"` // required
}

type NewUser struct {
	Email    string `json:"email"`    // "email": "jake@jake.jake" // required
	Password string `json:"password"` // "password": "jakejake" // required
	Username string `json:"username"` // "username": "Jacob" // required
}

type UpdateUserRequest struct {
	User UpdateUser `json:"user"`
}

type UpdateUser struct {
	Email    *string `json:"email"`    // "email": "jake@jake.jake" // optional and nullable
	Token    string  `json:"token"`    // "token": "...." // required but ignored
	Username string  `json:"username"` // "username": "Jacob" // optional and ignored
	Bio      *string `json:"bio"`      // "bio": "I like to skateboard" // optional and nullable
	Image    *string `json:"image"`    // "image": "https://i.stack.imgur.com/xHWG8.jpg" // optional and nullable
}
