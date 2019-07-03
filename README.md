# backend-core

Beep backend handling core relational information that is not updated often.

## Quickstart

```
cockroach start --insecure
echo "create database core;" | cockroach sql --insecure

migrate -database cockroach://root@localhost:26257/core?sslmode=disable -source file://migrations goto 1
go build && ./core
```

## Environment variables

Supply environment variables by either exporting them or editing ```.env```.

| ENV | Description | Default |
| ---- | ----------- | ------- |
| LISTEN | Host and port number to listen on | :8080 |
| POSTGRES | URL of Postgres | postgresql://root@localhost:26257/core?sslmode=disable |

## API

Unless otherwise noted, bodies and responses are with `Content-Type: application/json`. Endpoints marked with a ```*``` require a populated `X-User-Claim` header from `backend-auth`.

| Contents |
| -------- |
| [Create User](#Create-User) |
| [Get Users by Phone](#Get-Users-by-Phone) |
| [Get User by ID](#Get-User-by-ID) |
| [Get User by Username](#Get-User-by-Username) ]
| [Update User](#Update-User) |
| [Create Conversation](#Create-Conversation) |
| [Delete Conversation](#Delete-Conversation) |
| [Update Conversation](#Update-Conversation) |
| [Get Conversations](#Get-Conversations) |
| [Get Conversation](#Get-Conversation) |
| [Create Conversation Member](#Create-Conversation-Member) |
| [Get Conversation Members](#Get-Conversation-Members) |
| [Create Contact](#Create-Contact) |
| [Get Contacts](#Get-Contacts) |

---

### Create User

```
POST /user
```

Create a new user.

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| username | String | Username of the added user. Must be unique. | ✓ |
| bio | String | Bio of the added user | ✓ |
| profile_pic | String | URL of added user's profile picture | ✓ |
| first_name | String | First name of the added user. | ✓ |
| last_name | String | Last name of the added user. | ✓ |
| phone_number | String | Phone number of the added user. Shouldn't be needed but makes life easier. | X |

#### Success Response (200 OK)

Created user object.

```json
{
  "id": "<id>",
  "username": "<username>",
  "bio": "<bio>",
  "profile_pic: "<profile_pic>",
  "first_name": "<first_name>",
  "last_name": "<last_name>",
  "phone_number": "<phone_number>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error parsing submitted body, or fields first_name or last_name have a length of 0. |
| 500 | Error occurred inserting entry into database. |

---

### Get Users by Phone

```
GET /user
```

Get user(s) associated with the supplied phone number.

#### Querystring

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| phone_number | String | Phone number to be queried. | ✓ |

#### Success Response (200 OK)

List of users.

```json
[
  {
    "id": "<id>",
    "username": "<username>",
    "bio": "<bio>",
    "profile_pic": "<profile_pic>",
    "first_name": "<first_name>",
    "last_name": "<last_name>"
  },
  ...
]
```

#### Errors


| Code | Description |
| ---- | ----------- |
| 400 | Supplied phone_number is absent/an invalid phone number. |
| 500 | Error occurred retrieving entries from database. |

---

### Get User by ID

```
GET /user/id/:user
```

Get a specific user by ID.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |

#### Success Response (200 OK)

User object.

```json
{
  "id": "<id>",
  "username": "<username>",
  "bio": "<bio>",
  "profile_pic": "<profile_pic>",
  "first_name": "<first_name>",
  "last_name": "<last_name>",
  "phone_number": "<phone_number>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 404 | User with supplied ID could not be found in database |
| 500 | Error occurred retrieving entries from database. |

---

### Get User by Username

```
GET /user/username/:username
```

Get a specific user by username.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| username | String | User's username. | ✓ |

#### Success Response (200 OK)

User object.

```json
{
  "id": "<id>",
  "username": "<username>",
  "bio": "<bio>",
  "profile_pic": "<profile_pic>",
  "first_name": "<first_name>",
  "last_name": "<last_name>",
  "phone_number": "<phone_number>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 404 | User with supplied username could not be found in database |
| 500 | Error occurred retrieving entries from database. |

---

### Update User

```
PATCH /user
```

Update an existing user. User ID is taken from header supplied by `backend-auth`. If one does not wish to update a field, leave it the value acquire from Get-User.

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| username | String | Updated username. | X |
| bio | String | Updated bio. | X |
| profile_pic | String | Updated URL of profile picture. | X |
| first_name | String | Updated first name. | X |
| last_name | String | Updated last name | X |

#### Success (200 OK)

Empty body.

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error parsing body/User with username already exists. |
| 500 | Error occurred updating database. |

---

### Create Conversation*

```
POST /user/conversation
```

Create a new conversation for a user.

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| title | String | Title of the conversation | X |
| dm | Boolean | Whether the conversation is a DM or not | X |
| picture | String | URL of the group's picture | X |

#### Success Response (200 OK)

Conversation object.

```json
{
  "id": "<id>",
  "title": "<title>",
  "picture": "<picture>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/Invalid `X-User-Claim` header |
| 404 | User with supplied ID could not be found in database. |
| 500 | Error occurred inserting entries into the database. |

---

### Delete Conversation*

```
DELETE /user/conversation/:conversation
```

Delete the specified conversation.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| conversation | String | Conversation's ID. | ✓ |

#### Success Response (200 OK)

Empty body.

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Invalid `X-User-Claim` header. |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred deleting entries from the database. |

---

### Update Conversation*

```
PATCH /user/conversation/:conversation
```

Update a conversation's details (mainly just title for now).

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| conversation | String | Conversation's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| title | String | New title of the conversation. | X |
| picture | String | New URL of the group's picture | X |

#### Success Response (200 OK)

Empty Body. (TODO: Updated conversation)

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/Invalid `X-User-Claim` header. |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversations*

```
GET /user/conversation
```

Get the conversations of the specified user.

#### Success Response (200 OK)

List of conversations.

```json
[
  {
    "id": "<id>",
    "title": "<title>"
    "picture": "<picture>"
  },
  ...
]
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Invalid `X-User-Claim` header. |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversation*

```
GET /user/conversation/:conversation
```

Get a specific conversation of a specific user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| conversation | String | Conversation's ID. | ✓ |

#### Success Response (200 OK)

Conversation object.

```json
{
  "id": "<id>",
  "title": "<title>",
  "picture": "<picture>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Invalid `X-User-Claim` header. |
| 404 | Conversation with supplied ID could not be found in database. |
| 500 | Error occurred retrieving entries from the database. |

---

### Create Conversation Member*

```
POST /user/conversation/:conversation/member
```

Add a member to the specified conversation.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| conversation | String | Conversation's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | String | ID of the user to be added. | ✓ |

#### Success Response (200 OK)

The conversation ID of the conversation the user is added to.

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/The length of the ID supplied in the body is less than 1/Invalid `X-User-Claim` header. |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversation Members*

```
GET /user/conversation/:conversation/member
```

Get the members of the specified conversation.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| conversation | String | Conversation's ID. | ✓ |

#### Success (200 OK)

List of user objects in conversation.

```json
[
  {
    "id": "<id>",
    "username": "<username>",
    "bio": "<bio>",
    "profile_pic": "<profile_pic>",
    "first_name": "<first_name>",
    "last_name": "<last_name>",
    "phone_number": "<phone_number>"
  },
  ...
]
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Invalid `X-User-Claim` header. |
| 500 | Error occurred retrieving entries from the database. |

---

### Create Contact*

```
POST /user/contact
```

Add a new contact.

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| phone_number | String | New contact's phone number. A blank user object will be created if no such ID exists in the database. | ✓ |

#### Success Response (200 OK)

Empty body

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/The length of the ID supplied in the body is less than 1 or equal to the user's ID/Invalid `X-User-Claim` header. |
| 500 | Error occurred updating entries in the database. |

---

### Get Contacts

```
GET /user/contact
```

Get the user's contacts.

#### Success (200 OK)

List of user objects in user's contacts.

```json
[
  {
    "id": "<id>",
    "username": "<username>",
    "bio": "<bio>",
    "profile_pic": "<profile_pic",
    "first_name": "<first_name>",
    "last_name": "<last_name>"
  },
  ...
]
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Invalid `X-User-Claim` header. |
| 500 | Error occurred retrieving entries from the database. |
