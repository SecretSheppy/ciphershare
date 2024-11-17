# CipherShare
Quickly store and send files with only a link, from a web app. Files are securely encrypted. 
Give people access by list of emails with a simple, easy-to-use UI.


### Development

This project was developed as part of HackSheffield 9.
It allows the user to upload a file and specify a list of recipient email addresses. 
The file is then encrypted and can be retrieved by a recipient confirming their email address using Auth0 OTP(sent to their email)

### Libraries

 - Gorilla web Framework
 - MongoDB
 - Go Resty
 - Go.env
 - llimitman/Tint
