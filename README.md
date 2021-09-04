# jv

Rich JSON valiator

# Installation

```
$ go get https://github.com/teru01/jv
```

# Usage

```
$ jv [FILE]
```

# Feature

- validate JSON including `"${Variable}", ${Number}`
- extremely simple

# Example

sample.json.template
```
{
    "name": "main",
    "image": "${IMAGE_NAME}",
    "essential": true,
    "portMappings": [
        {
            "protocol": "tcp",
            "containerPort": ${PORT},
        }
    ]
}
```

It is tedious to collect all the environment variables and then validate them with jq, etc. jv provides a simple syntax check for such cases.

```
$ jv sample.json.template
syntax error at line 9

        }
          ^^^
invalid character '}' looking for beginning of object key string
```

