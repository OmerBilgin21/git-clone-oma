package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Env string

const (
	Dev  Env = "DEV"
	Prod Env = "PRODUCTION"
)

type FileIO interface {
	CreateRepoInitInfo(repositoryId int) error
	GetRepositoryId() (int, error)
	WriteToFile(filename string, content string) error
	DeleteFile(filename string) error
}

type FileIOImpl struct {
	currentDirectory string
	infoFile         string
	infoDir          string
	env              Env
}

func NewFileIO() *FileIOImpl {
	currDirr, err := os.Getwd()
	if err != nil {
		panic("please give oma execution rights or check your directory permissions")
	}

	env := os.Getenv("ENV")
	var ourEnv Env
	if Env(env) == Dev {
		fmt.Printf("RUNNING ON DEV\n")
		ourEnv = Dev
	} else {
		ourEnv = Prod
	}

	return &FileIOImpl{
		currentDirectory: currDirr,
		infoFile:         "repository_info.txt",
		infoDir:          ".oma",
		env:              ourEnv,
	}
}

func (repoFileIO *FileIOImpl) CreateRepoInitInfo(repositoryId int) error {
	// permissions in Go are in octal notation apparently
	// hence the 0 prefix
	err := os.Mkdir(".oma", 0755)
	if err != nil {
		return fmt.Errorf("error while creating repository info file parent:\n%v", err)
	}

	file, err := os.Create(filepath.Join(repoFileIO.currentDirectory, repoFileIO.infoDir, repoFileIO.infoFile))
	defer file.Close()

	if err != nil {
		return fmt.Errorf("error while creating repository info file:\n%v", err)
	}

	repoIdInBytes := []byte("repositoryId=" + strconv.Itoa(repositoryId))
	file.Write(repoIdInBytes)

	return nil
}

func (repoFileIO *FileIOImpl) GetRepositoryId() (int, error) {
	contentBytes, err := os.ReadFile(filepath.Join(repoFileIO.infoDir, repoFileIO.infoFile))

	if err != nil {
		return -1, fmt.Errorf("error while reading the info file:\n%v", err)
	}

	content := strings.ReplaceAll(string(contentBytes), "\r", "")
	content = strings.ReplaceAll(content, "\t", "  ")

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "repositoryId") {
			repoIdString := strings.Split(line, "=")[1]
			repoId, err := strconv.ParseInt(repoIdString, 10, 64)

			if err != nil {
				return -1, fmt.Errorf("error while parsing repository ID:\n%v", err)
			}

			return int(repoId), nil
		}
	}

	return -1, fmt.Errorf("repository ID could not be found")
}

func (repoFileIO *FileIOImpl) WriteToFile(filename string, content string) error {
	if repoFileIO.env == Dev {
		fmt.Printf("would have written to file :%v\nthe ingredients:\n%v\n", filename, content)
		return nil
	}
	// TODO: implement actual write to file
	return fmt.Errorf("not implemented yet")
}

func (repoFileIO *FileIOImpl) DeleteFile(filename string) error {
	if repoFileIO.env == Dev {
		fmt.Printf("file: %v would have been deleted\n", filename)
		return nil
	}
	if _, err := os.Stat(filename); err == nil {
		if err := os.RemoveAll(filename); err != nil {
			return fmt.Errorf("error while removing the .oma directory:\n%w", err)
		}

		fmt.Printf("file removed\n")
	}

	return nil
}
