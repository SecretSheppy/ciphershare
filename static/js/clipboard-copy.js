'use strict'

/**
 * displays the success message in the DOM
 */
function displaySuccessMessage() {
    $('#copied-success-message').show();

    setTimeout(() => {
        $('#copied-success-message').fadeOut();
    }, 500);
}

/**
 * Puts the text into the users clipboard
 *
 * @param {string} valueToCopy the string to put into the clipboard
 */
function putTextInClipboard(valueToCopy) {
    navigator.clipboard.writeText(valueToCopy)
        .then(r => {
            console.log('Text copied to clipboard:', valueToCopy);
            displaySuccessMessage();
        }).catch(e => console.error('failed to copy text:', e));
}

$(() => {
    $('#copy-to-clipboard').on( 'click', ()=> {
        putTextInClipboard($('#link-to-copy').attr('href'));
    });

});