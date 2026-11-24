package server_test

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const module = "github.com/hadden-me/learnoo"

// changeMethods are the methods that change permissions or account status,
// by the package path of the Store type that declares them.
var changeMethods = map[string]map[string]bool{
	module + "/internal/permissions": {"Grant": true, "Revoke": true, "GrantReviewer": true, "RevokeReviewer": true, "GrantDefaults": true, "Bootstrap": true},
	module + "/internal/accounts":    {"Disable": true, "Enable": true},
}

// allowedChangeFiles may call changeMethods: the operator CLI and the
// defining files. cmd/server/main.go also starts the HTTP server, but its
// handlers live in internal/server, and any call inside an HTTP handler is
// flagged in every file anyway.
var allowedChangeFiles = map[string]bool{
	"cmd/server/main.go":                  true,
	"internal/permissions/permissions.go": true,
	"internal/accounts/accounts.go":       true,
}

// sqlOwners lists, per protected write, the only package that may contain
// it in a SQL string literal.
var sqlOwners = []struct {
	re    *regexp.Regexp
	owner string
}{
	{regexp.MustCompile(`(insert\s+into|update|delete\s+from|truncate(\s+table)?)\s+(user_permissions|reviewer_assignments|permission_audit)\b`), "internal/permissions"},
	{regexp.MustCompile(`(insert\s+into|update|delete\s+from|truncate(\s+table)?)\s+account_status_audit\b`), "internal/accounts"},
	{regexp.MustCompile(`update\s+users\b[^;]*\bdisabled_at\b`), "internal/accounts"},
}

type goPackage struct {
	path, dir string
	files     []string
}

func goList(t *testing.T, root string, args ...string) []string {
	t.Helper()
	cmd := exec.Command("go", append([]string{"list"}, args...)...)
	cmd.Dir = root
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list %v: %v", args, err)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

// TestNoHTTPPathChangesPermissions type-checks every non-test package and
// fails if code outside the operator CLI and the owning packages changes
// permissions or account status: calls to the change methods of
// permissions.Store or accounts.Store, any use of permissions.Operator, any
// such call or use inside an HTTP handler, or SQL literals that write the
// permission, audit or account status tables directly.
func TestNoHTTPPathChangesPermissions(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	exports := map[string]string{}
	for _, line := range goList(t, root, "-deps", "-export", "-f", "{{.ImportPath}}={{.Export}}", "./...") {
		if path, file, ok := strings.Cut(line, "="); ok && file != "" {
			exports[path] = file
		}
	}
	var pkgs []goPackage
	for _, line := range goList(t, root, "-f", "{{.ImportPath}}|{{.Dir}}|{{join .GoFiles \",\"}}", "./...") {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 && parts[2] != "" {
			pkgs = append(pkgs, goPackage{parts[0], parts[1], strings.Split(parts[2], ",")})
		}
	}
	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		file, ok := exports[path]
		if !ok {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(file)
	})
	checked := 0
	for _, pkg := range pkgs {
		var files []*ast.File
		for _, name := range pkg.files {
			f, err := parser.ParseFile(fset, filepath.Join(pkg.dir, name), nil, parser.SkipObjectResolution)
			if err != nil {
				t.Fatal(err)
			}
			files = append(files, f)
		}
		info := &types.Info{Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}}
		if _, err := (&types.Config{Importer: imp}).Check(pkg.path, fset, files, info); err != nil {
			t.Fatalf("type-check %s: %v", pkg.path, err)
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(pkg.path, module), "/")
		for _, f := range files {
			checkFile(t, fset, root, rel, f, info)
			checked++
		}
	}
	if checked < 10 {
		t.Fatalf("only %d files checked", checked)
	}
}

func checkFile(t *testing.T, fset *token.FileSet, root, pkgRel string, f *ast.File, info *types.Info) {
	t.Helper()
	pos := func(p token.Pos) string {
		position := fset.Position(p)
		rel, _ := filepath.Rel(root, position.Filename)
		return fmt.Sprintf("%s:%d", filepath.ToSlash(rel), position.Line)
	}
	fileRel, _ := filepath.Rel(root, fset.Position(f.Pos()).Filename)
	fileRel = filepath.ToSlash(fileRel)

	// Source ranges of functions that take an http.ResponseWriter.
	var handlers [][2]token.Pos
	ast.Inspect(f, func(n ast.Node) bool {
		var ft *ast.FuncType
		switch fn := n.(type) {
		case *ast.FuncDecl:
			ft = fn.Type
		case *ast.FuncLit:
			ft = fn.Type
		}
		if ft != nil && ft.Params != nil {
			for _, field := range ft.Params.List {
				if sel, ok := field.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "ResponseWriter" {
					if obj := info.Uses[sel.Sel]; obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "net/http" {
						handlers = append(handlers, [2]token.Pos{n.Pos(), n.End()})
					}
				}
			}
		}
		return true
	})
	inHandler := func(p token.Pos) bool {
		for _, h := range handlers {
			if p >= h[0] && p < h[1] {
				return true
			}
		}
		return false
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			if sel := info.Selections[x]; sel != nil && sel.Kind() == types.MethodVal {
				recv := sel.Recv()
				if p, ok := recv.(*types.Pointer); ok {
					recv = p.Elem()
				}
				if named, ok := recv.(*types.Named); ok && named.Obj().Name() == "Store" && named.Obj().Pkg() != nil &&
					changeMethods[named.Obj().Pkg().Path()][x.Sel.Name] {
					if !allowedChangeFiles[fileRel] || inHandler(x.Pos()) {
						t.Errorf("%s: %s.%s changes permissions or account status outside the operator CLI", pos(x.Pos()), named.Obj().Name(), x.Sel.Name)
					}
				}
			}
		case *ast.Ident:
			if obj := info.Uses[x]; obj != nil && obj.Name() == "Operator" && obj.Pkg() != nil &&
				obj.Pkg().Path() == module+"/internal/permissions" {
				if (fileRel != "cmd/server/main.go" && fileRel != "internal/permissions/permissions.go") || inHandler(x.Pos()) {
					t.Errorf("%s: permissions.Operator used outside the operator CLI", pos(x.Pos()))
				}
			}
		case *ast.BasicLit:
			if x.Kind != token.STRING {
				return true
			}
			s, err := strconv.Unquote(x.Value)
			if err != nil {
				return true
			}
			s = strings.ToLower(s)
			for _, o := range sqlOwners {
				if o.re.MatchString(s) && pkgRel != o.owner {
					t.Errorf("%s: SQL write %q belongs in %s", pos(x.Pos()), o.re.FindString(s), o.owner)
				}
			}
		}
		return true
	})
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
