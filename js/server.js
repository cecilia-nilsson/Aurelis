class Server {
    listMessages() {
        return JSON.parse(localStorage['messages'] ? localStorage['messages'] : '[]')
    }

    getMessage(id) {
        return this.listMessages()[id]
    }

    addMessage(message) {
        const messages = this.listMessages()
        message['id'] = messages.length
        messages.push(message)
        localStorage['messages'] = JSON.stringify(messages)
    }

    updateMessage(id, message) {
        const messages = this.listMessages()
        messages[id] = message
        localStorage['messages'] = JSON.stringify(messages)
    }
}

const server = new Server()
