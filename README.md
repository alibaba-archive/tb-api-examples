# Teambition API in Example

## Current supported languages

* [nodejs](nodejs)
* [golang](golang)
* [python](python)
* [java](java)
* [postman(curl and others)](postman)
* ...

## Useful attributes

```
CLIENT_ID = 'xxxxxx'
CLIENT_SECRET = 'xxxxxx'
REDIRECT_URI = 'http://localhost:3000/tb/callback'
API_HOST = 'https://www.teambition.com/api'
ACCOUNT_HOST = 'https://account.teambition.com'
```

## Contribute

1. Create a directory named with the language (e.g. nodejs).
2. Copy README.md from nodejs/README.md into your directory, modify the content as you wish.
3. Start coding, the example should contain an embedded http server, an OAuth2 redirect process and a group of api calls.
4. Testing, visit your site at http://localhost:3000/auth, it should trigger a 301 redirect and finally return the project, task, tasklists, event and members in format of json.
5. Push your code and make a pull request.
