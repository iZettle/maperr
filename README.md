# maperr
[![][languagego img]][languagego]
[![][buildstatus img]][buildstatus]
[![][coverage img]][coverage]

`maperr` is a library that allow you to define a list of errors which you want to map to some other errors.

## Motivation

When writing a service which adopts a multi-layer architecture (e.g.: presentation, domain and storage) errors which are
returned in lower layers of the application are often wrapped with other errors. 
This not only add context but also allow you to assertion in your higher layer (e.g.: presentation) 
without having to directly perform checks for errors which are held withing the lower layers of your application (e.g.: storage).

This works fine if all your error will be mapped to a `5xx` error. However, sometimes some of this error are more user
errors which have been detected lower down the application layers which could fall in the `4xx` range.

This requires a mapping of errors, between layers, and if you have many errors to map, you would end up with one `if` statement
for each error. This library allow you to handle all of them in a single `if` statement.

## Mapping errors to status code

```go
var errMapper = maperr.NewMultiErr(
	maperr.NewHashableMapper().
		Append(domain.ErrOne, maperr.WithStatus("err one happened", http.StatusInternalServerError)).
		Append(domain.ErrTwo, maperr.WithStatus("err two happened", http.StatusBadRequest)).
		Append(domain.ErrThree, maperr.WithStatus("err three happened", http.StatusConflict)).
		Append(domain.ErrDeadlineExceeded, maperr.WithStatus("deadline exceeded", http.StatusBadRequest)).
		Append(domain.ErrCanceled, maperr.WithStatus("context was cancelled", http.StatusBadRequest)),
    
	maperr.NewIgnoreListMapper().
        Append(domain.ErrThatNeedsIgnored),
)

func (h Handler) Update(rw http.ResponseWriter, r *http.Request) {
    ...
    entity, err := h.Controller.Update(r.Context(), id)
    if mappedErr := errMapper.MappedWithStatus(err, maperr.WithStatusInternalServerError); mappedErr != nil {
    	 // mappedErr.Error() -> error to send in response
         // mappedErr.Status() -> status code for response
         // mappedErr.Unwrap() -> cause to log
    }
```

### Mapping errors to other errors

```go
    var errMapper = maperr.NewMultiErr(
        maperr.NewHashableMapper()
            Append(storage.ErrOne, ErrOne).
            Append(storage.ErrTwo, ErrTwo).
            Append(storage.ErrThree, ErrThree).
            Append(context.DeadlineExceeded, ErrDeadlineExceeded).
            Append(context.Canceled, ErrCanceled).
            Append(sql.ErrNoRows, ErrorUserNotFound))

    func (c Controller) Update(ctx context.Context, user User) error {
        ...
        err := s.storage.Update(user)
        if appendedErr := errMapper.Mapped(err, ErrSomethingBadHappened); appendedErr != nil {
            return appendedErr
        }
    }
```

[buildstatus img]:https://travis-ci.com/iZettle/maperr.svg?token=Gc7Chex1j1M4SzP7wjCm&branch=master
[buildstatus]:https://travis-ci.com/iZettle/maperr
[coverage img]:https://coveralls.io/repos/github/iZettle/maperr/badge.svg?branch=master&t=CxfFwY
[coverage]:https://coveralls.io/github/iZettle/maperr?branch=master
[languagego]:https://golang.org
[languagego img]:https://img.shields.io/badge/language-golang-77CDDD.svg?style=flat
