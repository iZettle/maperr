# maperr
[![Build Status](https://travis-ci.org/iZettle/maperr.svg?branch=master)](https://travis-ci.org/iZettle/maperr)
[![Coverage Status](https://coveralls.io/repos/github/iZettle/maperr/badge.svg)](https://coveralls.io/github/iZettle/maperr)

<!-- vim-markdown-toc GFM -->

* [Code Owners](#code-owners)
* [Hashable mapper](#hashable-mapper)
* [List mapper](#list-mapper)
* [Example usage:](#example-usage)

<!-- vim-markdown-toc -->

Small library that allow to separate map errors through layers

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

See the example below:

```go
    ErrFoo := FooError{}.SomeCustomBehaviour("foo")

    var Errors = maperr.NewMultiErr(
        maperr.NewListMapper().
            Appendf("element with %d was not found", ErrBar). // wraps the error in a error type which holds the format
            Append(ErrFoo, ErrBar), // add the error as it is
    )
    
    // maperr.Errorf wraps the error in a type which holds the format
    // this means that the mapper can match when the format is the same
    err = maperr.Errorf("element with %d was not found", 12345)
    
    if appendedErr := Errors.Mapped(err, ErrFoo); appendedErr != nil {
        return nil, appendedErr
    }
```

## Example usage:

Handler layer
```go
    // HTTPHandlerErrors associates a ControllerManager error with a http api handler layer error
    var HTTPHandlerErrors = maperr.NewMultiErr(maperr.HashableMapper{
    	skurecipe.ErrControllerRecipeForSKUNotFound:                   maperr.WithStatus(errorTextRecipeForSKUNotFound, http.StatusNotFound),
    	skurecipe.ErrControllerSKURecipeAssociationAlreadyExists:      maperr.WithStatus(errorTextRecipeForSKUAlreadyExists, http.StatusBadRequest),
    	skurecipe.ErrControllerBusinessIDMismatch:                     maperr.WithStatus(errorTextResourceBusinessMismatch, http.StatusUnauthorized),
    })

    func (h *HTTPHandler) Get() (*Foo, error) {
        model, err := h.controller.Get(id, businessID)
        if errWithStatus := HTTPHandlerErrors.LastMappedWithStatus(err); errWithStatus != nil {
            return jsonhandler.NewResponseWithError(errWithStatus.Status(), errWithStatus.Error(), &err)
        }
        if err != nil {
            return jsonhandler.ServerError(err)
        }
    }
```

Controller layer
```go
    // ControllerErrors associates a Repository error with a controller layer error
    var ControllerErrors = maperr.NewMultiErr(maperr.HashableMapper{
        ErrRepositorySKURecipeNotFound:                     ErrControllerRecipeForSKUNotFound,
        ErrRepositorySKURecipeAssociationCouldNotBeCreated: ErrControllerCouldNotAssociateSKUToRecipe,
        ErrRepositorySKURecipeAssociationCouldNotBeUpdated: ErrControllerCouldNotRemoveSKUAssociationWithRecipe,
    })

    func (r *Controller) Get() (*Foo, error) {
        ...
        if appendedErr := ControllerErrors.Mapped(err, ErrorControllerRecipeForSKUNotFound); appendedErr != nil {
            return nil, appendedErr
        }
    }
```

Repository layer
```go
    // RepositoryErrors associates a storage error with a Repository layer error
    var RepositoryErrors = maperr.NewMultiErr(maperr.HashableMapper{
        storage.ErrSKURecipeNotFound:                                ErrRepositorySKURecipeNotFound,
        storage.ErrDatabaseSKURecipeQuerySelectFailed:               ErrRepositorySKURecipeNotFound,
        storage.ErrDatabaseSKURecipeQueryInsertFailed:               ErrRepositorySKURecipeAssociationCouldNotBeCreated,
        storage.ErrDatabaseSKURecipeQueryUpdateFailed:               ErrRepositorySKURecipeAssociationCouldNotBeUpdated,
        storage.ErrDatabaseSKURecipeQueryUpdateFailedNoRowsAffected: ErrRepositorySKURecipeAssociationCouldNotBeUpdated,
    })

    func (r *Repository) Get() (*Foo, error) {
        ...
        if appendedErr := RepositoryErrors.Mapped(err, ErrorRepositorySKURecipeNotFound); appendedErr != nil {
            return nil, appendedErr
        }
    }
```

Storage layer
```go
    // Errors associates a possible error with a storage layer error
    var Errors = maperr.NewMultiErr(maperr.HashableMapper{
        sql.ErrNoRows: storage.ErrorSKURecipeNotFound,
    })

    func (s *Storage) Get() (*Foo, error) {
        ...
        err := s.db.Get(model, query, args...)

        // if the error is sql.ErrNoRows, wraps storage.ErrorSKURecipeNotFound
        // otherwise wraps storage.ErrorDatabaseQuerySelectFailed
        appendedErr := Errors.Mapped(err, storage.ErrorDatabaseQuerySelectFailed)
        if appendedErr != nil {
            return nil, appendedErr
        }
    }
```
