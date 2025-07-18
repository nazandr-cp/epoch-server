// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://lend.fam/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://lend.fam/support",
            "email": "support@lend.fam"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/epochs/distribute": {
            "post": {
                "description": "Initiates the distribution of subsidies for the current epoch",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "epochs"
                ],
                "summary": "Distribute subsidies",
                "responses": {
                    "202": {
                        "description": "Subsidy distribution accepted",
                        "schema": {
                            "$ref": "#/definitions/handlers.DistributeSubsidiesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/epochs/force-end": {
            "post": {
                "description": "Forcibly ends an epoch with zero yield distribution",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "epochs"
                ],
                "summary": "Force end epoch",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Epoch ID to force end",
                        "name": "epochId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Epoch force end accepted",
                        "schema": {
                            "$ref": "#/definitions/handlers.ForceEndEpochResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request - missing or invalid epochId",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/epochs/start": {
            "post": {
                "description": "Initiates the start of a new epoch for yield distribution",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "epochs"
                ],
                "summary": "Start epoch",
                "responses": {
                    "202": {
                        "description": "Epoch start accepted",
                        "schema": {
                            "$ref": "#/definitions/handlers.StartEpochResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{address}/merkle-proof": {
            "get": {
                "description": "Generates a merkle proof for a user's current earnings",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user merkle proof",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User wallet address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vault address (optional, uses default if not provided)",
                        "name": "vault",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Merkle proof generated successfully",
                        "schema": {
                            "$ref": "#/definitions/merkle.UserMerkleProofResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request - invalid address",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{address}/merkle-proof/epoch/{epochNumber}": {
            "get": {
                "description": "Generates a merkle proof for a user's earnings at a specific epoch",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get historical merkle proof",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User wallet address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Epoch number",
                        "name": "epochNumber",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Vault address (optional, uses default if not provided)",
                        "name": "vault",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Historical merkle proof generated successfully",
                        "schema": {
                            "$ref": "#/definitions/merkle.UserMerkleProofResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request - invalid address or epoch",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User or epoch not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{address}/total-earned": {
            "get": {
                "description": "Retrieves the total amount earned by a user across all epochs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user total earned",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User wallet address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User earnings information",
                        "schema": {
                            "$ref": "#/definitions/epoch.UserEarningsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request - invalid address",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Returns the current health status of the epoch server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "$ref": "#/definitions/handlers.HealthResponse"
                        }
                    },
                    "503": {
                        "description": "Service is unhealthy",
                        "schema": {
                            "$ref": "#/definitions/handlers.HealthResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "epoch.UserEarningsResponse": {
            "type": "object",
            "properties": {
                "calculatedAt": {
                    "type": "integer"
                },
                "dataTimestamp": {
                    "description": "Timestamp used for calculations",
                    "type": "integer"
                },
                "totalEarned": {
                    "type": "string"
                },
                "userAddress": {
                    "type": "string"
                },
                "vaultAddress": {
                    "type": "string"
                }
            }
        },
        "handlers.DistributeSubsidiesResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Subsidy distribution initiated successfully"
                },
                "status": {
                    "type": "string",
                    "example": "accepted"
                },
                "vaultID": {
                    "type": "string",
                    "example": "0x1234567890123456789012345678901234567890"
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "details": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "handlers.ForceEndEpochResponse": {
            "type": "object",
            "properties": {
                "epochId": {
                    "type": "integer",
                    "example": 1
                },
                "message": {
                    "type": "string",
                    "example": "Force end epoch initiated successfully"
                },
                "status": {
                    "type": "string",
                    "example": "accepted"
                },
                "vaultID": {
                    "type": "string",
                    "example": "0x1234567890123456789012345678901234567890"
                }
            }
        },
        "handlers.HealthResponse": {
            "type": "object",
            "properties": {
                "checks": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string",
                    "example": "ok"
                }
            }
        },
        "handlers.StartEpochResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Epoch start initiated successfully"
                },
                "status": {
                    "type": "string",
                    "example": "accepted"
                }
            }
        },
        "merkle.UserMerkleProofResponse": {
            "type": "object",
            "properties": {
                "epochNumber": {
                    "type": "string"
                },
                "generatedAt": {
                    "type": "integer"
                },
                "leafIndex": {
                    "type": "integer"
                },
                "merkleProof": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "merkleRoot": {
                    "type": "string"
                },
                "totalEarned": {
                    "type": "string"
                },
                "userAddress": {
                    "type": "string"
                },
                "vaultAddress": {
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
	Schemes:          []string{"http", "https"},
	Title:            "Epoch Server API",
	Description:      "Epoch Server for managing NFT collection-backed lending epochs, subsidies, and merkle proofs",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
