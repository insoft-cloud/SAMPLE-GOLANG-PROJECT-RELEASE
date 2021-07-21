package swag

import (
	"encoding/json"
	goparser "go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultParseDepth = 100

func TestNew(t *testing.T) {
	swagMode = test
	New()
}

func TestParser_ParseGeneralApiInfo(t *testing.T) {
	expected := `{
    "schemes": [
        "http",
        "https"
    ],
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.\nIt has a lot of beautiful features.",
        "title": "Swagger Example API",
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
        "version": "1.0",
        "x-logo": {
            "altText": "Petstore logo",
            "backgroundColor": "#FFFFFF",
            "url": "https://redocly.github.io/redoc/petstore-logo.png"
        }
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {},
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            },
            "x-tokenName": "id_token"
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-google-endpoints": [
        {
            "allowCors": true,
            "name": "name.endpoints.environment.cloud.goog"
        }
    ],
    "x-google-marks": "marks values"
}`
	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)
	p := New()
	err := p.ParseGeneralAPIInfo("testdata/main.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralApiInfoTemplated(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        }
    },
    "paths": {},
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-google-endpoints": [
        {
            "allowCors": true,
            "name": "name.endpoints.environment.cloud.goog"
        }
    ],
    "x-google-marks": "marks values"
}`
	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)
	p := New()
	err := p.ParseGeneralAPIInfo("testdata/templated.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralApiInfoExtensions(t *testing.T) {
	// should be return an error because extension value is not a valid json
	func() {
		expected := "annotation @x-google-endpoints need a valid json value"
		gopath := os.Getenv("GOPATH")
		assert.NotNil(t, gopath)
		p := New()
		err := p.ParseGeneralAPIInfo("testdata/extensionsFail1.go")
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	}()

	// should be return an error because extension don't have a value
	func() {
		expected := "annotation @x-google-endpoints need a value"
		gopath := os.Getenv("GOPATH")
		assert.NotNil(t, gopath)
		p := New()
		err := p.ParseGeneralAPIInfo("testdata/extensionsFail2.go")
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	}()
}

func TestParser_ParseGeneralApiInfoWithOpsInSameFile(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.\nIt has a lot of beautiful features.",
        "title": "Swagger Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "1.0"
    },
    "paths": {}
}`

	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)
	p := New()
	err := p.ParseGeneralAPIInfo("testdata/single_file_api/main.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralApiInfoFailed(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)
	p := New()
	assert.Error(t, p.ParseGeneralAPIInfo("testdata/noexist.go"))
}

func TestGetAllGoFileInfo(t *testing.T) {
	searchDir := "testdata/pet"

	p := New()
	err := p.getAllGoFileInfo("testdata", searchDir)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(p.packages.files))
}

func TestParser_ParseType(t *testing.T) {
	searchDir := "testdata/simple/"

	p := New()
	err := p.getAllGoFileInfo("testdata", searchDir)
	assert.NoError(t, err)

	_, err = p.packages.ParseTypes()

	assert.NoError(t, err)
	assert.NotNil(t, p.packages.uniqueDefinitions["api.Pet3"])
	assert.NotNil(t, p.packages.uniqueDefinitions["web.Pet"])
	assert.NotNil(t, p.packages.uniqueDefinitions["web.Pet2"])
}

func TestGetSchemes(t *testing.T) {
	schemes := getSchemes("@schemes http https")
	expectedSchemes := []string{"http", "https"}
	assert.Equal(t, expectedSchemes, schemes)
}

func TestParseSimpleApi1(t *testing.T) {
	expected, err := ioutil.ReadFile("testdata/simple/expected.json")
	assert.NoError(t, err)
	searchDir := "testdata/simple"
	mainAPIFile := "main.go"
	p := New()
	p.PropNamingStrategy = PascalCase
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "  ")
	assert.Equal(t, string(expected), string(b))
}

