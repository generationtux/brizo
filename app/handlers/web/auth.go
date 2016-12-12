package web

import "net/http"

const htmlAddUser = `<!doctype html><html><body>
<form action="/api/v1/users" method="post">
  <label>Add user:</label>
  <input type="text" name="username">
</form>
</body></html>
`

func AuthAddNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlAddUser))
}
