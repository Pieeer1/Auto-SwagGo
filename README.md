# Auto Swaggo

For those that hate maintaining swagger docs manually, this is a http handler that auto generates open api and swagger documentation for a net/http http server.

## Mux Features

- Invalid HTTP Methods Automatically Respond with a 405 (Method not Allowed)
- Invalid Request Bodies Response With a 422 (Unprocessable Entity)
- Auth Callback Failure Responds with a 401 (Unauthorized)
- Authorization Callback Failure Responds with a 403 (Forbidden)
- Version Handling and Multiple Swagger Docs For Versions
- Automatic OpenAPI and Swagger Endpoint Creation
- Ability to automatically open a browser on app run (for local and debugging purposes)

## Usage

