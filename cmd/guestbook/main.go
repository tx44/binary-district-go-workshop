package main

import (
  "flag"
  "html/template"
  "io"
  "log"
  "net/http"
)

var (
  addr = flag.String("addr", "127.0.0.1:8081", "addr to bind to")
)

const index = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
    <form method="POST" action="/">
      <input type="text" name="name" />
      <textarea name="message"></textarea>
      <input type="submit" />
    </form>

		{{range .Items}}
      <div>{{ . }}</div>
    {{else}}
    <div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

type Index struct{
  Title string
  Items []string
}

type Server struct {
  Pages map[string]func(io.Writer) error
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  log.Printf("got request %s", req.URL)
  if req.Method == "POST" {
    if err := req.ParseForm(); err != nil {
      res.WriteHeader(400)
      return
    }
    for key, values := range req.PostForm {
      log.Printf("key=%q values=%v", key, values)
    }
  }

  fn, ok := s.Pages[req.URL.Path]
  if !ok {
    res.WriteHeader(404)
    return
  }
  if err := fn(res); err != nil {
    res.WriteHeader(500)
  }
  //res.Write([]byte("hi there!"))
}

func main() {
  tmpl, err := template.New("index").Parse(index)
  if err != nil {
    log.Fatal(err)
  }

  s := Server{
    Pages: map[string]func(io.Writer) error {
      "/": func(w io.Writer) error {
        tmpl.Execute(w, Index{
          Title: "My Cool Guetbook",
        })
        return nil
      },
    },
  }

  log.Fatal(http.ListenAndServe(*addr, &s))
}
