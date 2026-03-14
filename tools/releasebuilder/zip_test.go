package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteDeterministicZip(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sourceDir := filepath.Join(root, "source")
	writeZipFixture(t, sourceDir)

	firstZip := filepath.Join(root, "first.zip")
	secondZip := filepath.Join(root, "second.zip")
	writeZipArchive(t, sourceDir, firstZip)
	writeZipArchive(t, sourceDir, secondZip)

	firstBytes := mustReadFile(t, firstZip)
	secondBytes := mustReadFile(t, secondZip)

	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("zip output is not deterministic")
	}

	reader, err := zip.OpenReader(firstZip)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer reader.Close()

	assertZipEntries(t, reader.File)
}

func TestWriteDeterministicZipRejectsSymlinks(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sourceDir := filepath.Join(root, "source")
	writeZipFixture(t, sourceDir)

	symlinkPath := filepath.Join(sourceDir, "README-link")
	if err := os.Symlink(filepath.Join(sourceDir, "README.md"), symlinkPath); err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}

	targetPath := filepath.Join(root, "release.zip")
	err := writeDeterministicZip(sourceDir, targetPath)
	if err == nil || !strings.Contains(err.Error(), "must not contain symlinks") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}

	if _, statErr := os.Stat(targetPath); !os.IsNotExist(statErr) {
		t.Fatalf("expected failed zip build to remove %s, got %v", targetPath, statErr)
	}
}

func TestValidateZipRelativePathRejectsEscapes(t *testing.T) {
	t.Parallel()

	if _, err := validateZipRelativePath(filepath.Join("..", "escape.txt")); err == nil {
		t.Fatal("expected escape path validation to fail")
	}
}

func writeZipFixture(t *testing.T, sourceDir string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(sourceDir, "bin"), 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}

	mustWriteFile(t, filepath.Join(sourceDir, "README.md"), []byte("readme\n"), 0o644)
	mustWriteFile(t, filepath.Join(sourceDir, "bin", "runecode-launcher"), []byte("launcher\n"), 0o755)
}

func writeZipArchive(t *testing.T, sourceDir, targetPath string) {
	t.Helper()

	if err := writeDeterministicZip(sourceDir, targetPath); err != nil {
		t.Fatalf("write zip %s: %v", targetPath, err)
	}
}

func assertZipEntries(t *testing.T, files []*zip.File) {
	t.Helper()

	if len(files) != 2 {
		t.Fatalf("expected 2 files in zip, got %d", len(files))
	}

	if files[0].Name != "source/README.md" {
		t.Fatalf("expected first entry source/README.md, got %q", files[0].Name)
	}
	if files[1].Name != "source/bin/runecode-launcher" {
		t.Fatalf("expected second entry source/bin/runecode-launcher, got %q", files[1].Name)
	}

	for _, file := range files {
		if !file.Modified.Equal(deterministicArchiveTime) {
			t.Fatalf("entry %q has modified time %v", file.Name, file.Modified)
		}
	}

	if mode := files[1].Mode().Perm(); mode != 0o755 {
		t.Fatalf("expected executable mode 0755, got %o", mode)
	}
	if mode := files[0].Mode().Perm(); mode != 0o644 {
		t.Fatalf("expected README mode 0644, got %o", mode)
	}
}

func mustWriteFile(t *testing.T, path string, contents []byte, mode os.FileMode) {
	t.Helper()

	if err := os.WriteFile(path, contents, mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return contents
}
