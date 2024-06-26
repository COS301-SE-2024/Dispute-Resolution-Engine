// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Login an existing user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User Credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.Credentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "description": "Reset a user's password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset a user's password",
                "responses": {
                    "200": {
                        "description": "Password reset not available yet...",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User Details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "post": {
                "description": "Verifies the user's email by checking the provided pin code against stored values.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify user email",
                "parameters": [
                    {
                        "description": "Verify User",
                        "name": "pinReq",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.VerifyUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Email verified successfully - Example response: { 'message': 'Email verified successfully' }",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Invalid pin - Example error response: { 'error': 'Invalid pin' }",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Error verifying pin - Example error response: { 'error': 'Error verifying pin' }",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/dispute": {
            "get": {
                "description": "Get a summary list of disputes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dispute"
                ],
                "summary": "Get a summary list of disputes",
                "responses": {
                    "200": {
                        "description": "Dispute Summary Endpoint",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/dispute/{id}": {
            "get": {
                "description": "Get a dispute",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dispute"
                ],
                "summary": "Get a dispute",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dispute ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Dispute Detail Endpoint",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update a dispute",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dispute"
                ],
                "summary": "Update a dispute",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dispute ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Dispute Patch Endpoint",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
                "description": "Get user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "User profile not available yet...",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Update user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user profile",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User updated successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/user/remove": {
            "delete": {
                "description": "Remove user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Remove user account",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.DeleteUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User account removed successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.Credentials": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.CreateUser": {
            "type": "object",
            "properties": {
                "address_type": {
                    "type": "string"
                },
                "birthdate": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "country": {
                    "description": "These are the user's address details",
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "description": "These are all the user details that are required to create a user",
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "nationality": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                },
                "preferred_language": {
                    "type": "string"
                },
                "province": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                },
                "street2": {
                    "type": "string"
                },
                "street3": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                },
                "timezone": {
                    "type": "string"
                }
            }
        },
        "models.DeleteUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "type": "string"
                }
            }
        },
        "models.UpdateUser": {
            "type": "object",
            "properties": {
                "address_type": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "code": {
                    "description": "This is the country code",
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "province": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                },
                "street2": {
                    "type": "string"
                },
                "street3": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "address_id": {
                    "description": "what the fuck",
                    "type": "integer"
                },
                "birthdate": {
                    "description": "check",
                    "type": "string"
                },
                "createdAt": {
                    "description": "Filled in by API",
                    "type": "string"
                },
                "email": {
                    "description": "check",
                    "type": "string"
                },
                "first_name": {
                    "description": "check",
                    "type": "string"
                },
                "gender": {
                    "description": "check",
                    "type": "string"
                },
                "id": {
                    "description": "Filled in by API",
                    "type": "integer"
                },
                "lastLogin": {
                    "description": "Filled in by API",
                    "type": "string"
                },
                "nationality": {
                    "description": "check",
                    "type": "string"
                },
                "password": {
                    "description": "Updated by API",
                    "type": "string"
                },
                "phone_number": {
                    "description": "need",
                    "type": "string"
                },
                "preferred_language": {
                    "description": "worked on",
                    "type": "string"
                },
                "role": {
                    "description": "Filled in by API",
                    "type": "string"
                },
                "salt": {
                    "description": "Filled in by API",
                    "type": "string"
                },
                "status": {
                    "description": "Filled in by API",
                    "type": "string"
                },
                "surname": {
                    "description": "check",
                    "type": "string"
                },
                "timezone": {
                    "description": "need to be handled by me?",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "Filled in by API",
                    "type": "string"
                }
            }
        },
        "models.VerifyUser": {
            "type": "object",
            "properties": {
                "pin": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Dispute Resolution Engine - v1",
	Description:      "This is a description.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
