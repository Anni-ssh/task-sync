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
        "/people": {
            "get": {
                "description": "Get all people",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "List People",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.People"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing person",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "Update People",
                "parameters": [
                    {
                        "description": "Person to update",
                        "name": "person",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.People"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to update person",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new person record. Passport number should be 6 digits and passport series should be 4 digits.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "Create a new people",
                "parameters": [
                    {
                        "description": "Details of the person to create",
                        "name": "people",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.People"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "ID of the created people",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "422": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/people/filter": {
            "get": {
                "description": "Get people based on filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "Get People by Filter",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Person ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Passport Series",
                        "name": "passport_series",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Passport Number",
                        "name": "passport_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Surname",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Patronymic",
                        "name": "patronymic",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Address",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.People"
                            }
                        }
                    },
                    "422": {
                        "description": "Failed to fetch people by filter",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/people/{peopleID}": {
            "get": {
                "description": "Get details of a people by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "Get People by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "People ID",
                        "name": "peopleID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.People"
                        }
                    },
                    "400": {
                        "description": "Failed to fetch person by ID",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a people by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "People"
                ],
                "summary": "Delete people",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "People ID",
                        "name": "peopleID",
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
                    "422": {
                        "description": "Invalid people ID",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to delete person",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/task": {
            "get": {
                "description": "Get list of all tasks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "List Tasks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Task"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing task. FORMAT TIME - RFC 3339 \"2024-08-01T08:00:00Z\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "Update Task",
                "parameters": [
                    {
                        "description": "Task to update",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.taskUpdate"
                        }
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
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new task. FORMAT TIME - RFC 3339 \"2024-08-01T08:00:00Z\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "Create Task",
                "parameters": [
                    {
                        "description": "Task to create",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task ID",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/task/update-people": {
            "put": {
                "description": "Update people associated with a task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "Update People in Task",
                "parameters": [
                    {
                        "description": "People and task to update",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.PeopleAndTask"
                        }
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
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/task/{taskID}": {
            "get": {
                "description": "Get a task by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "Get Task by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a task by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task"
                ],
                "summary": "Delete Task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskID",
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
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/time/end": {
            "post": {
                "description": "End recording time for a task. FORMAT TIME - RFC 3339 \"2024-08-01T08:00:00Z\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Time"
                ],
                "summary": "End Time Entry",
                "parameters": [
                    {
                        "description": "Task to end time entry for",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.timeTask"
                        }
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
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/time/spent": {
            "post": {
                "description": "Get time spent on tasks by a person within a specific time range. FORMAT TIME - RFC 3339 \"2024-08-01T08:00:00Z\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Time"
                ],
                "summary": "Task Time Spent",
                "parameters": [
                    {
                        "description": "People id and time range",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.peopleTimeRange"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TaskTimeSpent"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/time/start": {
            "post": {
                "description": "Start recording time for a task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Time"
                ],
                "summary": "Start Time Entry",
                "parameters": [
                    {
                        "description": "Task to start time entry for",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.timeTask"
                        }
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
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.People": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "passport_number": {
                    "type": "integer"
                },
                "passport_series": {
                    "type": "integer"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "entities.Task": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "timeEntry": {
                    "$ref": "#/definitions/entities.TimeEntry"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entities.TaskTimeSpent": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "people_id": {
                    "type": "integer"
                },
                "surname": {
                    "type": "string"
                },
                "task_id": {
                    "type": "integer"
                },
                "task_title": {
                    "type": "string"
                },
                "time_spent": {
                    "type": "string"
                }
            }
        },
        "entities.TimeEntry": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                },
                "end_time": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "people_id": {
                    "type": "integer"
                },
                "start_time": {
                    "type": "string"
                }
            }
        },
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.PeopleAndTask": {
            "type": "object",
            "properties": {
                "peopleID": {
                    "type": "integer"
                },
                "taskID": {
                    "type": "integer"
                }
            }
        },
        "handler.peopleTimeRange": {
            "type": "object",
            "properties": {
                "end_time": {
                    "type": "string"
                },
                "people_id": {
                    "type": "integer"
                },
                "start_time": {
                    "type": "string"
                }
            }
        },
        "handler.taskUpdate": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "task_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "handler.timeTask": {
            "type": "object",
            "properties": {
                "task_id": {
                    "type": "integer"
                },
                "time": {
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
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "TaskSync API",
	Description:      "API Server for Task Tracking",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
