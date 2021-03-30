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

// all response types derived are from https://github.com/gothinkster/realworld/blob/9686244365bf5681e27e2e9ea59a4d905d8080db/api/swagger.json

type ArticleResponse struct {
	Article Article `json:"article"`
}

type MultipleArticlesResponse struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"` // "articlesCount": 2
}

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type GenericErrorModel struct {
	Errors ErrorResponse `json:"errors"`
}

// If a request fails any validations, expect a 422 and errors in the following format:
type ErrorResponse struct {
	Body     []string `json:"body,omitempty"`
	Email    []string `json:"email,omitempty"`
	Password []string `json:"password,omitempty"`
	Username []string `json:"username,omitempty"`
}

type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

type TagsResponse struct {
	Tags []string `json:"tags"`
}

type UserResponse struct {
	User User `json:"user"`
}
