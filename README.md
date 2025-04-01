# ozonTask
Проект представляет собой сервис публикации постов и написания комментариев к ним. Хранилище предоставляется на выбор при запуске программы - In-memory или PostgreSQL.


# Инструкция по запуску
```shell
#Склонировать репозиторий и перейти в рабочую директорию
https://github.com/timurdilek/ozonTask
cd ozonTask
#Выбор PostgreSQL в качестве хранилища и Docker в качестве инструмента для контейнеризации 
make postgres
#Выбор In-memory в качестве хранилища и Docker в качестве инструмента для контейнеризации 
make in-memory
```
# Запросы обрабатываемые сервисом
```graphql
query GetPost {
    getPost(first: 0) {
        id
        authorId
        content
        areCommentsAllowed
        createdAt
        updatedAt
    }
    getPostById(id: "1") {
        id
        authorId
        content
        areCommentsAllowed
        createdAt
        updatedAt
    }
    getCommentByPostId(postId: "1", first: 0) {
        id
        postId
        parentCommentId
        authorId
        content
        createdAt
        updatedAt
    }
    getCommentByParentCommentId(parentCommentId: "1", first: 0) {
        id
        postId
        parentCommentId
        authorId
        content
        createdAt
        updatedAt
    }
}
```

# Мутации
```graphql
    mutation CreatePost {
        createPost(input: { authorId: "1", content: "test", areCommentsAllowed: true }) {
            id
            authorId
            content
            areCommentsAllowed
            createdAt
            updatedAt
        }
        postComment(input: { postId: "1", authorId: "1", content: "test" }) {
            id
            postId
            parentCommentId
            authorId
            content
            createdAt
            updatedAt
        }
        putPost(input: { id: "1", content: "test2" }) {
            id
            authorId
            content
            areCommentsAllowed
            createdAt
            updatedAt
        }
        putComment(input: { id: "1", content: "test2" }) {
            id
            postId
            parentCommentId
            authorId
            content
            createdAt
            updatedAt
        }
        deletePost(id: "1")
        deleteComment(id: "1")
    }
```

### Subscription
```graphql
subscription SubscriptionForComment {
    subscriptionForComment(postId: "1") {
        id
        postId
        parentCommentId
        authorId
        content
        createdAt
        updatedAt
    }
}
```
