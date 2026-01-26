# Updating Gameboy Go Core Library

Everytime a new feature is added to the Gameboy Go Core Library, the following steps should be followed to make the feature available to the users through the gameboy-go server.

## gameboy-go repository

Once the feature/bug fix is ready to be released, the following steps should be followed from the gameboy-go repository.

### step 1: commit and push changes

Nothing special here, just commit all the changes locally and push the changes to the remote repository.

```bash
git add .
git commit -m "commit message"
git push
```

### step 2: tag the version and push the tag

In order for the new module version to be referenced in the gameboy-go-server repository, the version should be tagged and pushed to the remote repository.

```bash
git tag v0.4.18 // tag the version
git push -tags  // push the tag to the remote repository
```

## gameboy-go-server repository

### step 1: get the new version of the gameboy-go module

In order to use the newly pushed/tagged version of the gameboy-go module, we can run the following command

```bash
go get github.com/codefrite/gameboy-go@v0.4.18
```

With this step done, the version used in the go.mod file should have been updated to the newer version.

```go
module github.com/codefrite/gameboy-go-server

go 1.23.1

require (
	github.com/codefrite/gameboy-go v0.4.18
	github.com/gorilla/websocket v1.5.3
)
```

Another way of proceding would be to manually update the go.mod file with the new version of the gameboy-go module and run the following command:

```bash
go mod tidy
```