func TestParseSimpleApi_ForSnakecase(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "Swagger Example API",
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
        "version": "1.0"
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {
        "/file/upload": {
            "post": {
                "description": "Upload file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload file",
                "operationId": "file.upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "this is a test file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-string-by-int/{some_id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new pet to the store",
                "operationId": "get-string-by-int",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.Pet"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-struct-array-by-string/{some_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    },
                    {
                        "BasicAuth": []
                    },
                    {
                        "OAuth2Application": [
                            "write"
                        ]
                    },
                    {
                        "OAuth2Implicit": [
                            "read",
                            "admin"
                        ]
                    },
                    {
                        "OAuth2AccessCode": [
                            "read"
                        ]
                    },
                    {
                        "OAuth2Password": [
                            "admin"
                        ]
                    }
                ],
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "operationId": "get-struct-array-by-string",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            1,
                            2,
                            3
                        ],
                        "type": "integer",
                        "description": "Category",
                        "name": "category",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "default": 0,
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 50,
                        "type": "integer",
                        "default": 10,
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maxLength": 50,
                        "minLength": 1,
                        "type": "string",
                        "default": "\"\"",
                        "description": "q",
                        "name": "q",
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
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "web.APIError": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "error_code": {
                    "type": "integer"
                },
                "error_message": {
                    "type": "string"
                }
            }
        },
        "web.Pet": {
            "type": "object",
            "required": [
                "price"
            ],
            "properties": {
                "birthday": {
                    "type": "integer"
                },
                "category": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "integer",
                            "example": 1
                        },
                        "name": {
                            "type": "string",
                            "example": "category_name"
                        },
                        "photo_urls": {
                            "type": "array",
                            "format": "url",
                            "items": {
                                "type": "string"
                            },
                            "example": [
                                "http://test/image/1.jpg",
                                "http://test/image/2.jpg"
                            ]
                        },
                        "small_category": {
                            "type": "object",
                            "required": [
                                "name"
                            ],
                            "properties": {
                                "id": {
                                    "type": "integer",
                                    "example": 1
                                },
                                "name": {
                                    "type": "string",
                                    "example": "detail_category_name"
                                },
                                "photo_urls": {
                                    "type": "array",
                                    "items": {
                                        "type": "string"
                                    },
                                    "example": [
                                        "http://test/image/1.jpg",
                                        "http://test/image/2.jpg"
                                    ]
                                }
                            }
                        }
                    }
                },
                "coeffs": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "custom_string": {
                    "type": "string"
                },
                "custom_string_arr": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "data": {
                    "type": "object"
                },
                "decimal": {
                    "type": "number"
                },
                "id": {
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                },
                "is_alive": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "poti"
                },
                "null_int": {
                    "type": "integer"
                },
                "pets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet2"
                    }
                },
                "pets2": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet2"
                    }
                },
                "photo_urls": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "http://test/image/1.jpg",
                        "http://test/image/2.jpg"
                    ]
                },
                "price": {
                    "type": "number",
                    "example": 3.25
                },
                "status": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Tag"
                    }
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "web.Pet2": {
            "type": "object",
            "properties": {
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "middle_name": {
                    "type": "string"
                }
            }
        },
        "web.RevValue": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "integer"
                },
                "err": {
                    "type": "integer"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "web.Tag": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "format": "int64"
                },
                "name": {
                    "type": "string"
                },
                "pets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    }
}`
	searchDir := "testdata/simple2"
	mainAPIFile := "main.go"
	p := New()
	p.PropNamingStrategy = SnakeCase
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseSimpleApi_ForLowerCamelcase(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "Swagger Example API",
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
        "version": "1.0"
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {
        "/file/upload": {
            "post": {
                "description": "Upload file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload file",
                "operationId": "file.upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "this is a test file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-string-by-int/{some_id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new pet to the store",
                "operationId": "get-string-by-int",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.Pet"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        },
        "/testapi/get-struct-array-by-string/{some_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    },
                    {
                        "BasicAuth": []
                    },
                    {
                        "OAuth2Application": [
                            "write"
                        ]
                    },
                    {
                        "OAuth2Implicit": [
                            "read",
                            "admin"
                        ]
                    },
                    {
                        "OAuth2AccessCode": [
                            "read"
                        ]
                    },
                    {
                        "OAuth2Password": [
                            "admin"
                        ]
                    }
                ],
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "operationId": "get-struct-array-by-string",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            1,
                            2,
                            3
                        ],
                        "type": "integer",
                        "description": "Category",
                        "name": "category",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "default": 0,
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 50,
                        "type": "integer",
                        "default": 10,
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maxLength": 50,
                        "minLength": 1,
                        "type": "string",
                        "default": "\"\"",
                        "description": "q",
                        "name": "q",
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
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "web.APIError": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "errorCode": {
                    "type": "integer"
                },
                "errorMessage": {
                    "type": "string"
                }
            }
        },
        "web.Pet": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "integer",
                            "example": 1
                        },
                        "name": {
                            "type": "string",
                            "example": "category_name"
                        },
                        "photoURLs": {
                            "type": "array",
                            "format": "url",
                            "items": {
                                "type": "string"
                            },
                            "example": [
                                "http://test/image/1.jpg",
                                "http://test/image/2.jpg"
                            ]
                        },
                        "smallCategory": {
                            "type": "object",
                            "properties": {
                                "id": {
                                    "type": "integer",
                                    "example": 1
                                },
                                "name": {
                                    "type": "string",
                                    "example": "detail_category_name"
                                },
                                "photoURLs": {
                                    "type": "array",
                                    "items": {
                                        "type": "string"
                                    },
                                    "example": [
                                        "http://test/image/1.jpg",
                                        "http://test/image/2.jpg"
                                    ]
                                }
                            }
                        }
                    }
                },
                "data": {
                    "type": "object"
                },
                "decimal": {
                    "type": "number"
                },
                "id": {
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                },
                "isAlive": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "poti"
                },
                "pets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet2"
                    }
                },
                "pets2": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet2"
                    }
                },
                "photoURLs": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "http://test/image/1.jpg",
                        "http://test/image/2.jpg"
                    ]
                },
                "price": {
                    "type": "number",
                    "example": 3.25
                },
                "status": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Tag"
                    }
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "web.Pet2": {
            "type": "object",
            "properties": {
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "middleName": {
                    "type": "string"
                }
            }
        },
        "web.RevValue": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "integer"
                },
                "err": {
                    "type": "integer"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "web.Tag": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "format": "int64"
                },
                "name": {
                    "type": "string"
                },
                "pets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Pet"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "authorizationUrl": "",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    }
}`
	searchDir := "testdata/simple3"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseStructComment(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "Swagger Example API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:4000",
    "basePath": "/api",
    "paths": {
        "/posts/{post_id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new pet to the store",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "description": "Some ID",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "web.APIError": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "Error time",
                    "type": "string"
                },
                "error": {
                    "description": "Error an Api error",
                    "type": "string"
                },
                "errorCtx": {
                    "description": "Error ` + "`" + `context` + "`" + ` tick comment",
                    "type": "string"
                },
                "errorNo": {
                    "description": "Error ` + "`" + `number` + "`" + ` tick comment",
                    "type": "integer"
                }
            }
        }
    }
}`
	searchDir := "testdata/struct_comment"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseNonExportedJSONFields(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server.",
        "title": "Swagger Example API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:4000",
    "basePath": "/api",
    "paths": {
        "/so-something": {
            "get": {
                "description": "Does something, but internal (non-exported) fields inside a struct won't be marshaled into JSON",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Call DoSomething",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MyStruct"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.MyStruct": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Post data",
                    "type": "object",
                    "properties": {
                        "name": {
                            "description": "Post tag",
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                },
                "id": {
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                },
                "name": {
                    "description": "Post name",
                    "type": "string",
                    "example": "poti"
                }
            }
        }
    }
}`

	searchDir := "testdata/non_exported_json_fields"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParsePetApi(t *testing.T) {
	expected := `{
    "schemes": [
        "http",
        "https"
    ],
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.  You can find out more about     Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).      For this sample, you can use the api key 'special-key' to test the authorization     filters.",
        "title": "Swagger Petstore",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "email": "apiteam@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {}
}`
	searchDir := "testdata/pet"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseModelAsTypeAlias(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "Swagger Example API",
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
        "version": "1.0"
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {
        "/testapi/time-as-time-container": {
            "get": {
                "description": "test container with time and time alias",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get container with time and time alias",
                "operationId": "time-as-time-container",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/data.TimeContainer"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "data.TimeContainer": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        }
    }
}`
	searchDir := "testdata/alias_type"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseComposition(t *testing.T) {
	searchDir := "testdata/composition"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")

	//windows will fail: \r\n \n
	assert.Equal(t, string(expected), string(b))
}

func TestParseImportAliases(t *testing.T) {
	searchDir := "testdata/alias_import"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	//windows will fail: \r\n \n
	assert.Equal(t, string(expected), string(b))
}

func TestParseNested(t *testing.T) {
	searchDir := "testdata/nested"
	mainAPIFile := "main.go"
	p := New()
	p.ParseDependency = true
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, string(expected), string(b))
}

func TestParseDuplicated(t *testing.T) {
	searchDir := "testdata/duplicated"
	mainAPIFile := "main.go"
	p := New()
	p.ParseDependency = true
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.Errorf(t, err, "duplicated @id declarations successfully found")
}

func TestParseDuplicatedOtherMethods(t *testing.T) {
	searchDir := "testdata/duplicated2"
	mainAPIFile := "main.go"
	p := New()
	p.ParseDependency = true
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.Errorf(t, err, "duplicated @id declarations successfully found")
}

func TestParseConflictSchemaName(t *testing.T) {
	searchDir := "testdata/conflict_name"
	mainAPIFile := "main.go"
	p := New()
	p.ParseDependency = true
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	expected, err := ioutil.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParser_ParseStructArrayObject(t *testing.T) {
	src := `
package api

type Response struct {
	Code int
	Table [][]string
	Data []struct{
		Field1 uint
		Field2 string
	}
}

// @Success 200 {object} Response
// @Router /api/{id} [get]
func Test(){
}
`
	expected := `{
   "api.Response": {
      "type": "object",
      "properties": {
         "code": {
            "type": "integer"
         },
         "data": {
            "type": "array",
            "items": {
               "type": "object",
               "properties": {
                  "field1": {
                     "type": "integer"
                  },
                  "field2": {
                     "type": "string"
                  }
               }
            }
         },
         "table": {
            "type": "array",
            "items": {
               "type": "array",
               "items": {
                  "type": "string"
               }
            }
         }
      }
   }
}`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)
	_, err = p.packages.ParseTypes()
	assert.NoError(t, err)

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	out, err := json.MarshalIndent(p.swagger.Definitions, "", "   ")
	assert.NoError(t, err)
	assert.Equal(t, expected, string(out))

}

func TestParser_ParseEmbededStruct(t *testing.T) {
	src := `
package api

type Response struct {
	rest.ResponseWrapper
}

// @Success 200 {object} Response
// @Router /api/{id} [get]
func Test(){
}
`
	restsrc := `
package rest

type ResponseWrapper struct {
	Status   string
	Code     int
	Messages []string
	Result   interface{}
}
`
	expected := `{
   "api.Response": {
      "type": "object",
      "properties": {
         "code": {
            "type": "integer"
         },
         "messages": {
            "type": "array",
            "items": {
               "type": "string"
            }
         },
         "result": {
            "type": "object"
         },
         "status": {
            "type": "string"
         }
      }
   }
}`
	parser := New()
	parser.ParseDependency = true

	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)
	parser.packages.CollectAstFile("api", "api/api.go", f)

	f2, err := goparser.ParseFile(token.NewFileSet(), "", restsrc, goparser.ParseComments)
	assert.NoError(t, err)
	parser.packages.CollectAstFile("rest", "rest/rest.go", f2)

	_, err = parser.packages.ParseTypes()
	assert.NoError(t, err)

	err = parser.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	out, err := json.MarshalIndent(parser.swagger.Definitions, "", "   ")
	assert.NoError(t, err)
	assert.Equal(t, expected, string(out))

}

func TestParser_ParseStructPointerMembers(t *testing.T) {
	src := `
package api

type Child struct {
	Name string
}

type Parent struct {
	Test1 *string  //test1
	Test2 *Child   //test2
}

// @Success 200 {object} Parent
// @Router /api/{id} [get]
func Test(){
}
`

	expected := `{
   "api.Child": {
      "type": "object",
      "properties": {
         "name": {
            "type": "string"
         }
      }
   },
   "api.Parent": {
      "type": "object",
      "properties": {
         "test1": {
            "description": "test1",
            "type": "string"
         },
         "test2": {
            "description": "test2",
            "$ref": "#/definitions/api.Child"
         }
      }
   }
}`

	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)
	_, err = p.packages.ParseTypes()
	assert.NoError(t, err)

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	out, err := json.MarshalIndent(p.swagger.Definitions, "", "   ")
	assert.NoError(t, err)
	assert.Equal(t, expected, string(out))
}

func TestParser_ParseStructMapMember(t *testing.T) {
	src := `
package api

type MyMapType map[string]string

type Child struct {
	Name string
}

type Parent struct {
	Test1 map[string]interface{}  //test1
	Test2 map[string]string		  //test2
	Test3 map[string]*string	  //test3
	Test4 map[string]Child		  //test4
	Test5 map[string]*Child		  //test5
	Test6 MyMapType				  //test6
	Test7 []Child				  //test7
	Test8 []*Child				  //test8
	Test9 []map[string]string	  //test9
}

// @Success 200 {object} Parent
// @Router /api/{id} [get]
func Test(){
}
`
	expected := `{
   "api.Child": {
      "type": "object",
      "properties": {
         "name": {
            "type": "string"
         }
      }
   },
   "api.MyMapType": {
      "type": "object",
      "additionalProperties": {
         "type": "string"
      }
   },
   "api.Parent": {
      "type": "object",
      "properties": {
         "test1": {
            "description": "test1",
            "type": "object",
            "additionalProperties": true
         },
         "test2": {
            "description": "test2",
            "type": "object",
            "additionalProperties": {
               "type": "string"
            }
         },
         "test3": {
            "description": "test3",
            "type": "object",
            "additionalProperties": {
               "type": "string"
            }
         },
         "test4": {
            "description": "test4",
            "type": "object",
            "additionalProperties": {
               "$ref": "#/definitions/api.Child"
            }
         },
         "test5": {
            "description": "test5",
            "type": "object",
            "additionalProperties": {
               "$ref": "#/definitions/api.Child"
            }
         },
         "test6": {
            "description": "test6",
            "$ref": "#/definitions/api.MyMapType"
         },
         "test7": {
            "description": "test7",
            "type": "array",
            "items": {
               "$ref": "#/definitions/api.Child"
            }
         },
         "test8": {
            "description": "test8",
            "type": "array",
            "items": {
               "$ref": "#/definitions/api.Child"
            }
         },
         "test9": {
            "description": "test9",
            "type": "array",
            "items": {
               "type": "object",
               "additionalProperties": {
                  "type": "string"
               }
            }
         }
      }
   }
}`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)

	_, err = p.packages.ParseTypes()
	assert.NoError(t, err)

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	out, err := json.MarshalIndent(p.swagger.Definitions, "", "   ")
	assert.NoError(t, err)
	assert.Equal(t, expected, string(out))
}

func TestParser_ParseRouterApiInfoErr(t *testing.T) {
	src := `
package test

// @Accept unknown
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.EqualError(t, err, "ParseComment error in file  :unknown accept type can't be accepted")
}

func TestParser_ParseRouterApiGet(t *testing.T) {
	src := `
package test

// @Router /api/{id} [get]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Get)
}

func TestParser_ParseRouterApiPOST(t *testing.T) {
	src := `
package test

// @Router /api/{id} [post]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Post)
}

func TestParser_ParseRouterApiDELETE(t *testing.T) {
	src := `
package test

// @Router /api/{id} [delete]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)
	p := New()

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Delete)
}

func TestParser_ParseRouterApiPUT(t *testing.T) {
	src := `
package test

// @Router /api/{id} [put]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Put)
}

func TestParser_ParseRouterApiPATCH(t *testing.T) {
	src := `
package test

// @Router /api/{id} [patch]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Patch)
}

func TestParser_ParseRouterApiHead(t *testing.T) {
	src := `
package test

// @Router /api/{id} [head]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)
	p := New()

	err = p.ParseRouterAPIInfo("", f)

	assert.NoError(t, err)
	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Head)
}

func TestParser_ParseRouterApiOptions(t *testing.T) {
	src := `
package test

// @Router /api/{id} [options]
func Test(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Options)
}

func TestParser_ParseRouterApiMultiple(t *testing.T) {
	src := `
package test

// @Router /api/{id} [get]
func Test1(){
}

// @Router /api/{id} [patch]
func Test2(){
}

// @Router /api/{id} [delete]
func Test3(){
}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	ps := p.swagger.Paths.Paths

	val, ok := ps["/api/{id}"]

	assert.True(t, ok)
	assert.NotNil(t, val.Get)
	assert.NotNil(t, val.Patch)
	assert.NotNil(t, val.Delete)
}

// func TestParseDeterministic(t *testing.T) {
// 	mainAPIFile := "main.go"
// 	for _, searchDir := range []string{
// 		"testdata/simple",
// 		"testdata/model_not_under_root/cmd",
// 	} {
// 		t.Run(searchDir, func(t *testing.T) {
// 			var expected string

// 			// run the same code 100 times and check that the output is the same every time
// 			for i := 0; i < 100; i++ {
// 				p := New()
// 				p.PropNamingStrategy = PascalCase
// 				err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
// 				b, _ := json.MarshalIndent(p.swagger, "", "    ")
// 				assert.NotEqual(t, "", string(b))

// 				if expected == "" {
// 					expected = string(b)
// 				}

// 				assert.Equal(t, expected, string(b))
// 			}
// 		})
// 	}
// }

func TestApiParseTag(t *testing.T) {
	searchDir := "testdata/tags"
	mainAPIFile := "main.go"
	p := New(SetMarkdownFileDirectory(searchDir))
	p.PropNamingStrategy = PascalCase
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)

	if len(p.swagger.Tags) != 3 {
		t.Error("Number of tags did not match")
	}

	dogs := p.swagger.Tags[0]
	if dogs.TagProps.Name != "dogs" || dogs.TagProps.Description != "Dogs are cool" {
		t.Error("Failed to parse dogs name or description")
	}

	cats := p.swagger.Tags[1]
	if cats.TagProps.Name != "cats" || cats.TagProps.Description != "Cats are the devil" {
		t.Error("Failed to parse cats name or description")
	}

	if cats.TagProps.ExternalDocs.URL != "https://google.de" || cats.TagProps.ExternalDocs.Description != "google is super useful to find out that cats are evil!" {
		t.Error("URL: ", cats.TagProps.ExternalDocs.URL)
		t.Error("Description: ", cats.TagProps.ExternalDocs.Description)
		t.Error("Failed to parse cats external documentation")
	}
}

