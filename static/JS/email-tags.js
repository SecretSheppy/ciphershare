'use strict';

$(() => {
    //$('#emails').val()
    $('#emails').on('input', (e) => {
        let validEmails = $('#emails').val()
        console.log(validEmails);
        if (validEmails.length === 0) {
            return
        }
        const emails = validEmails.split(',');
        $('#email-tags').empty();

        emails.forEach((value, index) => {
            let email = value.trim();
            let emailRe = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}/
            console.log(emailRe.test(email));
            let isEmail = emailRe.test(email) ? emailRe.test(email) : false;
            if (!isEmail) {
                return;
            }
            

            //TODO create blocks
            $('#email-tags').append(
                $('<div>').addClass('email-tag').append(
                    $('p').addClass('email-text').text(email)
                ).append(
                    $('button').addClass('email-tag-removal-button').append(
                        $('img').attr('src','/static/gfx/fontawesome/xmark-solid.svg').attr('alt','xmark')
                    )
                )
            );

            //deleteButton.bind("click", (e) => {
            //    deleteButton.parent.remove();
            //});
        });
    });


})
/*
<div class="email-tag">
                    <p class="email-text">email@example.com</p>
                    <button class="email-tag-removal-button">
                        <img src="/static/gfx/fontawesome/xmark-solid.svg" alt="xmark" />
                    </button>
                </div>
 */