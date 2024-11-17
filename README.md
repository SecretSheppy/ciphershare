# CipherShare
Quickly store and send files with only a link, from a web app. Files are securely encrypted. 
Give people access by list of emails with a simple, easy-to-use UI.

## Table of Contents
- [Features](#features)
- [Libraries](#libraries)
- [Setup](#setup)
  
---

### Features
The index page of the app allows a user to upload a file and specify a list of recipient email addresses. 
The file is then encrypted and stored, and can be retrieved by anyone with an email in the recipients list by using the download link.
After providing your email, you are sent a one-time-passcode to download the file.

### Libraries

 - Gorilla web Framework
 - MongoDB
 - Go Resty
 - Go.env
 - llimitman/Tint

### Setup
- Create a "files" folder for the uploaded files to be stored in.
- Create a server/ssl folder and put your SSL certificate in there
- Create a .env file with the following variables:
    - AUTH0_CLIENT_SECRET: Your Auth0 Client Secret.
    - AUTH0_CLIENT_ID: Yur Auth0 Client ID.
    - AUTH0_DOMAIN: # The URL of your Auth0 Tenant Domain.
    - DOMAIN # The domain your application is hosted on.
- On the system the application is running on, creaate the MONGODB_URI variable with link to your MongoDB database.

Build main.go to run the project.
