# maperr
[![][languagego img]][languagego]
[![][buildstatus img]][buildstatus]
[![][coverage img]][coverage]

<!-- vim-markdown-toc GFM -->

* [Code Owners](#code-owners)
* [Hashable mapper](#hashable-mapper)
* [List mapper](#list-mapper)
* [Ignore list mapper](#ignore-list-mapper)
* [Usages:](#usages)
	* [Hashable mapper](#hashable-mapper-1)
		* [Errors with status with hashable mapper](#errors-with-status-with-hashable-mapper)
	* [List mapper](#list-mapper-1)
	* [Ignore list mapper](#ignore-list-mapper-1)

<!-- vim-markdown-toc -->

## Code Owners
| Owner                             | Slack Channel          |
|-----------------------------------|------------------------|
| `mse-understand`                  | `#edinburgh-understand`|

## Hashable mapper
`HashableMapper` only works when the mapped error is a comparable https://golang.org/ref/spec#Comparison_operators

This is because is defined as `type HashableMapper map[error]error`

Use this when you are mapping simple errors, is the fastest mapper when matching the errors.

## List mapper
`ListMapper` can be used in combination of the standard `HashableMapper` to allow to map any error that implements
`maperr.Error`

`maperr.Errorf` returns an error which implements `maperr.Error` and when used in combination with `ListMapper` allow 
you to match error's formats

## Ignore list mapper

`IgnoreListMapper` allow to define a list of errors that should be ignored therefore if they are found in the last error
`Mapped()` will return `nil`

## Usages:

### Hashable mapper

```go
    // storageErrors associates a possible error with a storage layer error
    var storageErrors = maperr.NewMultiErr(maperr.NewHashableMapper()
        Append(sql.ErrNoRows, storage.ErrorSKURecipeNotFound)

    func (s *Storage) Get() (*Foo, error) {
        ...
        err := s.db.Get(model, query, args...)

        // if the error is sql.ErrNoRows, wraps storage.ErrorSKURecipeNotFound
        // otherwise wraps storage.ErrorDatabaseQuerySelectFailed
        appendedErr := storageErrors.Mapped(err, storage.ErrorDatabaseQuerySelectFailed)
        if appendedErr != nil {
            return nil, appendedErr
        }
    }
```


#### Errors with status with hashable mapper

The hashable mapper supports errors that are decorated with a status code by using `maperr.WithStatus`

```go
var handlerErrors = maperr.NewMultiErr(
	maperr.NewHashableMapper().
		Append(layoutsetview.ErrLayoutSetViewControllerCancelledContext, maperr.WithStatus(errTextCancelledRequest, http.StatusBadRequest)))

func (h Handler) GetByTerminal(rw http.ResponseWriter, r *http.Request) jsonhandler.JSONResponse {
    ...
    layoutSet, err := h.Controller.GetListByTerminal(r.Context(), siteID, terminalID)
    
    if mappedErr := errMapper.MappedWithStatus(err, maperr.WithStatusInternalServerError); mappedErr != nil {
         return jsonhandler.NewLoggableResponseFromErrorWithStatus(mappedErr)
    }
```

### List mapper

```go
    var errTextElementNotFound := "element with %d was not found"

    var repositoryErrors = maperr.NewMultiErr(
        maperr.NewListMapper().
            Appendf(errTextElementNotFound, ErrBar). // wraps the error in a error type which holds the format
            Append(ErrFoo, ErrBar), // add the error as it is
    )
    
    // maperr.Errorf wraps the error in a type which holds the format
    // this means that the mapper can match when the format is the same
    err = maperr.Errorf(errTextElementNotFound, 12345)
    
    if appendedErr := repositoryErrors.Mapped(err, ErrFoo); appendedErr != nil {
        return nil, appendedErr
    }
```

### Ignore list mapper

```go
    var Errors = maperr.NewMultiErr(
    	maperr.NewIgnoreListMapper().
            .Append(errors.New("this need ignored, use sparingly"))
```

[buildstatus img]:https://travis-ci.com/iZettle/maperr.svg?token=Gc7Chex1j1M4SzP7wjCm&branch=master
[buildstatus]:https://travis-ci.com/iZettle/maperr
[coverage img]:https://coveralls.io/repos/github/iZettle/maperr/badge.svg?branch=master&t=CxfFwY
[coverage]:https://coveralls.io/github/iZettle/maperr?branch=master
[languagego]:https://golang.org
[languagego img]:https://img.shields.io/badge/language-golang-77CDDD.svg?style=flat
