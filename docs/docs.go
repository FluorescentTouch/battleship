// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-04-05 14:23:50.565002777 +0300 MSK m=+0.022600299

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Shamil Garatuev",
            "email": "garatuev@gmail.com"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/clear": {
            "post": {
                "description": "clear the battlefield",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "BattleField"
                ],
                "summary": "clear the battlefield",
                "responses": {
                    "200": {},
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/create-matrix": {
            "post": {
                "description": "create new battlefield with provided size",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "BattleField"
                ],
                "summary": "create new battlefield",
                "parameters": [
                    {
                        "description": "createParams",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/battlefield.CreateFieldRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/battlefield.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/battlefield.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ship": {
            "post": {
                "description": "add ships to battlefield\ninput params should be like this:\n\"A1 B2,C4 C6,E7 F8\" where first coordinate is one corner of ship, second - other.\nships can be square or rectangular\nships can't be placed on top of each other and near each other.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Ships"
                ],
                "summary": "add ships to battlefield",
                "parameters": [
                    {
                        "description": "coordinates",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/battlefield.AddShipsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/battlefield.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/battlefield.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "battlefield.AddShipsRequest": {
            "type": "object",
            "properties": {
                "coords": {
                    "type": "string"
                }
            }
        },
        "battlefield.CreateFieldRequest": {
            "type": "object",
            "properties": {
                "range": {
                    "type": "integer"
                }
            }
        },
        "battlefield.HTTPError": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "2.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Swagger Example API",
	Description: "This is a battleships game server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
