# ozonTask
Проект представляет собой сервис публикации постов и написания комментариев к ним. Хранилище предоставляется на выбор при запуске программы - In-memory или PostgreSQL.


# Инструкция по запуску
```shell
#Склонировать репозиторий и перейти в рабочую директорию
https://github.com/timurdilek/ozonTask
cd ozonTask
#Выбор PostgreSQL в качестве хранилища и Docker в качестве инструмента для контейнеризации 
make postgres-docker
#Выбор In-memory в качестве хранилища и Docker в качестве инструмента для контейнеризации 
make in-memory-docker
#Выбор PostgreSQL в качестве хранилища, и запуск на локальной машине
make postgres-local
#Выбор In-memory в качестве хранилища, и запуск на локальной машине
make in-memory-local
```
# Запросы обрабатываемые сервисом
`GET /graphql?query={getPost(first:5){id,content,createdAt}}`<br/>
`GET /graphql?query={getPostById(id:"1"){id,authorId,content,areCommentsAllowed,comments{id,content}}}`<br/>
`GET /graphql?query={getCommentByPostId(postId:"1",first:10){id,authorId,content,replies{id,content}}}`<br/>
`GET /graphql?query={getCommentByParentCommentId(parentCommentId:"2",first:5){id,authorId,content}}`<br/>
`POST /graphql body {"query":"mutation { createPost(input: { authorId: \"1\", content: \"Hello, world!\", areCommentsAllowed: true }) { id content createdAt } }"}`<br/>
`POST /graphql body {"query":"mutation { postComment(input: { postId: \"1\", authorId: \"2\", content: \"Great post!\" }) { id content createdAt } }"}`<br/>
`POST /graphql body {"query":"mutation { postComment(input: { postId: \"1\", parentCommentId: \"3\", authorId: \"4\", content: \"I agree!\" }) { id content createdAt } }"}`<br/>
`POST /graphql body {"query":"mutation { putPost(input: { id: \"1\", content: \"Updated content\", areCommentsAllowed: false }) { id content areCommentsAllowed } }"}`<br/>
`POST /graphql body {"query":"mutation { putComment(input: { id: \"2\", content: \"Updated comment text\" }) { id content updatedAt } }"}`<br/>
`POST /graphql body {"query":"mutation { deletePost(id: \"1\") }"}`<br/>
`POST /graphql body {"query":"mutation { deleteComment(id: \"2\") }"}`<br/>
`POST /graphql body {"query":"subscription { subscriptionForComment(postId: \"1\") { id content authorId createdAt } }"}`<br/>

# Примеры запросов
```shell
# Получить список постов
curl -X GET "http://localhost:8080/graphql?query={getPost(first:5){id,content,createdAt}}"

# Получить пост по ID
curl -X GET "http://localhost:8080/graphql?query={getPostById(id:\"1\"){id,authorId,content,areCommentsAllowed,comments{id,content}}}"

# Получить комментарии для поста
curl -X GET "http://localhost:8080/graphql?query={getCommentByPostId(postId:\"1\",first:10){id,authorId,content,replies{id,content}}}"

# Получить ответы на комментарий
curl -X GET "http://localhost:8080/graphql?query={getCommentByParentCommentId(parentCommentId:\"2\",first:5){id,authorId,content}}"

# Создать новый пост
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { createPost(input: { authorId: \"1\", content: \"Hello, world!\", areCommentsAllowed: true }) { id content createdAt } }"}'

# Добавить комментарий к посту
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { postComment(input: { postId: \"1\", authorId: \"2\", content: \"Great post!\" }) { id content createdAt } }"}'

# Ответить на комментарий
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { postComment(input: { postId: \"1\", parentCommentId: \"3\", authorId: \"4\", content: \"I agree!\" }) { id content createdAt } }"}'

# Обновить пост
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { putPost(input: { id: \"1\", content: \"Updated content\", areCommentsAllowed: false }) { id content areCommentsAllowed } }"}'

# Обновить комментарий
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { putComment(input: { id: \"2\", content: \"Updated comment text\" }) { id content updatedAt } }"}'

# Удалить пост
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { deletePost(id: \"1\") }"}'

# Удалить комментарий
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { deleteComment(id: \"2\") }"}'

# Подписаться на новые комментарии к посту
curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -d '{"query":"subscription { subscriptionForComment(postId: \"1\") { id content authorId createdAt } }"}'
```
