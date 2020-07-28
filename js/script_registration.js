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
    // console.log(message)
//    server.addMessage(message).then(() => server.listMessages().then(render))
    server.addMessage(message).then(() => server.listMessages())
}

if (document.getElementById('submit-button')) {
    document.getElementById('submit-button').onclick = submitMessage
}

