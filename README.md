# GO Flow

Call Akvo Flow Web-form from CLI.

![cover](https://user-images.githubusercontent.com/3245109/93149946-0733b780-f722-11ea-8a5d-f983ae16e84d.jpg)

### Installation

```bash
go get github.com/dedenbangkit/goflow
```

Setup the required environment below:

```bash
AUTH0_CLIENT_FLOW='XXXXXXXXXXXXXXXX'
AUTH0_USER='<your_email@emailprovider>'
AUTH0_PWD='<your_password>'
AUTH0_URL='https://akvo.eu.auth0.com/oauth'
FLOW_API_URL='https://api-auth0.akvo.org/flow/orgs/<your_organisation>'
```

### Usage

```bash
goflow
```
![goflow](https://user-images.githubusercontent.com/3245109/93194120-ed21c580-f771-11ea-806c-208a9cb4670f.gif)

### Related Links

- [Flow Web-form](https://flowsupport.akvo.org/article/show/111363-webforms)
- [Akvo Flow REST API](https://github.com/akvo/akvo-flow-api/wiki/Akvo-Flow-REST-API)
- [Akvo SSO login](https://github.com/akvo/akvo-flow-api/wiki/Akvo-SSO-login)
