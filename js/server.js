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
