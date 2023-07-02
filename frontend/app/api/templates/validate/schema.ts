const schema = {
    "$schema": "http://json-schema.org/draft-06/schema#",
    "$ref": "#/definitions/Template",
    "definitions": {
        "Template": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "name": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "type": {
                    "type": "string"
                },
                "requirements": {
                    "$ref": "#/definitions/Requirements"
                },
                "regex": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Regex"
                    }
                },
                "output": {
                    "$ref": "#/definitions/Output"
                },
                "match": {
                    "$ref": "#/definitions/Match"
                },
                "script": {
                    "$ref": "#/definitions/Script"
                },
                "pre_script": {
                    "$ref": "#/definitions/Script"
                },
                "post_script": {
                    "$ref": "#/definitions/Script"
                }
            },
            "required": [
                "description",
                "name",
                "requirements",
                "tags",
                "type"
            ],
            "title": "Template"
        },
        "Match": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "filename": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "exclude": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            },
            "required": [
                "exclude",
                "filename"
            ],
            "title": "Match"
        },
        "Output": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "uniq": {
                    "type": "boolean"
                }
            },
            "required": [
                "uniq"
            ],
            "title": "Output"
        },
        "Script": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "language": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                }
            },
            "required": [
                "code",
                "language"
            ],
            "title": "Script"
        },
        "Regex": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "expression": {
                    "type": "string"
                },
                "group": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "references": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "validation": {
                    "$ref": "#/definitions/Validation"
                }
            },
            "required": [
                "description",
                "expression",
                "group",
                "references",
                "validation"
            ],
            "title": "Regex"
        },
        "Validation": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "tests": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Test"
                    }
                }
            },
            "required": [
                "tests"
            ],
            "title": "Validation"
        },
        "Test": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "input": {
                    "type": "string"
                },
                "want": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            },
            "required": [
                "input",
                "want"
            ],
            "title": "Test"
        },
        "Requirements": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "tools": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "pip": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "npm": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            },
            "required": [
            ],
            "title": "Requirements"
        }
    }
}

export default schema