func TestApiParseTag_NonExistendTag(t *testing.T) {
	searchDir := "testdata/tags_nonexistend_tag"
	mainAPIFile := "main.go"
	p := New(SetMarkdownFileDirectory(searchDir))
	p.PropNamingStrategy = PascalCase
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.Error(t, err)
}

func TestParseTagMarkdownDescription(t *testing.T) {
	searchDir := "testdata/tags"
	mainAPIFile := "main.go"
	p := New(SetMarkdownFileDirectory(searchDir))
	p.PropNamingStrategy = PascalCase
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	if err != nil {
		t.Error("Failed to parse api description: " + err.Error())
	}

	if len(p.swagger.Tags) != 3 {
		t.Error("Number of tags did not match")
	}

	apes := p.swagger.Tags[2]
	if apes.TagProps.Description == "" {
		t.Error("Failed to parse tag description markdown file")
	}
}

func TestParseApiMarkdownDescription(t *testing.T) {
	searchDir := "testdata/tags"
	mainAPIFile := "main.go"
	p := New(SetMarkdownFileDirectory(searchDir))
	p.PropNamingStrategy = PascalCase
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	if err != nil {
		t.Error("Failed to parse api description: " + err.Error())
	}

	if p.swagger.Info.Description == "" {
		t.Error("Failed to parse api description: " + err.Error())
	}
}

