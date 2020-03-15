# ShortLink ![](https://img.shields.io/badge/language-golang-blue)  
Short link generation service by golang
# Install
`go build .`  
# Usage
`./ShortLink`  
# API
  - POST /api/shorten  
  - GET /api/info?shortlink=shortlink  
  - GET /:shortlink return 302 code  
# Note  
  Use redis as the storage backend by default  
  If you want to replace other storage backends (e.g. mysql), you can implement the method of the Storage interface, similar to redis  
  There are some default environment variables  
  - APP_REDIS_ADDR (default "localhost:6379")
  - APP_REDIS_PASSWD (default "")
  - APP_REDIS_DB (default "0")
# Thank
Learn from [Jacky_1024](https://www.imooc.com/learn/1150)
# License
![](https://img.shields.io/badge/License-MIT-blue.svg)
