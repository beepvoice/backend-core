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

Unless otherwise noted, bodies and responses are with ```Content-Type: application/json```.

| Contents |
| -------- |
| [Create User](#Create-User) |
| [Get Users by Phone](#Get-Users-by-Phone) |
| [Get User](#Get-User) |
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
| first_name | String | First name of the added user. | ✓ |
| last_name | String | Last name of the added user. | ✓ |
| phone_number | String | Phone number of the added user. Shouldn't be needed but makes life easier. | X |

#### Success Response (200 OK)

Created user object.

```json
{
  "id": "<id>",
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

### Get User

```
GET /user/:user
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

### Create Conversation

```
POST /user/:user/conversation
```

Create a new conversation for a user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| title | String | Title of the conversation | X |

#### Success Response (200 OK)

Conversation object.

```json
{
  "id": "<id>",
  "title": "<title>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body. |
| 404 | User with supplied ID could not be found in database. |
| 500 | Error occurred inserting entries into the database. |

---

### Delete Conversation

```
DELETE /user/:user/conversation/:conversation
```

Delete the specified conversation.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |
| conversation | String | Conversation's ID. | ✓ |

#### Success Response (200 OK)

Empty body.

#### Errors

| Code | Description |
| ---- | ----------- |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred deleting entries from the database. |

---

### Update Conversation

```
PATCH /user/:user/conversation/:conversation
```

Update a conversation's details (mainly just title for now).

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |
| conversation | String | Conversation's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| title | String | New title of the conversation. | X |

#### Success Response (200 OK)

Empty Body. (TODO: Updated conversation)

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body. |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversations

```
GET /user/:user/conversation
```

Get the conversations of the specified user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |

#### Success Response (200 OK)

List of conversations.

```json
[
  {
    "id": "<id>",
    "title": "<title>"
  },
  ...
]
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversation

```
GET /user/:user/conversation/:conversation
```

Get a specific conversation of a specific user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |
| conversation | String | Conversation's ID. | ✓ |

#### Success Response (200 OK)

Conversation object.

```json
{
  "id": "<id>",
  "title": "<title>"
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred retrieving entries from the database. |

---

### Create Conversation Member

```
POST /user/:user/conversation/:conversation/member
```

Add a member to the specified conversation of the specified member.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |
| conversation | String | Conversation's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | String | ID of the user to be added. | ✓ |

#### Success Response (200 OK)

Empty body.

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/The length of the ID supplied in the body is less than 1. |
| 404 | User/Conversation with supplied ID could not be found in database. |
| 500 | Error occurred updating entries in the database. |

---

### Get Conversation Members

```
GET /user/:user/conversation/:conversation/member
```

Get the members of the specified conversation of the specified member.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |
| conversation | String | Conversation's ID. | ✓ |

#### Success (200 OK)

List of user objects in conversation.

```json
[
  {
    "id": "<id>",
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
| 500 | Error occurred retrieving entries from the database. |

---

### Create Contact

```
POST /user/:user/contact
```

Add a new contact for the specified user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |

#### Body

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| phone_number | String | New contact's phone number. A blank user object will be created if no such ID exists in the database. | ✓ |

#### Success Response (200 OK)

Empty body

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error occurred parsing the supplied body/The length of the ID supplied in the body is less than 1 or equal to the user's ID. |
| 500 | Error occurred updating entries in the database. |

---

### Get Contacts

```
GET /user/:user/contact
```

Get the contacts of the specified user.

#### URL Params

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| user | String | User's ID. | ✓ |

#### Success (200 OK)

List of user objects in user's contacts.

```json
[
  {
    "id": "<id>",
    "first_name": "<first_name>",
    "last_name": "<last_name>"
  },
  ...
]
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 500 | Error occurred retrieving entries from the database. |
