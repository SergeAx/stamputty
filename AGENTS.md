# Rules for coding agents

## Common
 
* Use conventional commit messages
* Do not use emoji
* Make small atomic commits
* Create tests for business logic
* Always add a trailing newline at the end of text files

## Language specific

* Always use `go fmt`
* Use `go mod tidy` after dependencies change 
* Never ignore return error values, use blank identifier if error handling is irrelevant
