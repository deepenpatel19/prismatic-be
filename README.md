# prismatic-be
- It is open source repository which show case basic functionalities of social media platform.

### Installation
- Is is developed using Go with 1.23.0 version.
- You can install dependencies by executing command `go install`.

#### Features
- CRUD operations for user, friend management, post & comments.

#### Env example
- Please change values accourding to resource & environment.
```
{
    "environment": "local",
    "database": {
        "db_name": "",
        "db_password": "",
        "db_user": "",
        "db_host": "",
        "db_port": "",
        "db_ssl_mode": ""
    },
    "password_hash_cost": 0,
    "auth_realm": "",
    "auth_secret_key": ""
}

```


#### DB Schema
![Alt text](./docs/Screenshot%202024-08-29%20at%2010.04.43 PM.png)

#### API endpoint collection
- [Postman Collection](./docs/prismatic.postman_collection.json)

#### References
- Gin => https://github.com/gin-gonic/gin
- Timeout middleware => https://dev.to/jacobsngoodwin/13-gin-handler-timeout-middleware-4bhg
- DB migration management => https://github.com/golang-migrate/migrate