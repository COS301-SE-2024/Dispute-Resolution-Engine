# User Endpoints Contracts

---

**Note:** Responses are standardized to only contain error field if there is an error with a request, else only a data field will be present in the response.

---

## Dispute Summary
### Request


| Parameter     | Type   | Required | Description                    |
|---------------|--------|----------|--------------------------------|
| request_type  | string | Yes      | The type of request being made |
| body          | object | Yes      | The body of the request        |
| body.id       | string | Yes      | The id of the user             |

```json
{
    "request_type": "dispute_summary",
    "body": {
        "id": "some_id"
    }
}
```

### response

| Parameter     | Type    | Required | Description                      |
|---------------|---------|----------|----------------------------------|
| status        | integer | Yes      | HTTP status code                 |
| data          | array   | No       | Array of dispute objects         |
| data[].id     | string  | No       | Unique identifier for each dispute |
| data[].title  | string  | No       | Title of the dispute             |
| error         | string  | No       | Title of the dispute             |

#### success

```json
{
    "status": 200,
    "data": [
        {
            "id": "1",
            "title": "Dispute 1"
        },
        {
            "id": "2",
            "title": "Dispute 2"
        },
        ...
    ]
}
```

#### failure

```json
{
    "status": 400,
    "error": "missing required fields"
}
```
---

## Register

### Request
| Parameter         | Type   | Required | Description                    |
|-------------------|--------|----------|--------------------------------|
| request_type      | string | Yes      | The type of request being made |
| body              | object | Yes      | The body of the request        |
| body.first_name   | string | Yes      | The first name of the user     |
| body.surname      | string | Yes      | The surname of the user        |
| body.password_hash| string | Yes      | The hashed password of the user|
| body.email        | string | Yes      | The email address of the user  |

```json
{
    "request_type": "create_account",
    "body": {
        "first_name": "example",
        "surname": "example",
        "password_hash": "somePasswordHash",
        "email": "example@test.com"
    }
}
```

### Responses

| Parameter      | Type    | Required | Description              |
|----------------|---------|----------|--------------------------|
| status         | integer | Yes      | HTTP status code         |
| data           | object  | No       | The body of the response |
| data.message   | string  | No       | Response message         |
| error          | string  | No       | Title of the dispute     |


#### success
```json
{
    "status": 200,
    "data": {
        "message": "account created"
    }
}
```
#### failure
```json
{
    "status": 400,
    "error": "missing required fields"
}
```
---

## Login
### Request

| Parameter         | Type   | Required | Description                    |
|-------------------|--------|----------|--------------------------------|
| request_type      | string | Yes      | The type of request being made |
| body              | object | Yes      | The body of the request        |
| body.email        | string | Yes      | The email address of the user  |
| body.password_hash| string | Yes      | The hashed password of the user|

```json
{
    "request_type": "login",
    "body": {
        "email": "example@test.com",
        "password_hash": "somePasswordHash"
    }
}
```
### Response

| Parameter      | Type    | Required | Description              |
|----------------|---------|----------|--------------------------|
| status         | integer | Yes      | HTTP status code         |
| data           | object  | Yes      | The body of the response |
| data.message   | string  | Yes      | Response message         |

#### success

```json
{
    "status": 200,
    "data": {
        "message": "login successful"
    }
}
```

#### failure

```json
{
    "status": 500,
    "error": "user does not exist"
}
```
