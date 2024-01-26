// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/line-callback": {
            "get": {
                "description": "Handler the redirect from Line Login",
                "tags": [
                    "Auth"
                ],
                "summary": "Line Callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/line-login": {
            "get": {
                "description": "Redirect to Line Login",
                "tags": [
                    "Auth"
                ],
                "summary": "Line Login",
                "responses": {
                    "302": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "get": {
                "description": "Logout",
                "tags": [
                    "Auth"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "register information",
                        "name": "RegisterRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Succesfully registered",
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterResponse"
                        }
                    }
                }
            }
        },
        "/user-profiles": {
            "post": {
                "description": "Create user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Profile"
                ],
                "summary": "Create user profile",
                "parameters": [
                    {
                        "description": "user profile information",
                        "name": "createUserProfileRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.createUserProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.userProfileResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        },
        "/user-profiles/{userId}": {
            "get": {
                "description": "Get user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Profile"
                ],
                "summary": "Get user profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.userProfileResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        },
        "/user-uploaded-files": {
            "get": {
                "description": "Get paginated files",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Uploaded File"
                ],
                "summary": "Get paginated files",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "last id",
                        "name": "lastID",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.getPaginatedFilesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create user uploaded file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Uploaded File"
                ],
                "summary": "Create user uploaded file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email recipient",
                        "name": "emailRecipient",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.UserUploadedFile": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "emailRecipient": {
                    "description": "The email address of the recipient",
                    "type": "string"
                },
                "emailSent": {
                    "description": "Indicates if the email was sent successfully",
                    "type": "boolean"
                },
                "emailSentAt": {
                    "description": "The timestamp when the email was sent",
                    "type": "string"
                },
                "errorMessage": {
                    "description": "// Error message if the email was not sent successfully",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "v1.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "v1.RegisterRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "displayName": {
                    "type": "string",
                    "example": "billyang"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "username": {
                    "type": "string",
                    "example": "useraname"
                }
            }
        },
        "v1.RegisterResponse": {
            "type": "object",
            "properties": {
                "userID": {
                    "type": "string",
                    "example": "userID"
                }
            }
        },
        "v1.createUserProfileRequest": {
            "type": "object",
            "required": [
                "displayName",
                "pictureUrl"
            ],
            "properties": {
                "displayName": {
                    "type": "string",
                    "example": "John Doe"
                },
                "pictureUrl": {
                    "type": "string",
                    "example": "https://example.com/picture.jpg"
                }
            }
        },
        "v1.errorResponse": {
            "type": "object",
            "properties": {
                "fieldErrs": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "v1.getPaginatedFilesResponse": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.UserUploadedFile"
                    }
                },
                "totalRecords": {
                    "type": "integer"
                }
            }
        },
        "v1.userProfileResponse": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string",
                    "example": "John Doe"
                },
                "pictureUrl": {
                    "type": "string",
                    "example": "https://example.com/picture.jpg"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "",
	Description:      "This is a Go Flow Gateway API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
