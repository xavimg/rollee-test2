# instructions

Write an http service that will accept 2 endpoints:

1. The first endpoint will accept a word (eg. "animal", "house") and store it in the service. Definition for a word is the following regex `[a-zA-Z]+`
2. The second endpoint will accept the beginning of a word (eg. "an") and returns the most frequent word stored in the service
3. The service will be case insensitive

scenario example
```
POST /service/word="abc"
POST /service/word="ab"
POST /service/word="ab"

GET /service/prefix="a"     => response: "ab"
GET /service/prefix="ab"   => response: "ab"
GET /service/prefix="abc"   => response: "abc"
GET /service/prefix="d"     => response: null
```

**Its up to you to choose a design in the API that makes sense**

The service must be written in Go, expect that there is no restrictions in terms of technologies or designs

Bonus points:
- Use docker
- Scalability
- Performances

