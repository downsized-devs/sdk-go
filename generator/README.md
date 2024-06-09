## How to use the generator

generator command

```bash
  go run main.go -entity_name=$entity_name -file_location=$file_location -api=$api
```

- `$entity_name` is the entity name that you want to generate and writen in camel case
- `$file_location` is the file location of generic-service in your local path
- `$api` is optional api that you want to generate. The available options are `create, edit, delete, get, activate` if this value is empty generator will generate all of the API type. Api option is written with `,` as separator