func TestIgnoreInvalidPkg(t *testing.T) {
	searchDir := "testdata/deps_having_invalid_pkg"
	mainAPIFile := "main.go"
	p := New()
	if err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth); err != nil {
		t.Error("Failed to ignore valid pkg: " + err.Error())
	}
}

func TestFixes432(t *testing.T) {
	searchDir := "testdata/fixes-432"
	mainAPIFile := "cmd/main.go"

	p := New()
	if err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth); err != nil {
		t.Error("Failed to ignore valid pkg: " + err.Error())
	}
}

func TestParseOutsideDependencies(t *testing.T) {
	searchDir := "testdata/pare_outside_dependencies"
	mainAPIFile := "cmd/main.go"

	p := New()
	p.ParseDependency = true
	if err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth); err != nil {
		t.Error("Failed to parse api: " + err.Error())
	}
}

func TestParseStructParamCommentByQueryType(t *testing.T) {
	src := `
package main

type Student struct {
	Name string
	Age int
	Teachers []string
	SkipField map[string]string
}

// @Param request query Student true "query params"
// @Success 200
// @Router /test [get]
func Fun()  {

}
`
	expected := `{
    "info": {
        "contact": {}
    },
    "paths": {
        "/test": {
            "get": {
                "parameters": [
                    {
                        "type": "integer",
                        "name": "age",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "name": "teachers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    }
}`

	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)

	_, err = p.packages.ParseTypes()
	assert.NoError(t, err)

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseRenamedStructDefinition(t *testing.T) {
	src := `
package main

type Child struct {
	Name string
}//@name Student

type Parent struct {
	Name string
	Child Child
}//@name Teacher

// @Param request body Parent true "query params"
// @Success 200 {object} Parent
// @Router /test [get]
func Fun()  {

}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)
	_, err = p.packages.ParseTypes()
	assert.NoError(t, err)

	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	assert.NoError(t, err)
	teacher, ok := p.swagger.Definitions["Teacher"]
	assert.True(t, ok)
	ref := teacher.Properties["child"].SchemaProps.Ref
	assert.Equal(t, "#/definitions/Student", ref.String())
	_, ok = p.swagger.Definitions["Student"]
	assert.True(t, ok)
	path, ok := p.swagger.Paths.Paths["/test"]
	assert.True(t, ok)
	assert.Equal(t, "#/definitions/Teacher", path.Get.Parameters[0].Schema.Ref.String())
	ref = path.Get.Responses.ResponsesProps.StatusCodeResponses[200].ResponseProps.Schema.Ref
	assert.Equal(t, "#/definitions/Teacher", ref.String())
}

func TestParseJSONFieldString(t *testing.T) {
	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server.",
        "title": "Swagger Example API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:4000",
    "basePath": "/",
    "paths": {
        "/do-something": {
            "post": {
                "description": "Does something",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Call DoSomething",
                "parameters": [
                    {
                        "description": "My Struct",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.MyStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MyStruct"
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "main.MyStruct": {
            "type": "object",
            "properties": {
                "boolvar": {
                    "description": "boolean as a string",
                    "type": "string",
                    "example": "false"
                },
                "floatvar": {
                    "description": "float as a string",
                    "type": "string",
                    "example": "0"
                },
                "id": {
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                },
                "myint": {
                    "description": "integer as string",
                    "type": "string",
                    "example": "0"
                },
                "name": {
                    "type": "string",
                    "example": "poti"
                },
                "truebool": {
                    "description": "boolean as a string",
                    "type": "string",
                    "example": "true"
                }
            }
        }
    }
}`

	searchDir := "testdata/json_field_string"
	mainAPIFile := "main.go"
	p := New()
	err := p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParseSwaggerignoreForEmbedded(t *testing.T) {
	src := `
package main

type Child struct {
	ChildName string
}//@name Student

type Parent struct {
	Name string
	Child ` + "`swaggerignore:\"true\"`" + `
}//@name Teacher

// @Param request body Parent true "query params"
// @Success 200 {object} Parent
// @Router /test [get]
func Fun()  {

}
`
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, goparser.ParseComments)
	assert.NoError(t, err)

	p := New()
	p.packages.CollectAstFile("api", "api/api.go", f)
	p.packages.ParseTypes()
	err = p.ParseRouterAPIInfo("", f)
	assert.NoError(t, err)

	assert.NoError(t, err)
	teacher, ok := p.swagger.Definitions["Teacher"]
	assert.True(t, ok)

	name, ok := teacher.Properties["name"]
	assert.True(t, ok)
	assert.Len(t, name.Type, 1)
	assert.Equal(t, "string", name.Type[0])

	childName, ok := teacher.Properties["childName"]
	assert.False(t, ok)
	assert.Empty(t, childName)
}

