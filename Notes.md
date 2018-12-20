# Dependency Management
* mit [go dep](https://golang.github.io/dep/)

* Installation über `apt-get go-dep`

```
dep ensure -add github.com/pkg/errors
```

# Wahl der Router Bibliothek
* chi-Router 

https://itnext.io/structuring-a-production-grade-rest-api-in-golang-c0229b3feedc



# Benötigte Routes



https://github.com/go-chi/chi/tree/master/_examples/todos-resource
https://github.com/go-chi/chi/blob/master/_examples/rest/main.go


# Chi-Codesnippets

`"github.com/go-chi/chi"`

```go
r := chi.Router()
r.Use(
    render.SetContentType(render.ContentTypeJSON),
    middleware.Logger,
    //...
)
```

`r.With().Get(...)` setzt eine Middleware bei einer spezifischen Route ein.

`r.Use()` setzt eine Middleware beim ganzen Router ein.

## Cors
Das Cors-Management basiert auf https://github.com/rs/cors. Das chi-Package ist lediglich ein Fork mit Integration für chi-Router

* beim Rendern von mehreren gleichen structs wird ein Array aus Rendern erstellt

## Rendern
* für das Ausgeben von JSON über die Render Methode muss ein struct mit allen Response Feldern erstellt werden. Dieses struct muss außerdem die Methode Render implementieren. Die Render Methode kann beispielsweise Fehlercodes mit der Response verbinden.

# Datenbank
* bongo als ODM für die Datenbank


# Docker

Für den Wechsel muss das ui Image neu gebildet werden.

```
docker-compose -f docker-compose.dev.yml up
```




```go

func increment(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "name")
	category := chi.URLParam(r, "category")

	user, err := db.FindUserByName(username)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	if category == "scheisse" {
		user, err = db.IncrementScheisse(user)
	} else if category == "schwierig" {
		user, err = db.IncrementSchwierig(user)
	} else if category == "problem" {
		user, err = db.IncrementProblem(user)
	} else {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("Can not find category")))
		return
	}

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
	}

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func decrement(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "name")
	category := chi.URLParam(r, "category")

	user, err := db.FindUserByName(username)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	if category == "scheisse" {
		user, err = db.DecrementScheisse(user)
	} else if category == "schwierig" {
		user, err = db.DecrementSchwierig(user)
	} else if category == "problem" {
		user, err = db.DecrementProblem(user)
	} else {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("Can not find category")))
		return
	}

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
	}

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

```



# Routes

`/users/sign-in`: Leifert beim richtigen Einloggen eine JWT, sondst eine Fehlermeldung
    * benutze die jwt Middleware von chi
`/users/sign-up`: Legt einen neuen Benutzer in einer Datenbank an

`/users`: Liefert den Benutzernamen aller Benutzer 
`/users/{name}` Liefert den Name eines bestimmten Benutzers



`/games/create`
    * erstellt aus einer Post request ein neues Spiel
`games/{name}` 
    * liefert Information zu einem bestimmten Spiel

`games/user/{name}`
    * liefert alle Spiele, bei denen der Nutzer mitspielt.



`games/enter/`
    * fügt einen neuen Nutzer dem Spiel hinzu. (POST)

`games/leave/`
    * entfernt einen Nutzer aus einem Spiel (POST)

	
```json
{
	"username": test,
	"game": "payword-backend"
}
```