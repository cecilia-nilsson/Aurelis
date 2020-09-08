class Server {
    listMessages() {
        return fetch('http://localhost:8080/posts')
        .then(response => response.json())
    }

    addMessage(message) {
        return fetch('http://localhost:8080/posts',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(message),
        })
        .then(response => response.json())
    }

    searchMessages(params) {
        return fetch('http://localhost:8080/searchPosts',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(params),
        })
        .then(response => response.json())
        .catch(err => {console.log(err)})
    }

    getMessage(id) {
        return this.listMessages()[id]
    }

    addComment(id, name, message) {
        return fetch('http://localhost:8080/comments',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({post_id:parseInt(id),name:name,message:message}),
        })
        .then(response => response.json())
    }

    likePost(postId) {
        return fetch(`http://localhost:8080/posts/${postId}/like`,{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(response => response.json())
    }
}

const server = new Server()
