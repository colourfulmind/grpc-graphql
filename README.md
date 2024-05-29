### SaveUser
mutation {
    SaveUser(input: {email: "test@test.com", password: "1234567890"}) {
        user_id
    }
}

### Login
mutation {
    Login(input: {email: "test2@test.com", password: "1234567890"}) {
        token
    }
}

### SavePost
mutation {
    SavePost(input: {token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiZXhwIjoxNzE2OTg4MDkxLCJ1aWQiOjJ9.BaC3zfiz8M9K4z_IF1v13OY_X5yyh5SQ30bv_46Av30",
    title: "test2", content: "hello, world", comments: true}) {
        post_id
        created_at
    }
}

### ProvidePost
mutation {
    ProvidePost(input: {post_id: 3}) {
        id
        user_id
        title
        content
        created_at
        comments
    }
}

### ProvideAllPosts
mutation {
    ProvideAllPosts(input: {page: 1}) {
        posts {
        id
            user_id
            title
            content
            created_at
            comments
        }
    }
}

### SaveComment
mutation {
    SaveComment(input: {token:"", post_id: 1, content: "hello"}) {
        id
        created_at
    }
}

### SaveCommentToComment
mutation {
    SaveCommentToComment(input: {token:"2", post_id: 1, parent_id: 1, content: "hello"}) {
        id
        created_at
    }
}


### ProvideComment
mutation {
    ProvideComment(input: {post_id: 1}) {
        comments {
            id
            user_id
            post_id
            content
            created_at
            comment_id
        }
    }
}