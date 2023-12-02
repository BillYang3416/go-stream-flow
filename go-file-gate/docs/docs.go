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
        "v1.createUserProfileRequest": {
            "type": "object",
            "required": [
                "accessToken",
                "displayName",
                "pictureUrl",
                "userId"
            ],
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "1234567890abcdef1234567890abcdef"
                },
                "displayName": {
                    "type": "string",
                    "example": "John Doe"
                },
                "pictureUrl": {
                    "type": "string",
                    "example": "https://example.com/picture.jpg"
                },
                "userId": {
                    "type": "string",
                    "example": "U1234567890abcdef1234567890abcdef"
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
                },
                "userId": {
                    "type": "string",
                    "example": "U1234567890abcdef1234567890abcdef"
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
	Description:      "This is a Go File Gate API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
