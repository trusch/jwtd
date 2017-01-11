# jwtd - a JSON Web Token Daemon

jwtd is a service responsible for authenticating users

It provides cryptographically signed JSON Web Tokens based on a role based access control system. These tokens can then be validated using jwtd-proxy, providing the possibility to secure arbitary http based services.

## Concept

jwtd is a HTTP server. It provides an endpoint which implements token creation. You simply send your username + password and a specification of what you want to access and jwtd sends a token back to you if your user is allowed to do this.

To declare users, roles and access rights jwtd has a companion program called jwtd-ctl. This is a commandline tool which allows you to:

- set username + password combinations
- associate users with groups
- declare access rights for groups

You can also use a restfull API provided by jwtd.

Access rights are tuples of a service id and a list of labels. The service is typically the URI of the some service you (as a system architect) want to protect. The labels are string-string tuples and used to implement fine grained control. It's up to the usecase which labels exists.

## Features

- JWTD Core

  - User Management
  - Group Management
  - Multi Project Support
  - Commandline Interface
  - REST API Interface
  - Wildcard Rights

- JWTD Proxy

  - Token validation for multiple services
  - TLS Termination
  - Pattern based Requierements
  - Variables in Route->Requirement mapping
  - File based Configuration

## Example

The `example` directory contains a docker based example setup. It only assumes that the hostnames "jwtd" and "http-echo" are locally resolvable to 127.0.0.1\. You can achieve that by manipulation your /etc/hosts for example. You can start the example via `bash example.sh`.

A more complex setup is shown in the `test` directory. It contains the integration test setup for this project. Therefore it starts by running an example setup with docker and then fires HTTP requests against jwtd-proxy for testing purpuses. I think about 90% of the functionality of jwtd is covered here. The tests are written as golang tests and can therefore be run via `go test -v`.

## API

### Group Handling

- list groups

  - `GET /project/{project}/group`

- create group

  - `POST /project/{project}/group`
  - body should look like:

    ```json
    {
      "name": "boss-group",
      "rights": {
        "admin.hosting.com": {
          "role": "admin"
        }
      }
    }
    ```

- get group

  - `GET /project/{project}/group/{groupname}`

- delete group

  - `DELETE /project/{project}/group/{groupname}`

- update group

  - `PATCH /project/{project}/group/{groupname}`
  - the rights object of the group is completly replaced, so be carefull
  - body should look like:

    ```json
    {
      "rights": {
        "admin.hosting.com": {
          "role": "admin",
          "meta": "abc"
        }
      }
    }
    ```

### User Handling

- list users

  - `GET /project/{project}/user`

- create user

  - `POST /project/{project}/user`
  - body should look like:

    ```json
    {
      "username": "hans",
      "password": "wurst",
      "groups": [
        "boss-group"
      ]
    }
    ```

- get user

  - `GET /project/{project}/user/{username}`

- delete user

  - `DELETE /project/{project}/user/{username}`

- update user

  - `PATCH /project/{project}/user/{username}`
  - the rights object of the user is completly replaced, so be carefull
  - for updating the password body should look like:

    ```json
    {
      "password": "newpassword"
    }
    ```
  - for updating the groups body should look like:

    ```json
    {
      "groups": [
        "boss-group",
        "user-group"
      ]
    }
    ```

### Token Generation

- get token

    - `GET /token`
    - body should look like
      ```json
      {
        "project": "my-project",
        "username": "hans",
        "password": "wurst",
        "service": "admin.hosting.com",
        "lifetime": "1h",
        "labels": {
          "role": "admin"
        }
      }
      ```
    - `lifetime` is optional and defaults to 10min
      - valid units are s, m and h, also mixed like `1h30m10s` for a token lifetime of 1 hour 30 minutes and 10 seconds

## Setup

1. Install jwtd and helpers
```bash
> go get -u -v github.com/trusch/jwtd \
    github.com/trusch/jwtd/jwtd-ctl \
    github.com/trusch/jwtd/jwtd-proxy \
    github.com/trusch/pki/pkitool
```

2. Create PKI for jwtd
```bash
> mkdir -p /etc/jwtd
> pkitool -p /etc/jwtd/pki init
> pkitool -p /etc/jwtd/pki issue server jwtd
```
Your config dir should now look like this:
```bash
> tree /etc/jwtd
/etc/jwtd/
-- pki
    |-- ca.crt
    |-- ca.key
    |-- jwtd.crt
    |-- jwtd.key
     -- serial
```

3. Create initial jwtd config
```bash
> mkdir -p /etc/jwtd/projects
> jwtd-ctl -c /etc/jwtd/projects -p project1 init
```
This will generate a config file in /etc/jwtd/projects/project1.yaml like this:
  ```yaml
  users:
  - name: admin
    passwordhash: $2a$10$VEdYiAT/JfN18pQYJA0OTeoTmTtzxeQhEfQcezQWmZHJsUDA7rgyC
    groups:
    - admin
  groups:
  - name: admin
    rights:
      jwtd:
        scope: admin
  ```

4. Create jwtd-proxy config in /etc/jwtd/proxy.yaml
```yaml
listen: :443
cert: /etc/jwtd/pki/jwtd.crt
hosts:
  jwtd:
    backend: http://localhost:8080
    project: project1
    tls:
      cert: /etc/jwtd/pki/jwtd.crt
      key: /etc/jwtd/pki/jwtd.key
    routes:
      - path: /token
        require: {}
      - path: /
        require:
          scope: admin
```

5. Start servers
```bash
> jwtd -config /etc/jwtd/projects -key /etc/jwtd/pki/jwtd.key -listen :8080 &
> jwtd-proxy -config /etc/jwtd/proxy.yaml &
```

6. Test your setup
```bash
>request='{
  "project":"project1",
  "username":"admin",
  "password":"admin",
  "service":"jwtd",
  "labels":{"scope":"admin"}
}'
> token=$(curl -k -H "Host: jwtd" --data "$request" https://localhost/token 2>/dev/null)
> curl -k -H "Host: jwtd" -H "Authorization: bearer $token" https://localhost/project/project1/group
```
This should output a list of all groups currently present in your project

## Special Features which will make you happy

### Wildcard Rights in JWTD
You can specify wildcard rights in a group, so that each user of that group can aquire all labels matching the wildcard.
As an example you could specify the following:
```yaml
users:
- name: admin
  passwordhash: $2a$10$VEdYiAT/JfN18pQYJA0OTeoTmTtzxeQhEfQcezQWmZHJsUDA7rgyC
  groups:
  - admin
groups:
- name: admin
  rights:
    service1:
      role: '*'
    service2:
      '*': admin
    service3:
      '*': '*'
```
The admin user can now :
* aquire labels for service1 with the key "role" and arbitary values
* aquire labels for service2 with arbitary keys and the value "admin"
* aquire arbitary labels for service3

### Variables in Route-Requirement Mapping in jwtd-proxy
You can use variables in your route definition and use them to specify which labels are needed. This gives you great flexibility with an acceptable complexity-overhead.
Check the following jwtd-proxy config:
  ```yaml
  listen: :443
  cert: /etc/jwtd/pki/jwtd.crt
  hosts:
    service1:
      backend: http://localhost:8080
      project: project1
      tls:
        cert: /etc/jwtd/pki/service1.crt
        key: /etc/jwtd/pki/service1.key
      routes:
        - path: /{scope}/{action}
          require:
            scope: $scope
            action: $action
  ```
If you request /foo/bar this rule will say that you need the labels (scope, foo) and (action, bar).
If you request /baz/qux this rule will say that you need the labels (scope, baz) and (action, qux).
