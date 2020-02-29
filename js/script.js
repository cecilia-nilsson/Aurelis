const templateFn = message => `
<li class="message">
    <div class="row space-between">
        <div class="column">
            <div class="row">
                <p>
                    ${message['name']}
                </p>
                <span class="ml">-</span>
                <p class="ml">${message['date']}</p>
                <p class="likes">Gillningar: ${message['likes']}</p>
            </div>
            <h3>
                ${message['title']}
            </h3>
            <p class="row">${message['text']}</p>
            <p>
        </div>
        <img src="${message['image']}" class="profile-image">
    </div>
        Skriv en kommentar nedan eller skicka ett mail till <a href="${message['email']}"">
        ${message['name']}</a>.
    </p>
    <button class="like" data-id="${message.id}">Gilla</button>
    <div>
        <h2>Skriv en kommentar</h2>
        <ul id="comment-list">
            <li>
                <label>Kommentar:</label>
                <textarea name="comment.text" class="comment-inputarea"></textarea>
            </li>
            <li>
                <label>Namn:</label>
                <input type="text" name="comment.name" class="comment-inputarea">
            </li>
        </ul>
        <button class="submit-comment" data-id="${message.id}">Skicka kommentar</button>
        <p>${message['comment'].map(comment =>`
            <p class="comment-text-output">${comment.text}</p>
            <div class="row">
                <p class="comment-name-output">${comment.name}</p>
                <span class="comment-date-output">-</span>
                <p class="comment-date-output">${comment.date}</p>
            </div>
        `)}</p>
    </div>
</li>`

const render = () => {
    const mappedMessages = server.listMessages().map(templateFn)
    const html = mappedMessages.join('')
    document.getElementById("messages-list").innerHTML = html

    for(let button of document.querySelectorAll('.message .like')) {
        button.onclick = () => {
            message = server.getMessage(button.dataset.id)
            message.likes++
            server.updateMessage(button.dataset.id, message)
            render()
        }
    }

    for(let button of document.querySelectorAll('.message .submit-comment')) {
        button.onclick = (e) => {
            message = server.getMessage(button.dataset.id)
            //console.error(message.comment)
            //console.log(document.querySelector('[name="comment.text"]').value)
            message.comment.push({
                name: document.querySelector('[name="comment.name"]').value,
                text: document.querySelector('[name="comment.text"]').value,
                // date: new Date(),
                date: moment().format('YYYY-MM-DD  HH:mm'),
            })
            console.error(message.comment)
            server.updateMessage(button.dataset.id, message)

            render()
        }
    }
}

const submitMessage = e => {
    e.preventDefault()

    const FD = new FormData(document.getElementById('message-form'))
    const message = Object.fromEntries(FD)
    message['date'] = moment().format('YYYY-MM-DD  HH:mm')
    message['likes'] = 0
    message['comment'] = []
    server.addMessage(message)

    render()
}

render()

document.getElementById('submit-button').onclick = submitMessage

