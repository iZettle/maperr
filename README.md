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

Controller layer

```go
    // ControllerErrors associates a Repository error with a controller layer error
    var ControllerErrors = maperr.NewMultiErr(
    	maperr.NewHashableMapper().
            Append(ErrRepositorySKURecipeNotFound, ErrControllerRecipeForSKUNotFound).
            Append(ErrRepositorySKURecipeAssociationCouldNotBeCreated, ErrControllerCouldNotAssociateSKUToRecipe).
            Append(ErrRepositorySKURecipeAssociationCouldNotBeUpdated, ErrControllerCouldNotRemoveSKUAssociationWithRecipe)

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
    var RepositoryErrors = maperr.NewMultiErr(
    	maperr.NewHashableMapper().
            Append(storage.ErrSKURecipeNotFound, ErrRepositorySKURecipeNotFound).
            Append(storage.ErrDatabaseSKURecipeQuerySelectFailed, ErrRepositorySKURecipeNotFound).
            Append(storage.ErrDatabaseSKURecipeQueryInsertFailed, ErrRepositorySKURecipeAssociationCouldNotBeCreated).
            Append(storage.ErrDatabaseSKURecipeQueryUpdateFailed, ErrRepositorySKURecipeAssociationCouldNotBeUpdated).
            Append(storage.ErrDatabaseSKURecipeQueryUpdateFailedNoRowsAffected, ErrRepositorySKURecipeAssociationCouldNotBeUpdate)

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
    var Errors = maperr.NewMultiErr(maperr.NewHashableMapper()
        Append(sql.ErrNoRows, storage.ErrorSKURecipeNotFound)

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

### List mapper

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

### Ignore list mapper

```go
    var Errors = maperr.NewMultiErr(
    	maperr.NewIgnoreListMapper().
            .Append(errors.New("pq: canceling statement due to user request"))
```

[buildstatus img]:https://travis-ci.com/iZettle/maperr.svg?token=Gc7Chex1j1M4SzP7wjCm&branch=master
[buildstatus]:https://travis-ci.com/iZettle/maperr
[coverage img]:https://coveralls.io/repos/github/iZettle/maperr/badge.svg?branch=master&t=CxfFwY
[coverage]:https://coveralls.io/github/iZettle/maperr?branch=master
[languagego]:https://golang.org
[languagego img]:https://img.shields.io/badge/language-golang-77CDDD.svg?style=flat
