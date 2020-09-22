const confirm_templateFn = message => `
<div class="profile-block">
<li class="message">
    <div class="row space-between">
        <div class="column">
            <div class="row">
                <p class="profile-name">
                    ${message['name']}
                </p>
                <span class="ml">-</span>
                <p class="ml">${moment(message['created']).format('YYYY-MM-DD HH:mm')}</p>
                <div class="likes">
                <img src="pictures/white_heart.png" class="likes-icon"/>
                <p class="likes-number">${message['likes']}</p>
                </div>
                <div class="like" data-id="${message.id}">
                <img src="pictures/like.png" class="like-icon"/>
                <span class="like-text">Gilla</span>
                </div>
            </div>
            <h2>
                ${message['title']}
            </h2>
            <p class="row">${message['message']}</p>
            <p>${getTimeCheckboxLabels(message)}</p>
            <p>${getAgeCheckboxLabels(message)}</p>    
        </div>
        <img src="${message['image']}" class="profile-image"/>
    </div>
</li>
</div>`

const submitMessage = e => {
    e.preventDefault()

    if (document.getElementById('title').value == '') { alert('missing title'); return }
    if (document.getElementById('name').value == '') { alert('missing name'); return }
    if (document.getElementById('email').value == '') { alert('missing email'); return }

    const FD = new FormData(document.getElementById('message-form'))
    const message = Object.fromEntries(FD)

    // For checkboxes, value 'on' is set to 'true', and 'off' is set to 'false', to be readable for Backend.
    message.vardag_fm = message.vardag_fm == 'on'
    message.vardag_em = message.vardag_em == 'on'
    message.vardag_kvall = message.vardag_kvall == 'on'
    message.helg = message.helg == 'on'
    message.age_0_6 = message.age_0_6 == 'on'
    message.age_7_12 = message.age_7_12 == 'on'
    message.age_13_18 = message.age_13_18 == 'on'

    server.addMessage(message).then((response) => {
        const elem = document.getElementById('registration')
        const html = confirm_templateFn(response)
        elem.innerHTML = '<p>Tack för registreringen. Din profil är sparad.</p>' + html
    })
}

if (document.getElementById('submit-button')) {
    document.getElementById('submit-button').onclick = submitMessage
}

