# Image converter Microservice in GO

This microservices makes it possible to offer an image in different sizes and formates for every device accessing your site. It makes it simple by managing the different files required and offering a simple API to access them.

### Usage

To upload an image, make a POST request to `/upload` with your image in a form. The server handles the converting to 3 different resolutions: `1920px`, `1280px` and `720px` width (and the original size).

To request an image, make a GET request to `/image/{uuid}/{resolution}`.