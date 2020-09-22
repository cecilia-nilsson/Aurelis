const accountMessage = e => {
    e.preventDefault()

    // if (document.getElementById('account-email').value == '') { alert('missing email'); return }
    // if (document.getElementById('account-password').value == '') { alert('missing password'); return }

    const FD = new FormData(document.getElementById('account-form'))
    const message = Object.fromEntries(FD)

    server.addUser(message).then((response) => {
        const elem = document.getElementById('register_user')
        elem.innerHTML = '<p>Tack för registreringen. Ditt konto är nu skapat och du kan logga in. Detta är dock ej bekräftat av servern. Behöver fixas.</p>'
    })
}

const loginMessage = e => {
    e.preventDefault()

    // if (document.getElementById('email').value == '') { alert('missing email'); return }
    // if (document.getElementById('password').value == '') { alert('missing password'); return }

    const FD = new FormData(document.getElementById('login-form'))
    const message = Object.fromEntries(FD)

    server.login(message).then((response) => {
        const elem = document.getElementById('login')
        elem.innerHTML = '<p>Inloggningen är slutförd. Det syns dock bara i terminalen om inloggningen funkade eller ej. Detta behöver ändras.</p>'
    })
}

if (document.getElementById('account-button')) {
    document.getElementById('account-button').onclick = accountMessage
}

if (document.getElementById('login-button')) {
    document.getElementById('login-button').onclick = loginMessage
}