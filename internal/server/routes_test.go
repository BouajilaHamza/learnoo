package server_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

// changeCalls are the methods that change permissions or account status.
// Only the operator CLI (cmd/server) may call them; the HTTP application
// gains a path only through a reviewed change to this test (LRN-019).
var changeCalls = map[string]bool{
	"Grant": true, "Revoke": true, "GrantReviewer": true, "RevokeReviewer": true,
	"GrantDefaults": true, "Bootstrap": true,
}

// allowedCallers are directories whose non-test code may call changeCalls:
// the defining package and the operator CLI.
var allowedCallers = map[string]bool{
	"internal/permissions": true,
	"cmd/server":           true,
}

func TestNoHTTPPathChangesPermissions(t *testing.T) {
	root := filepath.Join("..", "..")
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == "node_modules" || strings.HasPrefix(d.Name(), ".")) && path != root {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, _ := filepath.Rel(root, filepath.Dir(path))
		if allowedCallers[filepath.ToSlash(rel)] {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok && changeCalls[sel.Sel.Name] {
				t.Errorf("%s: %s may change permissions outside the operator CLI", fset.Position(sel.Pos()), sel.Sel.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoPermissionRoutes(t *testing.T) {
	e := newEnv(t)
	token, csrf := e.mustLogin()
	for _, path := range []string{"/api/permissions", "/api/admin/permissions", "/api/users/1/permissions", "/api/admin/users/1"} {
		for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			w := e.do(m, path, map[string]string{"permission": "administrator"}, withCookie(token), withCSRF(csrf))
			if w.Code != http.StatusNotFound {
				t.Errorf("%s %s = %d, want 404", m, path, w.Code)
			}
		}
	}
}
