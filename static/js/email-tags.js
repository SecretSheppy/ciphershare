'use strict';

// Identifies an email address
const EMAIL_REGEX = /^[\w-.]+@([\w-]+\.)+[\w-]{2,4}/

// Identifies multiple commas in a row with only spaces in between them
const MULTIPLE_COMMAS_IN_ROW = /(, +)+/g

/**
 * Creates a new email tag in the DOM.
 *
 * @param {string} email the email to create the tag for
 * @param {number} index the email
 */
function createNewEmailTag(email, index) {
    $('#email-tags').append(
        $('<div>')
            .addClass('email-tag')
            .attr('id', `email-tag-${index}`)
            .append(
                $('<p>')
                    .addClass('email-text')
                    .text(email))
            .append(
                $('<button>')
                    .addClass('email-tag-removal-button')
                    .attr('type', 'button')
                    .on('click', () => {
                        $(`#email-tag-${index}`).remove();

                        let emails = $(`#emails`);

                        emails.val(emails.val().replaceAll(email, ''))
                        emails.val(emails.val().replaceAll(MULTIPLE_COMMAS_IN_ROW, ', '))
                    }).append(
                        $('<img>')
                            .attr('src','/static/gfx/fontawesome/xmark-solid.svg')
                            .attr('alt','xmark')
                            .addClass('email-tag-removal-icon'))));
}

$(() => {
    $('#emails').on('input', (e) => {
        let validEmails = $('#emails').val()

        if (validEmails.length === 0) {
            return
        }

        const emails = validEmails.split(',');

        $('#email-tags').empty();

        emails.forEach((value, index) => {
            let email = value.trim();

            if (!EMAIL_REGEX.test(email)) {
                return;
            }

            createNewEmailTag(email, index);
        });
    });
});
