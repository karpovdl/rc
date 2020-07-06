# Samples

## Register

POST /api/register - registering a new user and getting a JWT token

Sample request

```http
curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"user.1\",\"password\":\"user.pass.1\"}" http://localhost:8080/api/register
```

Sample response

```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE"}
```

```json
{"errors":[{"location":"body","param":"username","value":"user.1","msg":"already exists"}]}
```

## Login

POST /api/login - log in as an existing user and get a JWT token

Sample request

```http
curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"user.1\",\"password\":\"user.pass.1\"}" http://localhost:8080/api/login
```

Sample response

```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE"}
```

```json
{"message":"user not found"}
```

```json
{"message":"invalid password"}
```

## Post

### Add new post

POST /api/posts/ - adding a post with url or text

Sample request

```http
curl -X POST -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" -d "{\"category\":\"music\",\"type\":\"link\",\"title\":\"title_link\",\"url\":\"https://a.b.url\"}" http://localhost:8080/api/posts
```

Sample response

```json
{"id":"5f02e0e90bdfcf000103c913","category":"music","type":"link","title":"title_link","url":"https://a.b.url","author":{"id":"5f0241eb0bdfcf0001558375","username":"user.1"},"comments":[],"created":"2020-07-06T08:29:29.494Z","scope":1,"views":0,"upvotePercentage":100,"votes":[{"user":"5f0241eb0bdfcf0001558375","vote":1}]}
```

Sample request

```http
curl -X POST -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" -d "{\"category\":\"music\",\"type\":\"text\",\"title\":\"title_text\",\"text\":\"a.b.message\"}" http://localhost:8080/api/posts
```

Sample response

```json
{"id":"5f02e8890bdfcf000103c914","category":"music","type":"text","title":"title_text","text":"a.b.message","author":{"id":"5f0241eb0bdfcf0001558375","username":"user.1"},"comments":[],"created":"2020-07-06T09:02:01.504Z","scope":1,"views":0,"upvotePercentage":100,"votes":[{"user":"5f0241eb0bdfcf0001558375","vote":1}]}
```

### List posts

GET /api/posts/ - get all posts

Sample request

```http
curl -X GET -H "Content-Type: application/json" "http://localhost:8080/api/posts"
```

GET /api/user/{USER_LOGIN} - get all post by user

Sample request

```http
curl -X GET -H "Content-Type: application/json" "http://localhost:8080/api/user/user.1"
```

GET /api/post/{POST_ID} - get post by id

Sample request

```http
curl -X GET -H "Content-Type: application/json" "http://localhost:8080/api/posts/5f02435b0bdfcf0001558377"
```

GET /a/funny/{CATEGORY_NAME} - get all post by category

Sample request

```http
curl -X GET -H "Content-Type: application/json" "http://localhost:8080/a/funny/music"
```

### Vote

GET /api/post/{POST_ID}/upvote - rating post up vote by id

Sample request

```http
curl -X GET -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" "http://localhost:8080/api/posts/5f02e0d40bdfcf000103c912/upvote"
```

GET /api/post/{POST_ID}/downvote - rating post down vote by id

Sample request

```http
curl -X GET -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" "http://localhost:8080/api/posts/5f02e0d40bdfcf000103c912/downvote"
```

### Comment

POST /api/post/{POST_ID} - add comment to exists post by id

Sample request

```http
curl -X POST -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" -d "{\"comment\":\"a.b.comment\"}" "http://localhost:8080/api/posts/5f02e0d40bdfcf000103c912"
```

DELETE /api/post/{POST_ID}/{COMMENT_ID} - delete comment from exists post by id

Sample request

```http
curl -X DELETE -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" "http://localhost:8080/api/posts/5f02e0d40bdfcf000103c912/5f02edb50bdfcf000103c918"
```

### Delete post

DELETE /api/post/{POST_ID} - delete post by id

Sample request

```http
curl -X DELETE -H "Content-Type: application/json" -H "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQxMDg0NzEsImlhdCI6MTU5NDAyMjA3MSwidXNlciI6eyJpZCI6IjVmMDI0MWViMGJkZmNmMDAwMTU1ODM3NSIsInVzZXJuYW1lIjoidXNlci4xIn19.JDsTW5ywEc9avWJd-OfhvG0DvoqiRfWY7Lr2P2FQDhE" "http://localhost:8080/api/posts/5f02435b0bdfcf0001558377"
```
