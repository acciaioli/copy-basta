package bootstrap

import (
	"log"
	"os"
)

func Bootstrap(destDir string) error {
	err := bootstrap(destDir)
	if err != nil {
		cleanup(destDir)
	}
	return err
}

func bootstrap(destDir string) error {
	if err := os.Mkdir(destDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func cleanup(destDir string) {
	if err := os.RemoveAll(destDir); err != nil {
		log.Print("[ERROR] cleanup fail")
	}
}

const ()
