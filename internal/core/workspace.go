package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type WorkspaceKind string

const (
	WorkspaceGlobal  WorkspaceKind = "global"
	WorkspaceProject WorkspaceKind = "project"
)

var ErrNoProjectWorkspace = errors.New("project workspace not found")

type Workspace struct {
	Kind        WorkspaceKind
	RootPath    string
	ProjectPath string
}

type WorkspaceResolver struct {
	CWD     string
	HomeDir string
}

func NewWorkspaceResolver(cwd string) WorkspaceResolver {
	return WorkspaceResolver{CWD: cwd}
}

func (r WorkspaceResolver) ResolveInitial() (Workspace, error) {
	project, err := r.ResolveProject(false)
	if err == nil {
		return project, nil
	}
	if !errors.Is(err, ErrNoProjectWorkspace) {
		return Workspace{}, err
	}
	return r.ResolveGlobal()
}

func (r WorkspaceResolver) ResolveGlobal() (Workspace, error) {
	home, err := r.homeDir()
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{
		Kind:     WorkspaceGlobal,
		RootPath: filepath.Join(home, ".glyph"),
	}, nil
}

func (r WorkspaceResolver) ResolveProject(create bool) (Workspace, error) {
	cwd, err := r.cwd()
	if err != nil {
		return Workspace{}, err
	}
	globalRoot, err := r.globalRootPath()
	if err != nil {
		return Workspace{}, err
	}

	projectPath, found, err := findAncestorWithDirectory(cwd, ".glyph")
	if err != nil {
		return Workspace{}, err
	}
	if found && filepath.Clean(filepath.Join(projectPath, ".glyph")) == globalRoot {
		found = false
	}
	if !found {
		if !create {
			return Workspace{}, ErrNoProjectWorkspace
		}
		projectPath = cwd
		if err := os.MkdirAll(filepath.Join(projectPath, ".glyph"), 0o755); err != nil {
			return Workspace{}, err
		}
	}

	return Workspace{
		Kind:        WorkspaceProject,
		ProjectPath: projectPath,
		RootPath:    filepath.Join(projectPath, ".glyph"),
	}, nil
}

func (r WorkspaceResolver) cwd() (string, error) {
	if strings.TrimSpace(r.CWD) != "" {
		return filepath.Clean(r.CWD), nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Clean(cwd), nil
}

func (r WorkspaceResolver) homeDir() (string, error) {
	if strings.TrimSpace(r.HomeDir) != "" {
		return filepath.Clean(r.HomeDir), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Clean(home), nil
}

func (r WorkspaceResolver) globalRootPath() (string, error) {
	home, err := r.homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(home, ".glyph")), nil
}

func findAncestorWithDirectory(startDir, name string) (string, bool, error) {
	current := filepath.Clean(startDir)
	for {
		path := filepath.Join(current, name)
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return current, true, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", false, err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false, nil
		}
		current = parent
	}
}