func TestDefineTypeOfExample(t *testing.T) {
	var example interface{}
	var err error

	example, err = defineTypeOfExample("string", "", "example")
	assert.NoError(t, err)
	assert.Equal(t, example.(string), "example")

	example, err = defineTypeOfExample("number", "", "12.34")
	assert.NoError(t, err)
	assert.Equal(t, example.(float64), 12.34)

	example, err = defineTypeOfExample("boolean", "", "true")
	assert.NoError(t, err)
	assert.Equal(t, example.(bool), true)

	example, err = defineTypeOfExample("array", "", "one,two,three")
	assert.Error(t, err)
	assert.Nil(t, example)

	example, err = defineTypeOfExample("array", "string", "one,two,three")
	assert.NoError(t, err)
	arr := []string{}

	for _, v := range example.([]interface{}) {
		arr = append(arr, v.(string))
	}

	assert.Equal(t, arr, []string{"one", "two", "three"})

	example, err = defineTypeOfExample("object", "", "key_one:one,key_two:two,key_three:three")
	assert.Error(t, err)
	assert.Nil(t, example)

	example, err = defineTypeOfExample("object", "string", "key_one,key_two,key_three")
	assert.Error(t, err)
	assert.Nil(t, example)

	example, err = defineTypeOfExample("object", "oops", "key_one:one,key_two:two,key_three:three")
	assert.Error(t, err)
	assert.Nil(t, example)

	example, err = defineTypeOfExample("object", "string", "key_one:one,key_two:two,key_three:three")
	assert.NoError(t, err)
	obj := map[string]string{}

	for k, v := range example.(map[string]interface{}) {
		obj[k] = v.(string)
	}

	assert.Equal(t, obj, map[string]string{"key_one": "one", "key_two": "two", "key_three": "three"})

	example, err = defineTypeOfExample("oops", "", "")
	assert.Error(t, err)
	assert.Nil(t, example)
}