package main

import (
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("only 1 arg (dir) allowed")
	}
	dir := os.Args[1]
	log.Printf("checking dir %s", dir)
	modcount := make(map[string]map[string]int)
	modfiles := 0

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() == "go.mod" {
				log.Printf("found go.mod at %s", path)
				modfiles++
				dat, err := os.ReadFile(path)
				if err != nil {
					log.Fatalf("error %s reading file %s", err, path)
				}
				modfile, err := modfile.Parse(path, dat, nil)
				if err != nil {
					log.Fatalf("error parsing %s (%s)", path, err)
				}
				for _,r := range modfile.Require {
					if _,ok := modcount[r.Mod.Path] ; !ok {
						modcount[r.Mod.Path] = make(map[string]int)
						modcount[r.Mod.Path][r.Mod.Version] = 0
					}
					modcount[r.Mod.Path][r.Mod.Version] = modcount[r.Mod.Path][r.Mod.Version] + 1
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	log.Printf("anz modfiles %d", modfiles)
	newmodfile, err := modfile.Parse("go.mod", []byte("module golang-build-base\ngo 1.17\n"), nil)

	for modname, vercountmap := range modcount {
		for version, count := range vercountmap {
			if count >= (modfiles / 2) {
				log.Printf("%d * %s@%s", count, modname, version)
				newmodfile.AddNewRequire(modname,version,true)
			}
		}
	}
	if res, err := newmodfile.Format() ; err != nil {
		log.Printf("could not format go.mod: %s", err)
	} else {
		if err := os.WriteFile("go-dep.mod", res, 0644) ; err != nil {
			log.Fatalf("err? %s", err)
		}
	}
}