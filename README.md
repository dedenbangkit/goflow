# GO Flow

Call Akvo Flow Web-form from CLI.

![cover](https://user-images.githubusercontent.com/3245109/93150872-7f02e180-f724-11ea-963e-4fccf1be7f80.png)

### Installation

```bash
go get dedenbangkit/goflow
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

### Related Links

- [Flow Web-form](https://flowsupport.akvo.org/article/show/111363-webforms)
- [Akvo Flow REST API](https://github.com/akvo/akvo-flow-api/wiki/Akvo-Flow-REST-API)
