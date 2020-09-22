const templateFn = message => `
<li class="message">
    <div class="row space-between">
        <div class="column">
            <div class="row">
                <p>
                    ${message['name']}
                </p>
                <span class="ml">-</span>
                <p class="ml">${moment(message['date']).format('YYYY-MM-DD  HH:mm')}</p>
                <p class="likes">Gillningar: ${message['likes']}</p>
            </div>
            <h3>
                ${message['title']}
            </h3>
            <p class="row">${message['message']}</p>
        </div>
        <img src="${message['image']}" class="profile-image"/>
    </div>
    <p>${getTimeCheckboxLabels(message)}</p>
    <p>${getAgeCheckboxLabels(message)}</p>
    <p>
        Skriv en kommentar nedan eller skicka ett mail till <a href="mailto:${message['email']}"">
        ${message['name']}</a>.
    </p>
    <button class="like" data-id="${message.id}">Gilla</button>
    <div>
        <p>${message['comment'].map(comment =>`
            <p class="comment-text-output">${comment.message}</p>
            <div class="row">
                <p class="comment-name-output">${comment.name}</p>
                <span class="comment-date-output">-</span>
                <p class="comment-date-output">${moment(comment.date).format('YYYY-MM-DD  HH:mm')}</p>
            </div>
        `).join('')}</p>
        <h2>Skriv en kommentar</h2>
        <ul id="comment-list">
            <li>
                <label>Kommentar:</label>
                <textarea name="comment.message" class="comment-inputarea"></textarea>
            </li>
            <li>
                <label>Namn:</label>
                <input type="text" name="comment.name" class="comment-inputarea" maxlength="40">
            </li>
        </ul>
        <button class="submit-comment" data-id="${message.id}">Skicka kommentar</button>
    </div>
</li>`

const render = messages => {
    // messages.forEach(message => { console.log(message.image) })

    const mappedMessages = messages.map(templateFn)
    const html = mappedMessages.join('')
    if (document.getElementById("messages-list")) {
        document.getElementById("messages-list").innerHTML = html
    }

    for(let button of document.querySelectorAll('.message .like')) {
        button.onclick = () => {
            server.likePost(button.dataset.id).then(() => {
                server.listMessages().then(render)
            })
        }
    }

    for(let button of document.querySelectorAll('.message .submit-comment')) {
        button.onclick = (e) => {
            name = button.previousElementSibling.querySelector('[name="comment.name"]').value
            message = button.previousElementSibling.querySelector('[name="comment.message"]').value
            server.addComment(button.dataset.id,name,message).then(() => server.listMessages().then(render))
        }
    }
}

function getTimeCheckboxLabels(message) {
    let html = ''

    let checkboxes = {
        'vardag_fm': 'Vardagar förmiddag',
        'vardag_em': 'Vardagar eftermiddag',
        'vardag_kvall': 'Vardagar kväll',
        'helg': 'Helger',
    }

    for (let value in checkboxes) {
        if (!message[value]) continue
        
        html += `<li>${checkboxes[value]}</li>`
    }

    if (html) {
        html = '<h4 class="checkboxheader">Hjälp önskas främst följande tider:</h4>' + '<ul class="checkbox-list">' + html + '</ul>'
    }

    return html
}

function getAgeCheckboxLabels(message) {
    let html = ''

    let checkboxes = {
        'age_0_6': '0-6 år',
        'age_7_12': '7-12 år',
        'age_13_18': '13-18 år',
    }

    for (let value in checkboxes) {
        if (!message[value]) continue
        
        html += `<li>${checkboxes[value]}</li>`
    }
    if (html) {
        html = '<h4 class="checkboxheader">Ålder på barnet/barnen:</h4>' + '<ul class="checkbox-list">' + html + '</ul>'
    }

    return html
}


const searchMessages = e => {
    e.preventDefault()
    
    const FD = new FormData(document.getElementById('search-form'))
    const params = Object.fromEntries(FD)

    // For checkboxes, value 'on' is set to 'true', and 'off' is set to 'false', to be readable for Backend.
    params.vardag_fm = params.vardag_fm == 'on'
    params.vardag_em = params.vardag_em == 'on'
    params.vardag_kvall = params.vardag_kvall == 'on'
    params.helg = params.helg == 'on'
    params.age_0_6 = params.age_0_6 == 'on'
    params.age_7_12 = params.age_7_12 == 'on'
    params.age_13_18 = params.age_13_18 == 'on'
    // console.log(params)
    server.searchMessages(params).then(messages => {
        console.log(messages)
        render(messages)
    })
}

server.listMessages().then(render)

if (document.getElementById('search-button')) {
    document.getElementById('search-button').onclick = searchMessages
}