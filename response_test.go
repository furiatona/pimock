package main

import (
  "testing"
)

func TestNewResponses(t *testing.T)  {
  type unit struct {
    mapdcvr map[string][]string
    method  string
    path    string
    header  map[string][]string
    urlq    map[string][]string
    result  int
  }

  cases := []unit{
    unit{map[string][]string{"healthz": []string{}},"GET", "healthz", map[string][]string{}, map[string][]string{"key": []string{"value"}}, 0},
    unit{map[string][]string{"healthz": []string{}},"GET", "healthz", map[string][]string{"Pimock-Sleep": []string{"1"}}, map[string][]string{"key": []string{"value"}}, 0},
    unit{map[string][]string{"healthz": []string{}},"GET", "healthz/", map[string][]string{}, map[string][]string{"key": []string{"value"}}, 0},
    unit{map[string][]string{"hallo": []string{}},"GET", "healthz", map[string][]string{}, map[string][]string{"key": []string{"value"}}, 1},
    unit{map[string][]string{"healthz/(?!1)": []string{}},"GET", "healthz", map[string][]string{}, map[string][]string{"key": []string{"value"}}, 2},
    unit{map[string][]string{"healthz": []string{}},"GET", "healthz", map[string][]string{"Pimock-Sleep": []string{"Value","Key"}}, map[string][]string{"key": []string{"value"}}, 2},

  }

  for _, c := range cases {
    switch c.result {
    case 0:
      res := NewResponse(c.method, c.path, c.header, c.urlq, c.mapdcvr)
      if res == nil {
        t.Errorf("The code did not give correct result")
      }
    case 1:
      res := NewResponse(c.method, c.path, c.header, c.urlq, c.mapdcvr)
      if res != nil {
        t.Errorf("The code did not give correct result")
      }
    case 2:
      t.Run("Panic Test", func(t *testing.T) {
        defer func() {
          if r := recover(); r == nil {
            t.Errorf("The code did not panic")
          }
        }()
        NewResponse(c.method, c.path, c.header, c.urlq, c.mapdcvr)
      })
    }
  }

}

func TestHeaderStatus(t *testing.T) {
  type unit struct {
    test    []string
    result  int
  }

  cases := []unit{
    unit{[]string{"HTTP/1.1 200 OK"}, 200},
    unit{[]string{"HTTP/1.1 OK"}, 0},
  }

  for _, c := range cases {
    res := Response{resp: c.test}
    if c.result != 0 {
      if res.getHeaderStatus() != c.result {
        t.Errorf("The code did not give correct result")
      }
    } else {
      defer func() {
          if r := recover(); r == nil {
              t.Errorf("The code did not panic")
          }
      }()
      res.getHeaderStatus()
    }
  }
}

func TestHeaderData(t *testing.T) {
  type unit struct {
    test    []string
    result  map[string]string
  }

  cases := []unit{
    unit{[]string{"HTTP/1.1 200 OK", "Content-Type: text/plain; charset=utf-8", ""}, map[string]string{"Content-Type": "text/plain; charset=utf-8"}},
    unit{[]string{"HTTP/1.1 200 OK", "Content-Type: application/json", ""}, map[string]string{"Content-Type": "application/json"}},
  }

  for _, c := range cases {
    res := Response{resp: c.test}
    headerData := res.getHeaderData()
    for i, v := range headerData {
      if c.result[i] != v {
        t.Errorf("The code did not give correct result")
      }
    }
  }

}

func TestBody(t *testing.T) {
  type unit struct {
    test    []string
    vars  map[string]string
    result  string
  }

  cases := []unit{
    unit{[]string{"", "{{request.path.[0]}}"}, map[string]string{"{{request.path.[0]}}": "healthz"}, "healthz"},
    unit{[]string{"", "OK"}, map[string]string{"{{request.path.[0]}}": "healthz"}, "OK"},
    unit{[]string{"OK"}, map[string]string{}, ""},
  }

  for _, c := range cases {
    res := Response{resp: c.test, variables: c.vars}
    if res.getBody() != c.result {
      t.Errorf("The code did not give correct result")
    }
  }

}
