package fontcache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
)

type Cache map[string]*truetype.Font

func New() Cache {
	return make(Cache)
}

var FontPaths = []string{
	filepath.Join(os.Getenv("windir"), "Fonts"),
	filepath.Join(os.Getenv("localappdata"), "Microsoft", "Windows", "Fonts"),
}

// Init returns a list of all font files found on the system.
func (c Cache) Init(paths []string) {
	for _, dir := range paths {
		filepath.Walk(dir, c.loadFont)
	}
}

func (c Cache) loadFont(path string, info os.FileInfo, err error) error {
	// process ttf files only
	if info.IsDir() == false && strings.HasSuffix(strings.ToLower(path), ".ttf") {
		// if strings.ToLower(path[len(path)-4:]) != ".ttf" {
		// 	return nil
		// }
		ttfBytes, err := ioutil.ReadFile(path)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		fontFace, err := truetype.Parse(ttfBytes)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		name := strings.ToLower(fontFace.Name(truetype.NameIDFontFamily))
		c[name] = fontFace
	}
	return nil
}

// func expandUser(path string) (expandedPath string) {
// 	if strings.HasPrefix(path, "~") {
// 		if u, err := user.Current(); err == nil {
// 			return strings.Replace(path, "~", u.HomeDir, -1)
// 		}
// 	}
// 	return path
// }
