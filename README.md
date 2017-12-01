# mongodbcaptchastore

A storage engine for
[github.com/dchest/captcha](https://github.com/dchest/captcha) using MongoDB

### Dependencies

* [gopkg.in/mgo.v2](https://gopkg.in/mgo.v2)
* [github.com/stretchr/testify](https://github.com/stretchr/testify)

### How To Use

```
timeout, _ := time.ParseDuration("5s")
expiration, _ := time.ParseDuration("1m")
s, err := mongodbcaptchastore.New("mongodb://localhost:27017", "captcha-db", "capcha-collection", 100 * 1024 * 1024, 100000, timeout, expiration)
if err != nil {
    panic(err)
}

// then use as custom store
captcha.SetCustomStore(s)
```
