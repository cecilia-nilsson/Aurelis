const templateFn = message => `
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
    <p>
        <span class="show-comment-form link" data-id="${message.id}">Skriv en kommentar</span>
         eller skicka ett mail till <a class="link" href="mailto:${message['email']}"">
        ${message['name']}</a>.
    </p>
    <p class="show-comments-list link" data-id="${message.id}">Visa alla kommentarer</p>
    <button class="delete-post" data-id="${message.id}">Radera profilen</button>
    <div class="hide" id="comments-list${message.id}" data-id="${message.id}">
        <p>${message['comment'].map(comment => `
            <p class="comment-text-output">${comment.message}</p>
            <div class="row">
                <p class="comment-name-output">${comment.name}</p>
                <span class="comment-date-output">-</span>
                <p class="comment-date-output">${moment(comment.created).format('YYYY-MM-DD HH:mm')}</p>
            </div>
        `).join('')}</p>
    </div>
    <div class="hide" id="comment-form${message.id}" data-id="${message.id}">
        <h2>Skriv en kommentar</h2>
        <ul id="comment-form">
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
</li>
</div>`

{/* <button class="update-post" data-id="${message.id}">Uppdatera profilen</button> */ }

// const updateFn = message => `
// <h1>Uppdatera profilen</h1>
// <div id="update-profile">
//     <form id="update-message-form">
//         <div>
//             <label>Titel:</label>
//             <input type="text" id="title" name="title" class="inputarea" value=${message['title']} maxlength="40" required>
//         </div>
//         <div>
//             <label>Namn:</label>
//             <input type="text" id="name" name="name" class="inputarea" value=${message['name']} maxlength="40" required>
//         </div>
//         <div class="email-field">
//             <label>E-postadress:</label>
//             <input type="email" id="email" name="email" class="inputarea" value=${message['email']} required>
//         </div>
//         <div>
//             <label>Bild:</label>
//             <input type=text name="image" class="inputarea" value=${message['image']}></input>
//         </div>
//         <div>
//             <label>Meddelande:</label>
//             <textarea name="message" class="inputarea" value=${message['message']}></textarea>
//         </div>
//         <div><label>Hjälp behövs främst:</label>
//             <ul class="checkboxes">
//                 <li>
//                     <input type="checkbox" name="vardag_fm" value=${message['vardag_fm']}>
//                     <label>Vardagar förmiddag</label>
//                 </li>
//                 <li>
//                     <input type="checkbox" name="vardag_em" value=${message['vardag_em']}>
//                     <label>Vardagar eftermiddag</label>
//                 </li>
//                 <li>
//                     <input type="checkbox" name="vardag_kvall" value=${message['vardag_kvall']}>
//                     <label>Vardagar kväll</label>
//                 </li>
//                 <li>
//                     <input type="checkbox" name="helg" value=${message['helg']}>
//                     <label>Helger</label>
//                 </li>
//             </ul>
//         </div>
//         <div><label>Ålder på barnet/barnen:</label>
//             <ul class="checkboxes">
//                 <li>
//                     <input type="checkbox" name="age_0_6" value=${message['age_0_6']}>
//                     <label>0-6 år</label>
//                 </li>
//                 <li>
//                     <input type="checkbox" name="age_7_12" value=${message['age_7_12']}>
//                     <label>7-12 år</label>
//                 </li>
//                 <li>
//                     <input type="checkbox" name="age_13_18" value=${message['age_13_18']}>
//                     <label>13-18 år</label>
//                 </li>
//             </ul>
//         </div>
//         <button id="update-profile-button">Spara ändringarna</button>
//     </form>
// </div>
// `

const render = messages => {
    // messages.forEach(message => { console.log(message.image) })
    if (messages.length === 0) {
        if (document.getElementById("messages-list")) {
            document.getElementById("messages-list").innerHTML =
                '<p>Det finns tyvärr inga profiler som motsvarade din sökning. Prova att använda färre sökkriterier.</p>'
        }
    }
    else {
        const mappedMessages = messages.map(templateFn)
        const html = mappedMessages.join('')
        if (document.getElementById("messages-list")) {
            document.getElementById("messages-list").innerHTML = html
        }
    }

    for (let button of document.querySelectorAll('.message .like')) {
        button.onclick = () => {
            server.likePost(button.dataset.id).then(() => {
                server.listMessages().then(render)
            })
        }
    }

    for (let button of document.querySelectorAll('.message .delete-post')) {
        button.onclick = () => {
            server.deletePost(button.dataset.id).then(() => {
                server.listMessages().then(render)
            })
        }
    }

    // for (let button of document.querySelectorAll('.message .update-post')) {
    //     button.onclick = () => {
    //         message = server.getMessage(button.dataset.id)
    //         console.log("id: ", button.dataset.id)
    //         console.log(message)
    //         const elem = document.getElementById('search-form')
    //         const html = updateFn(message)
    //         elem.innerHTML = html
    //     }
    // }

    // for (let button of document.querySelectorAll('.message .update-profile-button')) {
    //     button.onclick = () => {
    //         const FD = new FormData(document.getElementById('update-message-form'))
    //         const message = Object.fromEntries(FD)
    //         server.updatePost(message, button.dataset.id).then(() => {
    //             server.listMessages().then(render)
    //         })
    //     }
    // }

    for (let button of document.querySelectorAll('.message .submit-comment')) {
        button.onclick = (e) => {
            name = button.previousElementSibling.querySelector('[name="comment.name"]').value
            message = button.previousElementSibling.querySelector('[name="comment.message"]').value
            server.addComment(button.dataset.id, name, message).then(() => server.listMessages().then(render))
            const elem = document.getElementById('comment-form')
            elem.innerHTML = '<p>Kommentaren är nu publicerad.</p>'
        }
    }

    for (let button of document.querySelectorAll('.message .show-comment-form')) {
        button.onclick = (e) => {
            let profile_comment_form = document.getElementById(`comment-form${button.dataset.id}`);
            // console.log('Profile_comments dataset.id: ', button.dataset.id)
            // console.log('Profile_comments värde: ', profile_comments)
            profile_comment_form.classList.remove("hide");
        }
    }

    for (let button of document.querySelectorAll('.message .show-comments-list')) {
        button.onclick = (e) => {
            let profile_comments_list = document.getElementById(`comments-list${button.dataset.id}`);
            // console.log('Profile_comments dataset.id: ', button.dataset.id)
            // console.log('Profile_comments värde: ', profile_comments)
            profile_comments_list.classList.remove("hide");
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
        html = '<h4 class="checkboxheader checkboxheader_time">Hjälp önskas främst följande tider:</h4>' + '<ul class="checkbox-list">' + html + '</ul>'
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
    // server.searchMessages(params).then(messages => {
    //     console.log(messages)
    //     render(messages)
    // })
    server.searchMessages(params).then(render)
}

// server.listMessages().then(render)

if (document.getElementById('search-button')) {
    document.getElementById('search-button').onclick = searchMessages
}