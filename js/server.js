class Server {
    listMessages() {
        return fetch('http://localhost:8080/posts')
        .then(response => response.json())
    }

    getMessage(id) {
        return this.listMessages()[id]
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

    updateMessage(id, message) {
        const messages = this.listMessages()
        messages[id] = message
        localStorage['messages'] = JSON.stringify(messages)
    }
}

const server = new Server()
