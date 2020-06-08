package keystore

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"

	ci "github.com/libp2p/go-libp2p-core/crypto"
)

type rr struct{}

func (rr rr) Read(b []byte) (int, error) {
	return rand.Read(b)
}

func privKeyOrFatal(t *testing.T) ci.PrivKey {
	priv, _, err := ci.GenerateEd25519Key(rr{})
	if err != nil {
		t.Fatal(err)
	}
	return priv
}

func TestKeystoreBasics(t *testing.T) {
	tdir, err := ioutil.TempDir("", "keystore-test")
	if err != nil {
		t.Fatal(err)
	}

	ks, err := NewFSKeystore(tdir)
	if err != nil {
		t.Fatal(err)
	}

	l, err := ks.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 0 {
		t.Fatal("expected no keys")
	}

	k1 := privKeyOrFatal(t)
	k2 := privKeyOrFatal(t)
	k3 := privKeyOrFatal(t)
	k4 := privKeyOrFatal(t)

	err = ks.Put("foo", k1)
	if err != nil {
		t.Fatal(err)
	}

	err = ks.Put("bar", k2)
	if err != nil {
		t.Fatal(err)
	}

	l, err = ks.List()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(l)
	if l[0] != "bar" || l[1] != "foo" {
		t.Fatal("wrong entries listed")
	}

	if err := assertDirContents(tdir, []string{"foo", "bar"}); err != nil {
		t.Fatal(err)
	}

	err = ks.Put("foo", k3)
	if err == nil {
		t.Fatal("should not be able to overwrite key")
	}

	if err := assertDirContents(tdir, []string{"foo", "bar"}); err != nil {
		t.Fatal(err)
	}

	exist, err := ks.Has("foo")
	if !exist {
		t.Fatal("should know it has a key named foo")
	}
	if err != nil {
		t.Fatal(err)
	}

	exist, err = ks.Has("nonexistingkey")
	if exist {
		t.Fatal("should know it doesn't have a key named nonexistingkey")
	}
	if err != nil {
		t.Fatal(err)
	}

	if err := ks.Delete("bar"); err != nil {
		t.Fatal(err)
	}

	if err := assertDirContents(tdir, []string{"foo"}); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("beep", k3); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("boop", k4); err != nil {
		t.Fatal(err)
	}

	if err := assertDirContents(tdir, []string{"foo", "beep", "boop"}); err != nil {
		t.Fatal(err)
	}

	if err := assertGetKey(ks, "foo", k1); err != nil {
		t.Fatal(err)
	}

	if err := assertGetKey(ks, "beep", k3); err != nil {
		t.Fatal(err)
	}

	if err := assertGetKey(ks, "boop", k4); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("..///foo/", k1); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("", k1); err == nil {
		t.Fatal("shouldn't be able to put a key with no name")
	}

	if err := ks.Put(".foo", k1); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidKeyFiles(t *testing.T) {
	tdir, err := ioutil.TempDir("", "keystore-test")

	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tdir)

	ks, err := NewFSKeystore(tdir)
	if err != nil {
		t.Fatal(err)
	}

	key := privKeyOrFatal(t)

	bytes, err := key.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	encodedName, err := encode("valid")
<<<<<<< HEAD
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(ks.dir, encodedName), bytes, 0644)
=======
>>>>>>> master
	if err != nil {
		t.Fatal(err)
	}

<<<<<<< HEAD
=======
	err = ioutil.WriteFile(filepath.Join(ks.dir, encodedName), bytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

>>>>>>> master
	err = ioutil.WriteFile(filepath.Join(ks.dir, "z.invalid"), bytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	l, err := ks.List()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(l)
	if len(l) != 1 {
		t.Fatal("wrong entry count")
	}

	if l[0] != "valid" {
		t.Fatal("wrong entries listed")
	}

	exist, err := ks.Has("valid")
	if !exist {
		t.Fatal("should know it has a key named valid")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNonExistingKey(t *testing.T) {
	tdir, err := ioutil.TempDir("", "keystore-test")
	if err != nil {
		t.Fatal(err)
	}

	ks, err := NewFSKeystore(tdir)
	if err != nil {
		t.Fatal(err)
	}

	k, err := ks.Get("does-it-exist")
	if err != ErrNoSuchKey {
		t.Fatalf("expected: %s, got %s", ErrNoSuchKey, err)
	}
	if k != nil {
		t.Fatalf("Get on nonexistant key should give nil")
	}
}

func TestMakeKeystoreNoDir(t *testing.T) {
	_, err := NewFSKeystore("/this/is/not/a/real/dir")
	if err == nil {
		t.Fatal("shouldnt be able to make a keystore in a nonexistant directory")
	}
}

func assertGetKey(ks Keystore, name string, exp ci.PrivKey) error {
	outK, err := ks.Get(name)
	if err != nil {
		return err
	}

	if !outK.Equals(exp) {
		return fmt.Errorf("key we got out didn't match expectation")
	}

	return nil
}

func assertDirContents(dir string, exp []string) error {
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(finfos) != len(exp) {
		return fmt.Errorf("expected %d directory entries", len(exp))
	}

	var names []string
	for _, fi := range finfos {
		decodedName, err := decode(fi.Name())
		if err != nil {
			return err
		}
		names = append(names, decodedName)
	}

	sort.Strings(names)
	sort.Strings(exp)
	if len(names) != len(exp) {
		return fmt.Errorf("directory had wrong number of entries in it")
	}

	for i, v := range names {
		if v != exp[i] {
			return fmt.Errorf("had wrong entry in directory")
		}
	}
	return nil
}

func TestEncodedKeystoreBasics(t *testing.T) {
	tdir, err := ioutil.TempDir("", "encoded-keystore-test")
	if err != nil {
		t.Fatal(err)
	}

	ks, err := NewEncodedFSKeystore(tdir)
	if err != nil {
		t.Fatal(err)
	}

	l, err := ks.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 0 {
		t.Fatal("expected no keys")
	}

	k1 := privKeyOrFatal(t)
	k1Name, err := encode("foo")
	if err != nil {
		t.Fatal(err)
	}

	k2 := privKeyOrFatal(t)
	k2Name, err := encode("bar")
	if err != nil {
		t.Fatal(err)
	}

	err = ks.Put("foo", k1)
	if err != nil {
		t.Fatal(err)
	}

	err = ks.Put("bar", k2)
	if err != nil {
		t.Fatal(err)
	}

	l, err = ks.List()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(l)
	if l[0] != "bar" || l[1] != "foo" {
		t.Fatal("wrong entries listed")
	}

	if err := assertDirContents(tdir, []string{k1Name, k2Name}); err != nil {
		t.Fatal(err)
	}

	exist, err := ks.Has("foo")
	if !exist {
		t.Fatal("should know it has a key named foo")
	}
	if err != nil {
		t.Fatal(err)
	}

	if err := ks.Delete("bar"); err != nil {
		t.Fatal(err)
	}

	if err := assertDirContents(tdir, []string{k1Name}); err != nil {
		t.Fatal(err)
	}

	if err := assertGetKey(ks, "foo", k1); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("..///foo/", k1); err != nil {
		t.Fatal(err)
	}

	if err := ks.Put("", k1); err == nil {
		t.Fatal("shouldnt be able to put a key with no name")
	}

	if err := ks.Put(".foo", k1); err != nil {
		t.Fatal(err)
	}
}
